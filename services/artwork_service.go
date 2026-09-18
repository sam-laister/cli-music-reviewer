package services

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"

	"github.com/BourgeoisBear/rasterm"
)

// ErrArtworkUnsupported indicates the current terminal has no known inline
// image protocol support (Kitty, iTerm2/WezTerm/Rio).
var ErrArtworkUnsupported = errors.New("terminal does not support inline images")

// kittyDeleteAllImages removes every placement the Kitty graphics protocol
// currently has on screen (and frees their image data). The Kitty protocol
// treats placed images as persistent overlays independent of normal text
// redraws, so a new placement never implicitly replaces an old one — without
// this, scrolling through albums (or navigating away from the picture
// entirely) leaves stale artwork stacked on screen.
const kittyDeleteAllImages = "\x1b_Ga=d,d=A\x1b\\"

type ArtworkService interface {
	// Supported reports whether the current terminal can render inline images.
	Supported() bool
	// Render downloads the image at url and returns the ready-to-print
	// terminal escape sequence sized to cols x rows character cells.
	Render(url string, cols, rows int) (string, error)
	// ClearImages returns the terminal control sequence needed to remove any
	// previously placed inline images, or "" if the current terminal's
	// protocol doesn't need explicit clearing.
	ClearImages() string
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
		buf.WriteString(kittyDeleteAllImages)
		err = writeKittyImage(&buf, img, cols, rows)
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

func (s *ArtworkServiceImpl) ClearImages() string {
	if rasterm.IsKittyCapable() {
		return kittyDeleteAllImages
	}
	return ""
}

const kittyChunkSize = 4096

// writeKittyImage transmits img using the Kitty graphics protocol with C=1
// ("do not move the cursor"). rasterm.KittyWriteImage doesn't expose that
// flag, and without it Kitty moves the terminal cursor to after the placed
// image — which desyncs every line Bubble Tea prints afterward in the same
// frame (it advances line-to-line with relative \r\n), producing corrupted,
// duplicated-looking redraws whenever anything is printed following the
// image in the same View() output.
func writeKittyImage(w io.Writer, img image.Image, cols, rows int) error {
	var pngBuf bytes.Buffer
	if err := png.Encode(&pngBuf, img); err != nil {
		return err
	}
	data := base64.StdEncoding.EncodeToString(pngBuf.Bytes())

	opts := rasterm.KittyImgOpts{DstCols: uint32(cols), DstRows: uint32(rows)}

	for offset, first := 0, true; offset < len(data) || first; first = false {
		end := min(offset+kittyChunkSize, len(data))
		chunk := data[offset:end]
		offset = end

		more := 0
		if offset < len(data) {
			more = 1
		}

		var header string
		if first {
			header = opts.ToHeader("a=T", "f=100", "t=d", "C=1", fmt.Sprintf("m=%d", more))
		} else {
			header = fmt.Sprintf("\x1b_Gm=%d;", more)
		}

		if _, err := fmt.Fprint(w, header, chunk, "\x1b\\"); err != nil {
			return err
		}
	}

	return nil
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
