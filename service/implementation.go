package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func (s Service) GetHandler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/", s.handleFile)

	var handler http.Handler
	handler = m
	return handler
}

func (s Service) handleFile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // download file
		http.FileServer(http.Dir(s.rootDir)).ServeHTTP(w, r)
	case "PUT": // upload file
		// read body and write it to filesystem
		out, err := os.Create(filepath.Join(s.rootDir, r.URL.Path))
		if err != nil {
			fmt.Printf("failed to create target file.\n")
			http.Error(w, fmt.Sprintf("failed to create target file."), http.StatusBadRequest)
		}
		defer out.Close()
		_, err = io.Copy(out, r.Body)
		if err != nil {
			fmt.Printf("failed to write to target file: %s\n", err)
			http.Error(w, fmt.Sprintf("failed to write to target file: %s", err), http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusNoContent)
	default: // return error
		http.Error(w, fmt.Sprintf("method %s not supported, only GET and PUT", r.Method), http.StatusBadRequest)
	}
}
