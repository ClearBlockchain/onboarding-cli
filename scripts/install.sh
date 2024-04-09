#!/usr/bin/env sh
set -e

echo "Using the OSS distribution..."
REPO_NAME="ClearBlockchain/onboarding-cli"
RELEASES_URL="https://github.com/${REPO_NAME}/releases"
API_RELEASES_URL="https://api.github.com/repos/${REPO_NAME}/releases"
FILE_BASENAME="onboarding-cli"
LATEST="$(curl -sf $API_RELEASES_URL/latest | grep -i 'tag_name' | awk -F '"' '{print $4}')"

test -z "$VERSION" && VERSION="$LATEST"

test -z "$VERSION" && {
	echo "Unable to get glide version." >&2
	exit 1
}

TMP_DIR="$(mktemp -d)"
# shellcheck disable=SC2064 # intentionally expands here
trap "rm -rf \"$TMP_DIR\"" EXIT INT TERM

OS="$(uname -s)"
ARCH="$(uname -m)"
test "$ARCH" = "aarch64" && ARCH="arm64"
TAR_FILE="${FILE_BASENAME}_${OS}_${ARCH}.tar.gz"

(
	cd "$TMP_DIR"
	echo "Downloading glide $VERSION..."
	curl -sfLO "$RELEASES_URL/download/$VERSION/$TAR_FILE"
)

tar -xf "$TMP_DIR/$TAR_FILE" -C "$TMP_DIR"
"$TMP_DIR/glide" "$@"
