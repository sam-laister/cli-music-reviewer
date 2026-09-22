package album_picker

import "cli-music-reviewer/models/dtos"

type AlbumsLoadedMsg struct {
	albums []dtos.SavedAlbumDTO
	err    error
}

type ArtworkRenderedMsg struct {
	requestID int
	url       string
	rendered  string
	err       error
}

type AlbumSelectedMsg struct {
	Album dtos.AlbumDTO
}
