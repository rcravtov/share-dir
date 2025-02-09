package views

import (
	"html/template"
	"io"
	"io/fs"
)

type Template struct {
	tmpl *template.Template
}

func New(views fs.FS) *Template {
	tmpl, err := template.New("").ParseFS(views, "views/index.html")
	if err != nil {
		panic(err)
	}
	return &Template{
		tmpl: tmpl,
	}
}

func (t *Template) FileList(w io.Writer, data any) error {
	return t.tmpl.ExecuteTemplate(w, "FileList", data)
}
