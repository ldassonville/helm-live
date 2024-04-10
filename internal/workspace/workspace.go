package workspace

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
)

type Workspace struct {
	Path string

	OnChangeFnc func()
}

func (w *Workspace) Live() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					w.OnChangeFnc()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}

	}()

	w.watchRecursive(w.Path, watcher)
	defer watcher.Close()

	// Block main goroutine forever.
	<-make(chan struct{})
}

func (w *Workspace) watchRecursive(root string, watcher *fsnotify.Watcher) {

	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				if err := watcher.Add(path); err != nil {
					log.Fatal(err)
				}
			}
			return nil
		})

	if err != nil {
		log.Println(err)
	}
}
