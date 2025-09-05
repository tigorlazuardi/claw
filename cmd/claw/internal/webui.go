package internal

import (
	"embed"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/olivere/vite"
)

var indexTemplate = template.Must(template.New("").Parse( /*html*/
	`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{if .WebUI.Title}}{{.WebUI.Title}}{{else}}Claw{{end}}</title>
	{{- .Vite.Tags -}}
</head>
<body>
	<div id="app"></div>
</body>
`))

func CreateWebuiHandler(fragment *vite.Fragment, logger *slog.Logger) http.Handler {
	var files fs.FS
	if cfg.Server.WebUI.Path == "" {
		files = embeddedWebUIFS
	} else {
		logger.Info("serving web UI from filesystem", "path", cfg.Server.WebUI.Path)
		files = os.DirFS(cfg.Server.WebUI.Path)
	}
	fileServer := http.FileServer(http.FS(files))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/assets") {
			// Serve static assets
			fileServer.ServeHTTP(w, r)
			return
		}
		if err := indexTemplate.Execute(w, map[string]any{
			"Vite":  fragment,
			"WebUI": cfg.Server.WebUI,
		}); err != nil {
			logger.Error("failed to execute http template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})
}

//go:embed all:webui
var webuiFiles embed.FS
var embeddedWebUIFS, _ = fs.Sub(webuiFiles, "webui")

// CreateViteFragment creates a Vite fragment from the given directory.
// If dir is empty, it uses the embedded webui files.
func CreateViteFragment() (*vite.Fragment, error) {
	var files fs.FS
	if cfg.Server.WebUI.Path == "" {
		files = embeddedWebUIFS
	} else {
		files = os.DirFS(cfg.Server.WebUI.Path)
	}
	fragment, err := vite.HTMLFragment(vite.Config{
		FS:    files,
		IsDev: cfg.Server.WebUI.DevMode,
	})
	if err != nil {
		return nil, err
	}
	if cfg.Server.WebUI.DevMode {
		slog.Info("running in development mode; make sure the Vite dev server is running before accessing the web UI")
	}
	return fragment, nil
}
