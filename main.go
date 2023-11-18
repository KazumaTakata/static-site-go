package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/yuin/goldmark"
)

func main() {

	tmpl, err := template.ParseFiles("template/template.html", "template/header.html")

	if err != nil {
		panic("Parse failed!!")
	}

	dat, _ := os.ReadFile("source/sample.md")
	var buf bytes.Buffer
	if err := goldmark.Convert(dat, &buf); err != nil {
		panic(err)
	}

	var htmlOutputBuffer bytes.Buffer

	tmpl.Execute(&htmlOutputBuffer, template.HTML(buf.String()))

	fmt.Println(string(htmlOutputBuffer.Bytes()))

	os.WriteFile("dist/index.html", htmlOutputBuffer.Bytes(), 0644)

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

				dat, _ := os.ReadFile("source/sample.md")
				var buf bytes.Buffer
				if err := goldmark.Convert(dat, &buf); err != nil {
					panic(err)
				}

				var htmlOutputBuffer bytes.Buffer

				tmpl.Execute(&htmlOutputBuffer, template.HTML(buf.String()))

				os.WriteFile("dist/index.html", htmlOutputBuffer.Bytes(), 0644)

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
