package main

import (
	"log"
	"net/http"
	"flag"
	"os"
	"strings"
	"path/filepath"
	"github.com/asafg6/sse_handler"
	"encoding/json"
	"messages"
)

var client_dir string
var httpAddr string
var pushPaths []string
var pubsubClient *messages.PubSubClient

func init(){
	flag.StringVar(&client_dir, "client-dir", "frontend/build" ,"the built app files")
	flag.StringVar(&httpAddr, "http", ":8080", "Listen address")
}


func handleEventsSSE(w http.ResponseWriter, r *http.Request, flusher *sse_handler.MessageFlusher, close <-chan bool) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		channel := pubsubClient.Subscribe("events")
		item := channel.Listen()
		for {
			item = item.GetNextMessageOrWaitWithClose(close)
			if item == nil {
				log.Println("Client disconnect")
				break
			}
			message := item.GetData().(messages.Message)
			bytes, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				break
			}
			eventMessage := sse_handler.EventMessage{Id: message.Id,
					                  	Data: string(bytes),
								Event: message.Type }
			log.Printf("Sending %v", eventMessage)
			flusher.Send(&eventMessage)


		}
		pubsubClient.UnSubscribe("events")
}

func handleFrontend(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || strings.HasSuffix(r.URL.Path, "/") {
		serveIndex(w, r)
		return
	}
	filePath := client_dir + "/" + r.URL.Path[1:]
	log.Println("filePath ", filePath)
	if _, err := os.Stat(filePath); err == nil {
		http.ServeFile(w, r, filePath)
	} else {
		serveIndex(w, r)
	}

}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("serving index.html")
	pusher, ok := w.(http.Pusher)
	if ok {
		log.Println("Push is supported")
		for _, path := range pushPaths {
			if err := pusher.Push(path , nil); err != nil {
				log.Printf("Failed to push: %v", err)
			} else {
				log.Printf("Pushed %v", path)
			}
		}
	}
	http.ServeFile(w, r, client_dir + "/index.html")
}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() && !strings.Contains(path, "index.html") {
		pushPaths = append(pushPaths, strings.TrimPrefix(path, client_dir))
	}
	return nil
}

func fillPushPaths() {
	pushPaths = make([]string, 0)
	err := filepath.Walk(client_dir, visit)
	if err != nil {
		log.Printf("filepath.Walk() returned %v\n", err)
	}
	log.Println(pushPaths)
}

func main() {
	flag.Parse()
	pubsubClient = messages.MakeNewClient("localhost:6379")
	fillPushPaths()
	http.HandleFunc("/messages", sse_handler.HandleSSE(handleEventsSSE))
	http.HandleFunc("/", handleFrontend)
	log.Fatal(http.ListenAndServeTLS(httpAddr, "cert.pem", "key.pem", nil))
}


