package watcher

import (
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

//Task stores a pointer to FSWatcher from fsnotify and the paths we will be watching
type Task struct {
	FSWatcher *fsnotify.Watcher
	Paths     []string
}

//Create builds our Task struct and creates our NewWatcher from fsNotify
func Create(p []string) (Task, error) {
	t := Task{
		Paths: p,
	}
	w, err := fsnotify.NewWatcher()
	t.FSWatcher = w
	return t, err
}

//Run watches the directories insides the Path value of the passed in Task, it will block until it detects changes.
//An error will be returned if invalid paths are set
func (t *Task) Run() error {
	for _, v := range t.Paths { //Loop through all our paths
		err := t.FSWatcher.Add(v) //Add path to watcher
		if err != nil {
			return err // Pass fail back to caller, we can't continue from this.
		}
	}
	defer t.FSWatcher.Close()

	ticker := time.NewTicker(time.Millisecond * 500) //Run every 0.5 secs
	defer ticker.Stop()

	var evts []fsnotify.Event //Create a slice to temp store all events
	var mux sync.Mutex        //Create a mux for our append lock

	for {
		select {
		case <-ticker.C:
			if len(evts) == 0 { //Not events added on this tick
				continue
			}
			return nil //Events found, stop blocking and return to program flow
		case event := <-t.FSWatcher.Events:
			mux.Lock()
			evts = append(evts, event) //Append the events to the event slice, this isn't thread-safe so we need a mux
			mux.Unlock()
		}
	}
}
