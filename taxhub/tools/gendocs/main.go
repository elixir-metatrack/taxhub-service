package main

import (
	"os"
	"text/template"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>{{.Title}}</title>
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    const spec = {{.Spec}};
    SwaggerUIBundle({
      spec: spec,
      dom_id: "#swagger-ui",
      deepLinking: true,
      presets: [SwaggerUIBundle.presets.apis],
    });
  </script>
</body>
</html>
`

func main() {
	specPath := "docs/swagger.json"
	output := "docs/index.html"
	title := "TaxHub API"

	spec, err := os.ReadFile(specPath)
	if err != nil {
		panic(err)
	}

	tmpl := template.Must(template.New("html").Parse(htmlTemplate))

	f, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, map[string]string{
		"Title": title,
		"Spec":  string(spec),
	}); err != nil {
		panic(err)
	}
}
