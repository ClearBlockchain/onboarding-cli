package cmd

import (
	"fmt"

	ui "github.com/ClearBlockchain/onboarding-cli/pkg/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// version as provided by goreleaser
	Version = ""

	// CommitSHA as provided by goreleaser
	CommitSHA = ""

	rootCmd = &cobra.Command{
		Use:   "glide",
		Short: "Manage global Network APIs, 5G, Edge resources with Glide",
		Long: ui.Paragraph.Render(
			fmt.Sprintf("\nExplore %s Integration Layer, a one-stop infrastructure offering access to Network APIs, 5G, and edge resources from worldwide CSPs and perform time-sensitive SIM Swap checks using a unified interface and API", ui.Keyword.Render("Glide")),
		),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(GetLongDescription())

			err := cmd.Help()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func Execute() {
	if len(CommitSHA) >= 7 {
		vt := rootCmd.VersionTemplate()
		rootCmd.SetVersionTemplate(vt[:len(vt)-1] + " (" + CommitSHA[0:7] + ")\n")
	}
	if Version == "" {
		Version = "unknown (built from source)"
	}

	rootCmd.Version = Version

	// subcommands
	rootCmd.AddCommand(LoginCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
