package music

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gosimple/slug"
)

type Album struct {
	Title  string
	Year   string
	Slug   string
	Cover  string
	Path   string
	Tracks []Track
}

type Track struct {
	Number string
	Title  string
	Path   string
}

func LoadAlbums() ([]Album, error) {

	const musicPath = "./static/music"

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
			Cover: "music/" + ap.Name() + "/cover.jpg",
			Path:  musicPath + "/" + ap.Name(),
		}

		albumTracksPath, err := ioutil.ReadDir(musicPath + "/" + ap.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read `%s` album tracks",
				musicPath+"/"+ap.Name())
		}

		for _, atp := range albumTracksPath {

			if atp.Name() == "cover.jpg" {
				continue
			}

			album.Tracks = append(album.Tracks, Track{
				Number: (atp.Name()[1:])[:2],
				Title:  atp.Name()[11:],
				Path:   "music/" + ap.Name() + "/" + atp.Name(),
			})
		}

		albums = append(albums, album)
	}

	return albums, nil
}