package hostfile

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type Watcher struct {
	watcher  *fsnotify.Watcher
	OnChange chan bool
	quit     chan bool
	qDispose func()
	current  string
	file     string
}

func (w *Watcher) StopWatch() {
	if w.qDispose != nil {
		w.qDispose()
	}
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
		file:     file,
		OnChange: make(chan bool),
		quit:     make(chan bool, 1),
	}

	go func() {
		debounced := debounce.New(1 * time.Second)
		dfunc := func() {
			tools.Debug("[FsWatch] modified file")
			if data, err := ioutil.ReadFile(file); err == nil {
				w.current = string(data)
				w.notifyChange()
			} else {
				tools.Debug("[FsWatch] fail read file: %s", err)
			}
		}

		emod := fsnotify.Write + fsnotify.Create + fsnotify.Remove
		for {
			select {
			case event, ok := <-fsWatch.Events:
				if !ok {
					return
				}
				// tools.Error("fsnotify event: %s", spew.Sdump(event))
				if event.Name == file && (event.Op&emod) != 0 {
					debounced(dfunc)
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
		w.current = string(data)
		w.notifyChange()
	}

	w.qDispose = tools.WaitExit(func(int) {
		w.StopWatch()
	})

	return &w
}

func (w *Watcher) notifyChange() {
	select {
	case w.OnChange <- true:
	default:
	}
}
