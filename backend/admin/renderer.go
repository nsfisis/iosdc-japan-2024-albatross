package admin

import (
	"embed"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	//go:embed templates
	templatesFS embed.FS
	//go:embed assets
	assetsFS embed.FS
)

type Renderer struct {
	templates map[string]*template.Template
}

func NewRenderer() *Renderer {
	return &Renderer{
		templates: make(map[string]*template.Template),
	}
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	tmpl, ok := r.templates[name]
	if !ok {
		t, err := template.ParseFS(templatesFS, "templates/base.html", "templates/"+name+".html")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		r.templates[name] = t
		tmpl = t
	}
	return tmpl.ExecuteTemplate(w, name+".html", data)
}

func newAssetsMiddleware() echo.MiddlewareFunc {
	return middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "/assets",
		Filesystem: http.FS(assetsFS),
		IgnoreBase: true,
	})
}
