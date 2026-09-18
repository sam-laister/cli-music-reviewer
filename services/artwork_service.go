package services

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/BourgeoisBear/rasterm"
)

// ErrArtworkUnsupported indicates the current terminal has no known inline
// image protocol support (Kitty, iTerm2/WezTerm/Rio).
var ErrArtworkUnsupported = errors.New("terminal does not support inline images")

type ArtworkService interface {
	// Supported reports whether the current terminal can render inline images.
	Supported() bool
	// Render downloads the image at url and returns the ready-to-print
	// terminal escape sequence sized to cols x rows character cells.
	Render(url string, cols, rows int) (string, error)
}

type ArtworkServiceImpl struct{}

func (s *ArtworkServiceImpl) Supported() bool {
	return rasterm.IsKittyCapable() || rasterm.IsItermCapable()
}

func (s *ArtworkServiceImpl) Render(url string, cols, rows int) (string, error) {
	img, err := s.download(url)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	switch {
	case rasterm.IsKittyCapable():
		err = rasterm.KittyWriteImage(&buf, img, rasterm.KittyImgOpts{
			DstCols: uint32(cols),
			DstRows: uint32(rows),
		})
	case rasterm.IsItermCapable():
		err = rasterm.ItermWriteImageWithOptions(&buf, img, rasterm.ItermImgOpts{
			Width:         fmt.Sprintf("%d", cols),
			Height:        fmt.Sprintf("%d", rows),
			DisplayInline: true,
		})
	default:
		return "", ErrArtworkUnsupported
	}

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *ArtworkServiceImpl) download(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("artwork download failed: %s", resp.Status)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func NewArtworkService() *ArtworkServiceImpl {
	return &ArtworkServiceImpl{}
}

var _ ArtworkService = (*ArtworkServiceImpl)(nil)
