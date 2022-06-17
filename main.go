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
	if storageType != "persistent" && storageType != "volatile" {
		panic("invalid storage type specified")
	}

	fileStorePersistent, err := storage.NewPersistentStore("archives")
	fileStoreVolatile := storage.NewVolatileStore()
	if err != nil {
		panic("cannot initialize storage layer")
	}

	http.Handle("/upload", server.UploadFileHandler(fileStorePersistent, fileStoreVolatile, storageType))
	http.Handle("/download", server.DownloadFileHandler(fileStorePersistent, fileStoreVolatile, storageType))
	log.Printf("Server started on localhost:8080 in %s storage mode, use /upload for uploading files and /download?fileName=<fileName> for downloading", storageType)
	log.Fatal(http.ListenAndServe(":8080", server.LogRequest(http.DefaultServeMux)))
}
