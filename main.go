package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuffalo/packr/v2"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/skip2/go-qrcode"
)

var (
	root = ""
)

// Directory contain files for html template
type Directory struct {
	Files []string
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if filepath.Ext(path) != ".cia" {
			return nil
		}
		filename, _ := filepath.Rel(root, path)
		*files = append(*files, filename)
		return nil
	}
}

func index(htmlTemplate string) func(w http.ResponseWriter, req *http.Request) {
	// Prepare template
	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, req *http.Request) {
		var files []string
		err = filepath.Walk(root, visit(&files))
		if err != nil {
			log.Fatal(err)
		}
		data := Directory{
			Files: files,
		}
		tmpl.Execute(w, data)
	}
}

func fileQRcode(w http.ResponseWriter, req *http.Request) {
	// fileName := chi.URLParam(req, "fileName")
	fileName := chi.URLParam(req, "*")

	var png []byte
	url := fmt.Sprintf("http://%s/files/%s", req.Host, url.PathEscape(fileName))
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(png)))
	w.Write(png)
}

// Midlleware to remove /files
func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		log.Fatal("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := neuter(http.StripPrefix(pathPrefix, http.FileServer(root)))
		fs.ServeHTTP(w, r)
	})
}

func main() {
	rootPtr := flag.String("path", "", "Path to shared directory")
	serverPort := flag.Int("port", 8090, "Web server port")
	flag.Parse()
	root = *rootPtr
	if len(root) == 0 {
		log.Fatal("Path have to be set")
	}
	root, _ = homedir.Expand(root)
	root, _ = filepath.Abs(root)

	// https://github.com/gobuffalo/packr/tree/master/v2
	box := packr.New("html templates", "./templates")
	html, err := box.FindString("index.html")
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", index(html))
	r.Get("/img/*", fileQRcode)

	fileServer(r, "/files", http.Dir(root))

	http.ListenAndServe(fmt.Sprintf(":%d", *serverPort), r)
}
