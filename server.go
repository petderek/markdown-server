package markdown_server

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// MarkdownServer is a file server that converts any markdown files
// into html.
type MarkdownServer struct {
	Root       http.FileSystem
	Extensions []string
	IndexFile  string
}

const (
	dummyFileName = "foo.html"
)

func (m *MarkdownServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
		r.URL.Path = path
	}

	var redirected bool
redirect:
	if strings.HasSuffix(path, "/") {
		// todo rediect
		path = path + m.IndexFile
	}

	f := m.tryFiles(path)
	if f == nil {
		notFound(w, r, "file not found")
		return
	}

	d, err := f.Stat()
	if err != nil {
		notFound(w, r, "stat error: "+err.Error())
		return
	}

	if d.IsDir() {
		if redirected {
			notFound(w, r, "entered the void: "+path)
			return
		}
		// todo redirect
		redirected = true
		path = path + "/"
		goto redirect
	}

	c, err := ioutil.ReadAll(f)
	if err != nil {
		// idk what to do here
		notFound(w, r, "idk: "+err.Error())
		return
	}
	content := renderMarkdown(c)

	http.ServeContent(w, r, dummyFileName, d.ModTime(), content)
}

// try index.md, index.mdown -- any extension in the list
func (m *MarkdownServer) tryFiles(base string) http.File {
	files := []string{base}
	for _, ext := range m.Extensions {
		files = append(files, base+"."+ext)
	}

	for _, filename := range files {
		if f, err := m.Root.Open(filename); err == nil {
			return f
		}
	}
	return nil
}

func notFound(w http.ResponseWriter, r *http.Request, msg string) {
	log.Println(msg)
	http.NotFound(w, r)
}
