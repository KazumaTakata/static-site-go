package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fsnotify/fsnotify"
)

func main() {

	renderHTML()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	defer watcher.Close()

	fileChanged := make(chan bool)

	done := make(chan bool)
	go func() {
		defer close(done)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				renderHTML()

				log.Printf("%s %s\n", event.Name, event.Op)

				fileChanged <- true

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}

	}()

	err = watcher.Add("./source")
	if err != nil {
		log.Fatal("Add failed:", err)
	}

	fs := http.FileServer(HTMLDir{http.Dir("dist")})
	http.Handle("/", fs)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, fileChanged)
	})

	log.Print("Listening on :3000...")

	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}

}

type HTMLDir struct {
	d http.Dir
}

func (d HTMLDir) Open(name string) (http.File, error) {
	// Try name as supplied
	f, err := d.d.Open(name)
	if os.IsNotExist(err) {
		// Not found, try with .html
		if f, err := d.d.Open(name + ".html"); err == nil {
			return f, nil
		}
	}
	return f, err
}
