package worldmanager

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log/slog"
	"os"
	"path/filepath"
	"server/server/utils"

	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
)

func World(dir string, isZip bool, log *slog.Logger) (*world.World, error) {
	var folder string
	var err error
	if isZip {
		folder, err = unzipFile(dir)
	} else {
		folder = dir
		err = nil
	}
	if err != nil {
		return nil, err
	}
	conf := mcdb.Config{Log: log}
	provider, err := conf.Open(folder)
	if err != nil {
		return nil, err
	}

	return world.Config{
		Log:      log,
		Dim:      world.Overworld,
		Provider: provider,
		ReadOnly: true,
		Entities: entity.DefaultRegistry,
	}.New(), nil
}

func unzipFile(zipFile string) (string, error) {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return "", err
	}
	defer utils.Panic(r.Close())

	dest, err := os.MkdirTemp("", "world-*")
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return "", err
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer utils.Panic(rc.Close())

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, f.Mode())
			if err != nil {
				return err
			}
		} else {
			err = os.MkdirAll(filepath.Dir(path), f.Mode())
			if err != nil {
				return err
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer utils.Panic(f.Close())
			if _, err = io.Copy(f, rc); err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return "", err
		}
	}

	return dest, nil
}

func zipDir(dir string) (string, error) {
	archive, err := os.CreateTemp("", fmt.Sprintf("%v.*.zip", filepath.Dir(dir)))
	if err != nil {
		return "", err
	}
	defer utils.Panic(archive.Close())
	w := zip.NewWriter(archive)
	defer utils.Panic(w.Close())

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// add a trailing slash for creating dir
			path = fmt.Sprintf("%s%c", path, os.PathSeparator)
			_, err = w.Create(path)
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer utils.Panic(file.Close())
		// Ensure that `path` is not absolute; it should not start with "/".
		// This snippet happens to work because I don't use
		// absolute paths, but ensure your real-world code
		// transforms path into a zip-root relative path.
		f, err := w.Create(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}
		return nil
	}
	if err = filepath.Walk(dir, walker); err != nil {
		return "", err
	}

	return archive.Name(), nil
}

func CopyAndRenameFolder(src, dst string, newName string) error {
	// Get source folder information
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source folder: %w", err)
	}

	// Check if source is a directory
	if !srcInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", src)
	}

	// Create destination folder with new name
	dstPath := filepath.Join(dst, newName)
	if err := os.MkdirAll(dstPath, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create destination folder: %w", err)
	}

	// Walk through the source folder
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the source folder itself
		if path == src {
			return nil
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		newPath := filepath.Join(dstPath, relPath)

		// Check if it's a file or directory
		if info.IsDir() {
			// Create subdirectory in destination
			if err := os.MkdirAll(newPath, info.Mode()); err != nil {
				return fmt.Errorf("failed to create subdirectory: %w", err)
			}
		} else {
			// Copy the file
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			if err := ioutil.WriteFile(newPath, data, info.Mode()); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
		}

		return nil
	})

	return err
}
