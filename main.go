package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/fsnotify/fsnotify"
)

func main() {
	fmt.Println("golang-ipa-renamer-watch started.")

	// Read config from environment variables
	watchDir := os.Getenv("WATCH_DIR")
	if watchDir == "" {
		watchDir = "."
	}
	ipaRenamer := os.Getenv("IPA_RENAMER")
	if ipaRenamer == "" {
		ipaRenamer = "./ipa_renamer"
	}
	outputDir := os.Getenv("OUTPUT_DIR")
	if outputDir == "" {
		outputDir = watchDir
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					if filepath.Ext(event.Name) == ".ipa" {
						filename := filepath.Base(event.Name)
						matched, _ := regexp.MatchString(`^[^@]+@[^@]+\.ipa$`, filename)
						if !matched {
							cmd := exec.Command(ipaRenamer, event.Name, "-o", outputDir)
							cmd.Stdout = os.Stdout
							cmd.Stderr = os.Stderr
							err := cmd.Run()
							if err != nil {
								log.Printf("Failed to rename %s: %v", filename, err)
							}
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(watchDir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
