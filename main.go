package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

// Precompile the IPA filename pattern for performance
var ipaPattern = regexp.MustCompile(`^[^@]+@[^@]+\.ipa$`)

// Config struct and renameIPA, getInfoPlist, matchPlist, copyFile are imported from ipa_renamer.go

func main() {
	watchDir := os.Getenv("WATCH_DIR")
	if watchDir == "" {
		watchDir = "./watched"
	}
	outputDir := os.Getenv("OUTPUT_DIR")
	if outputDir == "" {
		outputDir = "./output"
	}
	template := os.Getenv("TEMPLATE")
	if template == "" {
		template = "$raw@$CFBundleIdentifier"
	}
	tempDir := os.Getenv("TEMP_DIR")
	if tempDir == "" {
		tempDir = "./temp"
	}
	cfg := Config{
		Template: template,
		Out:      outputDir,
		Temp:     tempDir,
	}

	log.Printf("[INFO] golang-ipa-renamer-watch starting with\nWATCH_DIR=%s,\nTEMPLATE=%s,\nOUTPUT_DIR=%s,\nTEMP_DIR=%s", watchDir, template, outputDir, tempDir)

	if err := os.MkdirAll(watchDir, 0755); err != nil {
		log.Fatalf("[ERROR] failed to create watch dir: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("[ERROR] failed to create output dir: %v", err)
	}
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		log.Fatalf("[ERROR] failed to create temp dir: %v", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("[ERROR] failed to create watcher: %v", err)
	}
	defer watcher.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write != 0 {
					log.Printf("[DEBUG] write event for %s", event.Name)
					if filepath.Ext(event.Name) == ".ipa" {
						filename := filepath.Base(event.Name)
						if ipaPattern.MatchString(filename) {
							log.Printf("[DEBUG] skipping already correct name %s", filename)
							continue
						}
						cfg.Glob = event.Name // single file
						if err := renameIPA(cfg, event.Name); err != nil {
							log.Printf("[ERROR] failed to rename %s: %v", filename, err)
						} else {
							log.Printf("[INFO] renamed %s successfully", filename)
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("[ERROR] watcher error: %v", err)
			case sig := <-sigs:
				log.Printf("[INFO] received signal %s, shutting down", sig)
				done <- true
				return
			}
		}
	}()

	err = watcher.Add(watchDir)
	if err != nil {
		log.Fatalf("[ERROR] failed to watch directory %s: %v", watchDir, err)
	}

	log.Printf("[INFO] monitoring directory %s for .ipa changes", watchDir)

	<-done
	log.Println("[INFO] shutdown complete")
}
