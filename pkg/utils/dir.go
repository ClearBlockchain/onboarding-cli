package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func CopyDirectory(scrDir, dest string, grace bool) error {
    entries, err := os.ReadDir(scrDir)
    if err != nil {
        return err
    }
	log.Debugf("Copying directory: %s to %s", scrDir, dest)
    for _, entry := range entries {
        sourcePath := filepath.Join(scrDir, entry.Name())
        destPath := filepath.Join(dest, entry.Name())

        fileInfo, err := os.Stat(sourcePath)
        if err != nil {
            log.Debugf("failed to get file info for '%s'", sourcePath)
			if grace {
				continue
			} else {
				return err
			}
        }

        stat, ok := fileInfo.Sys().(*syscall.Stat_t)
        if !ok {
			err := fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
			if grace {
				log.Debugf("%s", err)
				continue
			} else {
				return err
			}
        }

        switch fileInfo.Mode() & os.ModeType{
        case os.ModeDir:
            if err := CreateIfNotExists(destPath, 0755); err != nil {
                log.Debugf("failed to create directory: '%s'", destPath)
				if grace {
					continue
				} else {
					return err
				}
            }
            if err := CopyDirectory(sourcePath, destPath, grace); err != nil {
                log.Debugf("failed to copy directory: '%s'", sourcePath)
				if grace {
					continue
				} else {
					return err
				}
            }
        case os.ModeSymlink:
            if err := CopySymLink(sourcePath, destPath); err != nil {
                log.Debugf("failed to copy symlink: '%s'", sourcePath)
				if grace {
					continue
				} else {
					return err
				}
            }
        default:
            if err := Copy(sourcePath, destPath); err != nil {
                log.Debugf("failed to copy file: '%s'", sourcePath)
				if grace {
					continue
				} else {
					return err
				}
            }
        }

        if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
            log.Debugf("failed to change owner for '%s'", destPath)
			if grace {
				continue
			} else {
				return err
			}
        }

        fInfo, err := entry.Info()
        if err != nil {
            log.Debugf("failed to get file info for '%s'", sourcePath)
			if grace {
				continue
			} else {
				return err
			}
        }

        isSymlink := fInfo.Mode()&os.ModeSymlink != 0
        if !isSymlink {
            if err := os.Chmod(destPath, fInfo.Mode()); err != nil {
                log.Debugf("failed to change mode for '%s'", destPath)
				if grace {
					continue
				} else {
					return err
				}
            }
        }
    }
    return nil
}

func Copy(srcFile, dstFile string) error {
    out, err := os.Create(dstFile)
    if err != nil {
        return err
    }

    defer out.Close()

    in, err := os.Open(srcFile)
    if err != nil {
        return err
    }

    defer in.Close()

    _, err = io.Copy(out, in)
    if err != nil {
        return err
    }

    return nil
}

func Exists(filePath string) bool {
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return false
    }

    return true
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
    if Exists(dir) {
        return nil
    }

    if err := os.MkdirAll(dir, perm); err != nil {
        return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
    }

    return nil
}

func CopySymLink(source, dest string) error {
    link, err := os.Readlink(source)
    if err != nil {
        return err
    }
    return os.Symlink(link, dest)
}
