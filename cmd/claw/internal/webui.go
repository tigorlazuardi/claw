package internal

import (
	"embed"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

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

type WebUIConfig struct {
	Fragment *vite.Fragment
	DistFS   fs.FS
	Logger   *slog.Logger
}

func CreateWebuiHandler(config WebUIConfig) http.Handler {
	fileServer := http.FileServer(http.FS(config.DistFS))
	mux := http.NewServeMux()
	if cfg.Server.WebUI.DevMode {
		mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))
	} else {
		mux.Handle("/assets", http.StripPrefix("/assets", fileServer))
	}
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := indexTemplate.Execute(w, map[string]any{
			"Vite":  config.Fragment,
			"WebUI": cfg.Server.WebUI,
		}); err != nil {
			config.Logger.Error("failed to execute http template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}))
	return mux
}

//go:embed all:webui
var webuiFiles embed.FS
var embeddedWebUIFS, _ = fs.Sub(webuiFiles, "webui")

// CreateViteFragment creates a Vite fragment from the given directory.
// If dir is empty, it uses the embedded webui files.
func CreateViteFragment() (*vite.Fragment, fs.FS, error) {
	var files fs.FS
	if cfg.Server.WebUI.DevMode {
		files = os.DirFS("./webui") // Relative to root working directory. Assuming this golang server runs from project root.
	} else if cfg.Server.WebUI.Path == "" {
		files = embeddedWebUIFS
	} else {
		files = os.DirFS(cfg.Server.WebUI.Path)
	}
	fragment, err := vite.HTMLFragment(vite.Config{
		FS:        files,
		IsDev:     cfg.Server.WebUI.DevMode,
		ViteEntry: "src/main.ts", // This is used in dev mode to load the Vite dev server script.
	})
	if err != nil {
		return nil, files, err
	}
	if cfg.Server.WebUI.DevMode {
		slog.Info("running in development mode; make sure the Vite dev server is running before accessing the web UI")
	}
	return fragment, files, nil
}
