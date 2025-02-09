package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"share-dir/internal/files"
	"share-dir/internal/views"
	"strings"
)

type Handler struct {
	fileService     *files.FileService
	templateService *views.Template
	prefix          string
}

type FileListData struct {
	Prefix   string
	FileList []files.File
}

func New(f *files.FileService, t *views.Template, prefix string) *Handler {
	return &Handler{
		fileService:     f,
		templateService: t,
		prefix:          prefix,
	}
}

func PathFromURL(urlPath string) (string, error) {
	path, err := url.QueryUnescape(urlPath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	path = strings.TrimPrefix(path, "/files/")
	path = strings.TrimSuffix(path, "/")
	return path, nil
}

func (h Handler) ListFiles(w http.ResponseWriter, r *http.Request) {

	path, err := PathFromURL(r.URL.Path)
	if err != nil {
		fmt.Println(err)
		return
	}

	isDir, err := h.fileService.IsDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	if isDir {
		list, err := h.fileService.GetFileList(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		data := FileListData{
			Prefix:   h.prefix,
			FileList: list,
		}
		err = h.templateService.FileList(w, data)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	} else {
		fullPath := h.fileService.GetFullPath(path)
		http.ServeFile(w, r, fullPath)
	}
}

func (h Handler) CreateFile(w http.ResponseWriter, r *http.Request) {

	path, err := PathFromURL(r.URL.Path)
	if err != nil {
		fmt.Println(err)
		return
	}

	isDir, err := h.fileService.IsDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !isDir {
		fmt.Println("not a directory")
		return
	}

	r.ParseMultipartForm(10 << 20)
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	err = h.fileService.CreateFile(path, header.Filename, file)
	if err != nil {
		fmt.Println(err)
		return
	}

	h.ListFiles(w, r)
}
