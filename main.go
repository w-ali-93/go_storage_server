package main

import (
	"flag"
	server "go_storage_server/server"
	storage "go_storage_server/storage"
	"log"
	"net/http"
)

func main() {
	var storageType string
	flag.StringVar(&storageType, "storage_type", "persistent", "type of storage to use")
	flag.Parse()
	var fileStore storage.FileStore
	if storageType == "persistent" {
		var err error
		if fileStore, err = storage.NewPersistentStore("archives"); err != nil {
			panic("cannot initialize storage layer")
		}
	} else if storageType == "volatile" {
		fileStore = storage.NewVolatileStore()
	} else {
		panic("invalid storage type specified")
	}

	http.Handle("/upload", server.UploadFileHandler(fileStore))
	http.Handle("/download", server.DownloadFileHandler(fileStore))
	log.Printf("Server started on localhost:8080 in %s storage mode, use /upload for uploading files and /download?fileName=<fileName> for downloading", storageType)
	log.Fatal(http.ListenAndServe(":8080", server.LogRequest(http.DefaultServeMux)))
}
