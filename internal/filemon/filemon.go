package filemon

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

type FileMon struct {
	w       *fsnotify.Watcher
	DataDir string
}

type FileEvent struct {
	Name string
	Size int
}

func NewFileMon(dataDir string) (*FileMon, error) {
	w, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, err
	}

	err = w.Add(dataDir)
	if err != nil {
		return nil, err
	}

	return &FileMon{
		w:       w,
		DataDir: dataDir,
	}, nil
}

func (fm *FileMon) Watch() <-chan FileEvent {
	var ch = make(chan FileEvent)

	go func(w *fsnotify.Watcher, ch chan<- FileEvent) {
		for {
			select {
			case event, ok := <-w.Events:
				log.Println(event.Name, event.Op)
				if !ok || (event.Op != fsnotify.Write && event.Op != fsnotify.Chmod) {
					break
				}

				fh, err := os.Open(event.Name)
				if err != nil {
					log.Println(err)
					break
				}

				info, err := fh.Stat()
				if err != nil {
					log.Println(err)
					break
				}

				ch <- FileEvent{
					Name: event.Name,
					Size: int(info.Size()),
				}
			case err, _ := <-w.Errors:
				if err != nil {
					log.Println("not ok : ", err)
					break
				}
			}
		}
	}(fm.w, ch)

	fmt.Println("watching ", fm.DataDir)
	return ch
}
