package main

import (
	"embed"
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"share-dir/internal/files"
	"share-dir/internal/handler"
	"share-dir/internal/views"
)

//go:embed views/index.html
var viewsFS embed.FS

//go:embed static/*
var staticFS embed.FS

func main() {

	var addr, path, prefix string

	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exeDir := filepath.Dir(exePath)

	flag.StringVar(&addr, "addr", ":3000", "Server listning address")
	flag.StringVar(&path, "path", exeDir, "Working directory")
	flag.StringVar(&prefix, "prefix", "/", "URL prefix")
	flag.Parse()

	templateService := views.New(viewsFS)
	fileService := files.New(path)
	handlerService := handler.New(fileService, templateService, prefix)

	http.HandleFunc("GET /", handlerService.ListFiles)
	http.HandleFunc("POST /", handlerService.CreateFile)

	staticHandler := handler.DenyDirectoryListing(http.FileServerFS(staticFS))
	http.Handle("GET /static/", staticHandler)

	http.ListenAndServe(addr, nil)
}
