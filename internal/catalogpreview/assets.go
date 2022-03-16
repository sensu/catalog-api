package catalogpreview

import (
	"embed"
	"net/http"
	"path/filepath"
)

//go:generate go run ./cmd/update_assets

//go:embed assets/*
var assetsFS embed.FS
var assetsHTTPFS = chrootableFS{fs: http.FS(assetsFS)}

type FS interface {
	http.FileSystem
	Chroot(dir string) FS
}

type chrootableFS struct {
	fs   http.FileSystem
	root string
}

func (cfs *chrootableFS) Open(name string) (http.File, error) {
	path := filepath.Join(cfs.root, name)
	return cfs.fs.Open(path)
}

func (fs *chrootableFS) Chroot(dir string) *chrootableFS {
	return &chrootableFS{
		fs:   fs,
		root: dir,
	}
}
