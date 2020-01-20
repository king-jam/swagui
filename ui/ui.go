// Package ui holds all customization items for the Swagger UI release
package ui

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"text/template"

	_ "github.com/king-jam/swagui/ui/statik"

	"github.com/rakyll/statik/fs"
)

func UI() {
	dist, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	// Serve the contents over HTTP.
	//http.Handle("/", http.HandlerFunc(handler))
	http.Handle("/", http.StripPrefix("/", http.FileServer(swagui(dist))))
	http.ListenAndServe(":8080", nil)
}

type uiHandler struct {
	dist http.FileSystem
	afero.NewMemMapFs
	idx bytes.Buffer
}

func swagui(dist http.FileSystem) http.FileSystem {
	b := new(bytes.Buffer)

	genIndex(b)

	return &uiHandler{dist}
}

func (f *uiHandler) Open(name string) (http.File, error) {
	if name == "/index.html" {
		log.Print("My Special Handler")
	}
	return f.dist.Open(name)
}

const indexHtml = `<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="./swagger-ui.css" >
    <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }

      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }

      body
      {
        margin:0;
        background: #fafafa;
      }
    </style>
  </head>

  <body>
    <div id="swagger-ui"></div>

    <script src="./swagger-ui-bundle.js"> </script>
    <script src="./swagger-ui-standalone-preset.js"> </script>
    <script>
    window.onload = function() {
      // Begin Swagger UI call region
      const ui = SwaggerUIBundle({
        url: "{{.URL}}",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      })
      // End Swagger UI call region

      window.ui = ui
    }
  </script>
  </body>
</html>
`

func genIndex(w io.Writer) {
	t := template.New("index.html")
	t, _ = t.Parse(indexHtml)
	url := struct {
		URL string
	}{
		URL: "localhost:8080/swagger.json",
	}
	t.Execute(w, url)
}
