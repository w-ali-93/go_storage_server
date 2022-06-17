package server

import (
	"errors"
	"fmt"
	storage "go_storage_server/storage"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

const maxUploadSize = 512 * 1024 // 512 kb

func LogRequest(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func UploadFileHandler(store storage.FileStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// limit upload size
		r.Body = http.MaxBytesReader(w, r.Body, 2*512*1024)

		// serve webpage for testing
		if r.Method == "GET" {
			t, _ := template.ParseFiles("upload.html")
			t.Execute(w, nil)
			return
		}

		// parse entire multipart form in one go
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			fmt.Printf("Could not parse multipart form: %v\n", err)
			generateError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
			return
		}

		// parse and validate file and file header from form data
		file, fileHeader, err := r.FormFile("uploadFile")
		if err != nil {
			generateError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		defer file.Close()
		fileSize := fileHeader.Size
		if fileSize > maxUploadSize {
			generateError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			generateError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		// parse and validate file type
		detectedFileType := http.DetectContentType(fileBytes)
		switch detectedFileType {
		case "application/zip":
			break
		default:
			generateError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
			return
		}

		// upload file to storage
		if fileName, err := store.Upload(fileBytes); err != nil {
			generateError(w, err.Error(), http.StatusBadRequest)
		} else {
			w.Write([]byte(fileName))
		}
	})
}

func DownloadFileHandler(store storage.FileStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse requested fileName from request
		fileName, err := parseFileName(r)
		if err != nil {
			generateError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// download file from storage
		if fc, err := store.Download(fileName); err != nil {
			generateError(w, err.Error(), http.StatusNotFound)
			return
		} else {
			w.Write(fc)
		}
	})
}

func parseFileName(r *http.Request) (string, error) {
	fileNameRaw, ok := r.URL.Query()["filename"]
	if !ok || len(fileNameRaw) == 0 {
		return "", errors.New("MISSING_OR_INVALID_FILENAME")
	}
	return fileNameRaw[0], nil
}

func generateError(w http.ResponseWriter, message string, statusCode int) {
	log.Println("ERROR:", message)
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
