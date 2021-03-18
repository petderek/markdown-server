package main

import (
	"log"
	"net/http"
	"os"

	"github.com/petderek/dflag"
	. "github.com/petderek/markdown-server"
)

var flags = struct {
	Addr string
	Root string
}{}

func main() {
	flags.Addr = ":80"
	if defaultDir, err := os.Getwd(); err != nil {
		flags.Root = defaultDir
	}

	if err := dflag.Parse(&flags); err != nil {
		log.Fatal("error parsing: ", err)
	}

	md := &MarkdownServer{
		Root: http.Dir(flags.Root),
	}

	log.Fatal(http.ListenAndServe(flags.Addr, md))
}
