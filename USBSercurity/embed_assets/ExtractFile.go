package embed_assets

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed  assets/*
var EmbeddedFiles embed.FS // ✅ viết hoa để export

func ExtractAllAssetsTo(destDir string) (map[string]string, error) {
	extracted := make(map[string]string)

	err := fs.WalkDir(EmbeddedFiles, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		// Read embedded file
		data, err := EmbeddedFiles.ReadFile(path)
		if err != nil {
			return err
		}

		// Clean filename (strip "assets/")
		filename := strings.TrimPrefix(path, "assets/")
		destPath := filepath.Join(destDir, filename)

		// Write to disk
		err = os.WriteFile(destPath, data, 0755)
		if err != nil {
			return err
		}

		extracted[filename] = destPath
		return nil
	})

	if err != nil {
		return nil, err
	}
	return extracted, nil
}
