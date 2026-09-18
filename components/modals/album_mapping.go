package modals

import (
	"cli-music-reviewer/models/dtos"
	"strings"
)

// albumArtists joins an album's artist names into a single display string.
func albumArtists(album dtos.AlbumDTO) string {
	names := make([]string, len(album.Artists))
	for i, artist := range album.Artists {
		names[i] = artist.Name
	}
	return strings.Join(names, ", ")
}

// albumYear extracts the leading year from a Spotify release_date value,
// which may be "YYYY", "YYYY-MM", or "YYYY-MM-DD" depending on precision.
func albumYear(album dtos.AlbumDTO) string {
	year, _, _ := strings.Cut(album.ReleaseDate, "-")
	return year
}

// albumCoverArt maps an album's widest-first Images array onto the app's
// small/medium/large cover art fields, tolerating fewer than 3 images.
func albumCoverArt(album dtos.AlbumDTO) (small, medium, large string) {
	images := album.Images
	if len(images) > 0 {
		large = images[0].URL
	}
	if len(images) > 1 {
		medium = images[1].URL
	}
	if len(images) > 2 {
		small = images[2].URL
	}
	return small, medium, large
}

// albumArtworkURL picks the best available image to render inline art from,
// preferring the medium size to keep downloads light.
func albumArtworkURL(album dtos.AlbumDTO) string {
	_, medium, large := albumCoverArt(album)
	if medium != "" {
		return medium
	}
	return large
}
