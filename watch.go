package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsevents"
)

func watch(root string, es *fsevents.EventStream, names chan<- string, done chan<- error, reflexes []*Reflex) {
	rootDir, err := filepath.Abs(root)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		select {
		case msg := <-es.Events:
			for _, e := range msg {
				fullPath := "/" + e.Path
				if verbose {
					infoPrintln(-1, "fsnotify event:", fullPath, e.Flags)
				}
				stat, err := os.Stat(fullPath)
				if err != nil {
					continue
				}
				path := normalize(fullPath, rootDir, stat.IsDir())
				names <- path
			}
			// TODO: Cannot currently remove fsnotify watches
			// recursively, or for deleted files. See:
			// https://github.com/cespare/reflex/issues/13
			// https://github.com/go-fsnotify/fsnotify/issues/40
			// https://github.com/go-fsnotify/fsnotify/issues/41
			// case err := <-watcher:
			// 	done <- err
			// 	return
			// }
		}
	}
}

func normalize(path, rootDir string, dir bool) string {
	path = strings.TrimPrefix(path, rootDir)
	path = strings.TrimPrefix(path, "./")
	if dir && !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	return path
}
