package main

import (
	"log"
	"net/http"

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

	fs := http.FileServer(http.Dir("./dist"))
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
