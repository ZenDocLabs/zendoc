package export

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dterbah/zendoc/internal/parser"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	Exporter DocExporter
}

func shouldIgnore(path, docPath string) bool {
	absPath, _ := filepath.Abs(path)
	absDocPath, _ := filepath.Abs(docPath)

	if strings.HasPrefix(absPath, absDocPath) {
		return true
	}
	if filepath.Base(absPath) == "doc.json" {
		return true
	}
	return false
}

func (watcher FileWatcher) WatchDir(docParser parser.DocParser, dirName, docPath string) error {
	color.Green("Mode watched activated !")

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %w", err)
	}
	defer w.Close()

	err = filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && !shouldIgnore(path, docPath) {
			containsGoFile := false

			entries, err := os.ReadDir(path)
			if err != nil {
				return err
			}

			for _, entry := range entries {
				if !entry.IsDir() && filepath.Ext(entry.Name()) == ".go" {
					containsGoFile = true
					break
				}
			}

			if containsGoFile {
				err = w.Add(path)
				if err != nil {
					color.Red("failed to watch directory %s: %s", path, err)
				} else {
					color.Cyan("Watching directory: %s", path)
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	done := make(chan bool)
	errChan := make(chan error)

	var debounceTimer *time.Timer
	var debounceDelay = 500 * time.Millisecond
	var debounceChan = make(chan struct{}, 1)

	go func() {
		for range debounceChan {
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			debounceTimer = time.AfterFunc(debounceDelay, func() {
				color.Green("📝 Debounced export triggered")

				doc, err := docParser.ParseDocForDir(dirName, "")
				if err != nil {
					errChan <- fmt.Errorf("error during parsing: %w", err)
					return
				}

				err = watcher.Exporter.Export(*doc)
				if err != nil {
					errChan <- fmt.Errorf("error during export: %w", err)
				}
			})
		}
	}()

	go func() {
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					errChan <- fmt.Errorf("watcher event channel closed unexpectedly")
					return
				}

				if shouldIgnore(event.Name, docPath) {
					color.Magenta("Ignored event on: %s", event.Name)
					continue
				}

				color.Yellow("File event: %s", event)

				if event.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Remove|fsnotify.Rename) != 0 {
					color.Blue("Change detected on: %s", event.Name)

					select {
					case debounceChan <- struct{}{}:
					default:
					}
				}

			case err, ok := <-w.Errors:
				if !ok {
					errChan <- fmt.Errorf("watcher error channel closed unexpectedly")
					return
				}
				errChan <- fmt.Errorf("watcher error: %w", err)
				return
			}
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-done:
		return nil
	}
}
