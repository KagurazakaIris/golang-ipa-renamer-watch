package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	plist "howett.net/plist"
)

type Config struct {
	Glob     string
	Template string
	Out      string
	Temp     string
}

// func main() removed to avoid redeclaration error.
// Move this logic to another function if needed, or ensure only one main() exists in your package.

func renameIPA(cfg Config, path string) error {
	plistPath, err := getInfoPlist(cfg, path)
	if err != nil {
		return err
	}
	defer os.Remove(plistPath)

	f, err := os.Open(plistPath)
	if err != nil {
		return err
	}
	defer f.Close()

	var info map[string]interface{}
	decoder := plist.NewDecoder(f)
	if err := decoder.Decode(&info); err != nil {
		return err
	}
	cfBundleIdentifier, ok := info["CFBundleIdentifier"].(string)
	if !ok {
		return fmt.Errorf("CFBundleIdentifier not found")
	}

	rawName := strings.TrimSuffix(filepath.Base(path), ".ipa")
	if rawName == "" {
		rawName = "unknown"
	}
	// 固定命名规则：原文件名@CFBundleIdentifier.ipa
	newName := fmt.Sprintf("%s@%s.ipa", rawName, cfBundleIdentifier)

	newPath := filepath.Join(cfg.Out, newName)
	if _, err := copyFile(path, newPath); err != nil {
		return err
	}
	fmt.Printf("[renamed] %s to %s\n", path, newPath)
	return nil
}

func getInfoPlist(cfg Config, ipaPath string) (string, error) {
	r, err := zip.OpenReader(ipaPath)
	if err != nil {
		return "", err
	}
	defer r.Close()
	for _, f := range r.File {
		if !f.FileInfo().IsDir() && matchPlist(f.Name) {
			outPath := filepath.Join(cfg.Temp, "Info.plist")
			outFile, err := os.Create(outPath)
			if err != nil {
				return "", err
			}
			inFile, err := f.Open()
			if err != nil {
				outFile.Close()
				return "", err
			}
			_, err = io.Copy(outFile, inFile)
			inFile.Close()
			outFile.Close()
			if err != nil {
				return "", err
			}
			return outPath, nil
		}
	}
	return "", fmt.Errorf("info.plist not found")
}

func matchPlist(path string) bool {
	return strings.Count(path, "/") == 2 && strings.HasSuffix(path, "/Info.plist")
}

func copyFile(src, dst string) (int64, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()
	return io.Copy(dstFile, srcFile)
}
