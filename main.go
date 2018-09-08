package main

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"

	"github.com/dimuls/touhou-music.online/music"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	albums    []music.Album
	albumsMap map[string]music.Album
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))

	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("template/*.html")),
	}

	e.Static("/", "static")

	albums, err := music.LoadAlbums()
	if err != nil {
		e.Logger.Fatalf("Failed to load albums: %v", err)
	}

	useAlbums(albums)

	e.GET("/", indexHandler)
	e.GET("/:albumSlug", albumHandler)
	e.File("/google0c80a89b75247802.html", "static/google0c80a89b75247802.html")
	e.File("/sitemap.xml", "static/sitemap.xml")

	e.Logger.Fatal(e.Start(":80"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func useAlbums(as []music.Album) {
	albums = as
	albumsMap = map[string]music.Album{}
	for _, a := range albums {
		albumsMap[a.Slug] = a
	}
}

func indexHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index", albums)
}

type albumHandlerData struct {
	Album     music.Album
	AlbumJSON string
}

func albumHandler(c echo.Context) error {
	slug := c.Param("albumSlug")
	album, exists := albumsMap[slug]
	if !exists {
		return c.String(http.StatusNotFound, "404 Not found")
	}

	albumJSON, err := json.Marshal(album)
	if err != nil {
		c.Logger().Errorf("Failed to marshal album: %v", err)
		return c.String(http.StatusInternalServerError,
			"500 Internal server error")
	}

	return c.Render(http.StatusOK, "album", albumHandlerData{
		Album:     album,
		AlbumJSON: string(albumJSON),
	})
}
