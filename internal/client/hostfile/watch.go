package hostfile

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type Watcher struct {
	watcher  *fsnotify.Watcher
	OnChange chan string
	quit     chan bool
}

func (w *Watcher) StopWatch() {
	tools.Error("Stop hosts file watcher")
	w.watcher.Close()
	close(w.OnChange)
	w.quit <- true
	close(w.quit)
	tools.Error("hosts file watcher is stop")
}

func StartWatch(file string) *Watcher {
	fsWatch, err := fsnotify.NewWatcher()
	if err != nil {
		tools.Die("Failed create fsnotify", err.Error())
	}
	file = filepath.Clean(file)

	w := Watcher{
		watcher:  fsWatch,
		OnChange: make(chan string, 1),
		quit:     make(chan bool, 1),
	}

	go func() {
		emod := fsnotify.Write + fsnotify.Create + fsnotify.Remove
		for {
			select {
			case event, ok := <-fsWatch.Events:
				if !ok {
					return
				}
				// tools.Error("fsnotify event: %s", spew.Sdump(event))
				if event.Name == file && (event.Op&emod) != 0 {
					log.Println("modified file:", event.Name)
					if data, err := ioutil.ReadFile(file); err == nil {
						w.OnChange <- string(data)
					} else {
						w.OnChange <- ""
					}
				}
			case err, ok := <-fsWatch.Errors:
				if !ok {
					return
				}
				log.Println("fsnotify error:", err)
			case _ = <-w.quit:
				log.Println("fsnotify finishing...")
				return
			}
		}
	}()

	err = fsWatch.Add(filepath.Dir(file))
	if err != nil {
		tools.Die("failed add watch dir: %s", err.Error())
	}

	if data, err := ioutil.ReadFile(file); err == nil {
		err = fsWatch.Add(file)
		if err != nil {
			tools.Die("failed add watch file: %s", err.Error())
		}
		w.OnChange <- string(data)
	}

	return &w
}
