package markdown_server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

const (
	index       = "# index"
	about       = "## about"
	secondIndex = "# second index"
	secondAbout = "## second about"

	rIndex       = "<h1>index</h1>"
	rAbout       = "<h2>about</h2>"
	rSecondIndex = "<h1>second index</h1>"
	rSecondAbout = "<h2>second about</h2>"
)

func TestServer(t *testing.T) {
	handler := &MarkdownServer{
		Root:       http.Dir(setup(t)),
		IndexFile:  "index.md",
		Extensions: []string{"md"},
	}

	ts := httptest.NewServer(handler)
	defer ts.Close()

	assertContains(t, request(t, ts.URL), rIndex)
	assertContains(t, request(t, ts.URL+"/"), rIndex)
	assertContains(t, request(t, ts.URL+"/index.md"), rIndex)
	assertContains(t, request(t, ts.URL+"/about.md"), rAbout)

	assertContains(t, request(t, ts.URL+"/foo"), rSecondIndex)
	assertContains(t, request(t, ts.URL+"/foo/"), rSecondIndex)
	assertContains(t, request(t, ts.URL+"/foo/index.md"), rSecondIndex)
	assertContains(t, request(t, ts.URL+"/foo/about.md"), rSecondAbout)
}

func setup(t *testing.T) string {
	t.Helper()
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join(dir, "index.md"), []byte(index), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join(dir, "about.md"), []byte(about), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(filepath.Join(dir, "foo"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join(dir, "foo/index.md"), []byte(secondIndex), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join(dir, "foo/about.md"), []byte(secondAbout), 0644)
	if err != nil {
		t.Fatal(err)
	}

	return dir
}

func request(t *testing.T, url string) []byte {
	t.Helper()
	res, err := http.Get(url)
	if err != nil {
		t.Fatal("unable to make request", err)
	}
	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("unable to read response", err)
	}
	return page
}

func assertContains(t *testing.T, haystack []byte, needle string) {
	t.Helper()
	if !bytes.Contains(haystack, []byte(needle)) {
		t.Errorf("Mismatch: expect %s to contain %s", string(haystack), needle)
	}
}
