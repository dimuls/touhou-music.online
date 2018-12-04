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
	Disks []Disk `json:"disks"`
}

type Disk struct {
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

		disksPath, err := ioutil.ReadDir(album.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to read `%s` disks path: %v",
				musicPath+"/"+ap.Name(), err)
		}

		for _, dp := range disksPath {

			if !dp.IsDir() {
				if strings.Contains(dp.Name(), ".mp3") {
					// If we found mp3 file this means that we need to load
					// disks path itself as disk 1.

					disk, err := loadDisk(ap.Name(), "")
					if err != nil {
						return nil, err
					}

					disk.Number = "1"

					album.Disks = append(album.Disks, disk)

					break
				}
				continue
			}

			i := strings.Index(dp.Name(), "Disc")
			if i != 0 {
				continue
			}

			disk, err := loadDisk(ap.Name(), dp.Name())
			if err != nil {
				return nil, err
			}

			disk.Number = dp.Name()[5:]

			album.Disks = append(album.Disks, disk)
		}

		albums = append(albums, album)
	}

	return albums, nil
}

func loadDisk(albumPathName, diskPathName string) (Disk, error) {

	diskBasePath := filepath.Join(musicPath, albumPathName, diskPathName)

	diskTracksPath, err := ioutil.ReadDir(diskBasePath)
	if err != nil {
		return Disk{}, fmt.Errorf("failed to read `%s` disk tracks path: %v",
			diskBasePath, err)
	}

	var disk Disk

	for _, dtp := range diskTracksPath {

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

		disk.Tracks = append(disk.Tracks, Track{
			Number: number,
			Title:  title,
			Path:   filepath.Join("/static/music", albumPathName, diskPathName, dtp.Name()),
		})
	}

	return disk, nil
}
