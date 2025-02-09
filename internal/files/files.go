package files

import (
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"slices"
)

type FileService struct {
	baseDir string
}

type File struct {
	Name  string
	URL   string
	Path  string
	IsDir bool
}

func New(baseDir string) *FileService {
	return &FileService{
		baseDir: baseDir,
	}
}

func (fserv *FileService) GetFullPath(path string) string {
	return filepath.Join(fserv.baseDir, path)
}

func (fserv *FileService) IsDir(path string) (bool, error) {
	fullPath := fserv.GetFullPath(path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func (fserv *FileService) GetFileList(path string) ([]File, error) {
	dir := os.DirFS(fserv.GetFullPath(path))
	entries, err := fs.ReadDir(dir, ".")
	if err != nil {
		return []File{}, err
	}

	var list []File
	for _, e := range entries {

		name := e.Name()

		url, err := url.JoinPath("files", path, url.QueryEscape(name))
		if err != nil {
			return []File{}, err
		}

		file := File{
			Name:  name,
			URL:   url,
			Path:  filepath.Join(fserv.baseDir, name),
			IsDir: e.IsDir(),
		}

		list = append(list, file)
	}

	// place dirs at the top
	slices.SortStableFunc(list, func(a File, b File) int {
		if a.IsDir && b.IsDir {
			return 0
		} else {
			if a.IsDir {
				return -1
			}
		}
		return 1
	})

	return list, nil
}

func (fserv *FileService) CreateFile(path string, name string, src io.Reader) error {
	fullPath := fserv.GetFullPath(path)
	fullPath = filepath.Join(fullPath, name)

	dst, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
