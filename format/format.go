package format

import (
	"bytes"
	"fmt"
	"github.com/Kagami/go-avif"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/reflet-devops/go-media-resizer/types"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"slices"
)

func Format(file io.Reader, opts *types.ResizeOption) (io.Reader, error) {
	var errFormat error
	w := &bytes.Buffer{}

	if slices.Contains([]string{types.TypeAVIF, types.TypeWEBP}, opts.Format) {
		img, _, errDecode := image.Decode(file)
		if errDecode != nil {
			return nil, fmt.Errorf("failed to decode %s to format: %v", opts.Source, errDecode)
		}
		if opts.Format == types.TypeAVIF {
			errFormat = avif.Encode(w, img, nil)
		} else if opts.Format == types.TypeWEBP {
			errFormat = webp.Encode(w, img, nil)
		}

		return w, errFormat
	}

	return file, nil
}
