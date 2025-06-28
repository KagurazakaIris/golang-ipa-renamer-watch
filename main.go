package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

// Precompile the IPA filename pattern for performance
var ipaPattern = regexp.MustCompile(`^[^@]+@[^@]+\.ipa$`)

func main() {
	// INFO: print startup environment
	watchDir := os.Getenv("WATCH_DIR")
	if watchDir == "" {
		watchDir = "./watched"
	}
	iparenamer := os.Getenv("IPA_RENAMER")
	if iparenamer == "" {
		iparenamer = "./ipa_renamer"
	}
	outputDir := os.Getenv("OUTPUT_DIR")
	if outputDir == "" {
		outputDir = "./output"
	}
	log.Printf("[INFO] golang-ipa-renamer-watch starting with\nWATCH_DIR=%s,\nIPA_RENAMER=%s,\nOUTPUT_DIR=%s", watchDir, iparenamer, outputDir)

	// Initialize file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("[ERROR] failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Setup signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool)

	// event processing
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// DEBUG: capture file write events
				if event.Op&fsnotify.Write != 0 {
					log.Printf("[DEBUG] write event for %s", event.Name)
					if filepath.Ext(event.Name) == ".ipa" {
						filename := filepath.Base(event.Name)
						if ipaPattern.MatchString(filename) {
							log.Printf("[DEBUG] skipping already correct name %s", filename)
							continue
						}
						// INFO: renaming action
						cmd := exec.Command(iparenamer, "-o", outputDir, "--", event.Name)
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						err := cmd.Run()
						if err != nil {
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

	// Begin watching the target directory
	err = watcher.Add(watchDir)
	if err != nil {
		log.Fatalf("[ERROR] failed to watch directory %s: %v", watchDir, err)
	}

	// INFO: monitoring started
	log.Printf("[INFO] monitoring directory %s for .ipa changes", watchDir)

	<-done
	log.Println("[INFO] shutdown complete")
}
