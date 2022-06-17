package main

import (
	server "go_storage_server/server"
	storage "go_storage_server/storage"
	"log"
	"net/http"
)

func main() {
	s, err := storage.NewPersistentStore("archives")
	if err != nil {
		panic("cannot initialize storage layer")
	}

	http.Handle("/upload", server.UploadFileHandler(s))
	http.Handle("/download", server.DownloadFileHandler(s))
	log.Print("Server started on localhost:8080, use /upload for uploading files and /download?fileName=<fileName> for downloading")
	log.Fatal(http.ListenAndServe(":8080", server.LogRequest(http.DefaultServeMux)))
}
