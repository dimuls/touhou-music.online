package music

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
)

const musicPath = "./static/music"

type Album struct {
	Title string `json:"title"`
	Year  string `json:"year"`
	Slug  string `json:"slug"`
	Cover string `json:"cover"`
	Path  string `json:"path"`
	Discs []Disc `json:"discs"`
}

type Disc struct {
	Number string  `json:"number"`
	Tracks []Track `json:"tracks"`
}

type Track struct {
	Number string `json:"number"`
	Title  string `json:"title"`
	Path   string `json:"path"`
}

func LoadAlbums() ([]Album, error) {

	albumsPaths, err := ioutil.ReadDir(musicPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read `%s` music path", musicPath)
	}

	var albums []Album

	for _, ap := range albumsPaths {

		if !ap.IsDir() {
			continue
		}

		i := strings.Index(ap.Name(), "[")
		if i == -1 {
			continue
		}

		title := ap.Name()[:i]
		year := (ap.Name()[i+1:])[:4]

		album := Album{
			Title: title,
			Year:  year,
			Slug:  slug.Make(year + " " + title),
			Cover: filepath.Join("music", ap.Name(), "cover.jpg"),
			Path:  filepath.Join(musicPath, ap.Name()),
		}

		discsPath, err := ioutil.ReadDir(album.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to read `%s` discs path: %v",
				musicPath+"/"+ap.Name(), err)
		}

		for _, dp := range discsPath {

			if !dp.IsDir() {
				if strings.Contains(dp.Name(), ".mp3") {
					// If we found mp3 file this means that we need to load
					// discs path itself as disc 1.

					disc, err := loadDisc(ap.Name(), "")
					if err != nil {
						return nil, err
					}

					disc.Number = "1"

					album.Discs = append(album.Discs, disc)

					break
				}
				continue
			}

			i := strings.Index(dp.Name(), "Disc")
			if i != 0 {
				continue
			}

			disc, err := loadDisc(ap.Name(), dp.Name())
			if err != nil {
				return nil, err
			}

			disc.Number = dp.Name()[5:]

			album.Discs = append(album.Discs, disc)
		}

		albums = append(albums, album)
	}

	return albums, nil
}

func loadDisc(albumPathName, discPathName string) (Disc, error) {

	discBasePath := filepath.Join(musicPath, albumPathName, discPathName)

	discTracksPath, err := ioutil.ReadDir(discBasePath)
	if err != nil {
		return Disc{}, fmt.Errorf("failed to read `%s` disc tracks path: %v",
			discBasePath, err)
	}

	var disc Disc

	for _, dtp := range discTracksPath {

		if dtp.Name() == "cover.jpg" {
			continue
		}

		fileName := dtp.Name()

		dashIndex := strings.Index(fileName, "–")
		if dashIndex == -1 {
			continue
		}

		extIndex := strings.Index(fileName, ".mp3")
		if extIndex == -1 {
			continue
		}

		number := (fileName[1:])[:2]
		title := fileName[dashIndex+4 : extIndex]

		disc.Tracks = append(disc.Tracks, Track{
			Number: number,
			Title:  title,
			Path:   filepath.Join("/static/music", albumPathName, discPathName, dtp.Name()),
		})
	}

	return disc, nil
}
