package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

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

	e.Pre(middleware.NonWWWRedirect())

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))

	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("template/*.html")),
	}

	e.Static("/static", "static")

	albums, err := music.LoadAlbums()
	if err != nil {
		e.Logger.Fatalf("Failed to load albums: %v", err)
	}

	useAlbums(albums)

	e.File("/google0c80a89b75247802.html", "static/google0c80a89b75247802.html")
	e.File("/sitemap.xml", "static/sitemap.xml")

	e.GET("/", indexHandler)
	e.GET("/:language", indexHandler)
	e.GET("/:language/:albumSlug", albumHandler)

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

var indexPageL10n = map[string]map[string]string{
	"en": {
		"title":       "Touhou music online",
		"description": "Listen touhou music online",
	},
	"ru": {
		"title":       "Touhou музыка онлайн",
		"description": "Слушать touhou музыку онлайн",
	},
}

type indexPageData struct {
	L10n      map[string]string
	Language  string
	Languages []string
	Albums    []music.Album
}

func indexHandler(c echo.Context) error {
	lang := c.Param("language")

	if lang == "" {
		// Redirect to canonical paths.
		lang = languageFromHeader(c)
		if _, exists := indexPageL10n[lang]; exists {
			// Redirect to localized path.
			return c.Redirect(http.StatusFound, "/"+lang)
		} else {
			// Redirect to english page if language is not supported.
			return c.Redirect(http.StatusFound, "/en")
		}

	} else if _, exists := albumsMap[lang]; exists {
		// If lang is album slug we need to keep old album URL working.
		albumSlug := lang
		lang = languageFromHeader(c)
		if _, exists := albumPageL10n[lang]; !exists {
			lang = "en"
		}
		return c.Redirect(http.StatusFound,
			fmt.Sprintf("/%s/%s", lang, albumSlug))
	}

	l10n, exists := indexPageL10n[lang]
	if !exists {
		// Redirect to english page if language is not supported.
		return c.Redirect(http.StatusFound, "/en")
	}

	var langs []string
	for l := range indexPageL10n {
		langs = append(langs, l)
	}

	return c.Render(http.StatusOK, "index", indexPageData{
		L10n:      l10n,
		Language:  lang,
		Languages: langs,
		Albums:    albums,
	})
}

var albumPageL10n = map[string]map[string]string{
	"en": {
		"description": "Touhou music album, year %s, %s",
		"back":        "Back",
		"disc":        "Disc",
	},
	"ru": {
		"description": "Альбом touhou музыки, год %s, %s",
		"back":        "Назад",
		"disc":        "Диск",
	},
}

type albumHandlerData struct {
	L10n        map[string]string
	Description string
	Language    string
	Languages   []string
	Album       music.Album
	AlbumJSON   string
}

func albumHandler(c echo.Context) error {
	slug := c.Param("albumSlug")
	album, exists := albumsMap[slug]
	if !exists {
		return c.String(http.StatusNotFound, "404 Not found")
	}

	lang := c.Param("language")
	l10n, exists := albumPageL10n[lang]
	if !exists {
		// Redirect to english page if language is not supported.
		return c.Redirect(http.StatusFound, "/en/"+slug)
	}

	albumJSON, err := json.Marshal(album)
	if err != nil {
		c.Logger().Errorf("Failed to marshal album: %v", err)
		return c.String(http.StatusInternalServerError,
			"500 Internal server error")
	}

	var langs []string
	for l := range indexPageL10n {
		langs = append(langs, l)
	}

	return c.Render(http.StatusOK, "album", albumHandlerData{
		L10n:        l10n,
		Description: fmt.Sprintf(l10n["description"], album.Year, album.Title),
		Language:    lang,
		Languages:   langs,
		Album:       album,
		AlbumJSON:   string(albumJSON),
	})
}

func languageFromHeader(c echo.Context) string {
	al := c.Request().Header.Get("Accept-Language")
	ls := strings.Split(al, ",")
	lang := "en"
LOOP:
	for _, lq := range ls {
		l := strings.Split(lq, "-")
		switch l[0] {
		case "en":
			break LOOP
		case "ru":
			lang = "ru"
			break LOOP
		}
	}
	return lang
}
