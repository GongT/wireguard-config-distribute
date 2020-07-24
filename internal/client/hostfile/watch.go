package hostfile

import (
	"io/ioutil"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type watcher struct {
	watcher  *fsnotify.Watcher
	OnChange chan string
}

func (w *watcher) StopWatch() {
	w.watcher.Close()
}

func StartWatch(file string) watcher {
	fsWatch, err := fsnotify.NewWatcher()
	if err != nil {
		tools.Die("Failed create fsnotify", err.Error())
	}

	w := watcher{
		fsWatch,
		make(chan string, 1),
	}

	if data, err := ioutil.ReadFile(file); err == nil {
		w.OnChange <- string(data)
	}

	go func() {
		for {
			select {
			case event, ok := <-fsWatch.Events:
				if !ok {
					return
				}
				log.Println("fsnotify event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
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
			}
		}
	}()

	err = fsWatch.Add(file)
	if err != nil {
		log.Fatal(err)
	}

	return w
}
