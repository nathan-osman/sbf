package server

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jpillora/sizestr"
)

type fileEntry struct {
	Name    string
	Size    string
	ModTime time.Time
}

var (
	//go:embed templates/*
	tmplFS embed.FS

	sanitizeRE = regexp.MustCompile(`[^\w\s-.]`)
)

func render(c *gin.Context, name string, data any) {
	t, err := template.ParseFS(tmplFS, "templates/base.html", name)
	if err != nil {
		panic(err)
	}
	if err := t.Execute(c.Writer, data); err != nil {
		panic(err)
	}
}

func sanitize(v string) string {
	return sanitizeRE.ReplaceAllString(v, "")
}

func (s *Server) adminIndex(c *gin.Context) {
	render(c, "templates/admin/index.html", nil)
}

func (s *Server) index(c *gin.Context) {
	render(c, "templates/index.html", gin.H{
		"Year": time.Now().Year(),
	})
}

func (s *Server) srvFolderGET(c *gin.Context) {
	n := sanitize(c.Param("name"))
	entries, err := os.ReadDir(filepath.Join(s.dir, n))
	if err != nil {
		panic(err)
	}
	files := []*fileEntry{}
	for _, entry := range entries {
		if !entry.IsDir() {
			i, err := entry.Info()
			if err != nil {
				panic(err)
			}
			files = append(files, &fileEntry{
				Name:    i.Name(),
				Size:    sizestr.ToString(i.Size()),
				ModTime: i.ModTime(),
			})
		}
	}
	render(c, "templates/srv/index.html", gin.H{
		"Name":  n,
		"Files": files,
	})
}

func (s *Server) srvFolderPOST(c *gin.Context) {
	n := sanitize(c.Param("name"))
	f, err := c.FormFile("file")
	if err != nil {
		panic(err)
	}
	p := filepath.Join(s.dir, n, sanitize(f.Filename))
	if err := c.SaveUploadedFile(f, p); err != nil {
		panic(err)
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/%s", n))
}
