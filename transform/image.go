package transform

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"os/exec"
	"slices"
	"time"

	"github.com/disintegration/imaging"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/reflet-devops/go-media-resizer/types"
)

func Transform(file *bytes.Buffer, opts *types.ResizeOption) error {
	if !opts.NeedTransform() {
		return nil
	}

	img, _, errDecode := image.Decode(file)
	if errDecode != nil {
		return fmt.Errorf("failed to decode image %s: %w", opts.Source, errDecode)
	}

	if opts.NeedResize() {
		if opts.Fit == types.TypeFitCrop && (opts.Width == 0 || opts.Height == 0) {
			return fmt.Errorf("cannot crop without width and height")
		}
		img = Resize(img, opts)
	}

	if opts.NeedAdjust() {
		img = Adjust(img, opts)
	}

	errFormat := Format(file, img, opts)
	if errFormat != nil {
		return fmt.Errorf("failed to format image %s: %w", opts.Source, errFormat)
	}
	return nil
}

func Resize(img image.Image, opts *types.ResizeOption) image.Image {
	var imgResize *image.NRGBA

	switch opts.Fit {
	case types.TypeFitCrop:
		imgResize = imaging.CropCenter(img, opts.Width, opts.Height)
	case types.TypeFitScaleDown:
		imgResize = imaging.Fit(img, opts.Width, opts.Height, imaging.Lanczos)
	default:
		imgResize = imaging.Resize(img, opts.Width, opts.Height, imaging.Lanczos)
	}

	return imgResize
}

func Adjust(img image.Image, opts *types.ResizeOption) image.Image {
	for _, fn := range adjustFnList {
		img = fn(img, opts)
	}
	return img
}

func encodeAVIFExternal(buffer *bytes.Buffer, img image.Image, quality int) error {
	tmpDir := "/dev/shm"
	tmpInput := fmt.Sprintf("%s/avif_input_%d.png", tmpDir, time.Now().UnixNano())
	tmpOutput := fmt.Sprintf("%s/avif_output_%d.avif", tmpDir, time.Now().UnixNano())

	defer func() {
		_ = os.Remove(tmpInput)
		_ = os.Remove(tmpOutput)
	}()

	inputFile, err := os.Create(tmpInput)
	if err != nil {
		return fmt.Errorf("failed to create temp input file: %w", err)
	}

	err = png.Encode(inputFile, img)
	_ = inputFile.Close()
	if err != nil {
		return fmt.Errorf("failed to encode PNG to temp file: %w", err)
	}

	cmd := exec.Command("avifenc",
		"--speed", "10",
		"--jobs", "all",
		tmpInput,
		tmpOutput,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("avifenc failed: %w (stderr: %s)", err, stderr.String())
	}

	avifData, err := os.ReadFile(tmpOutput)
	if err != nil {
		return fmt.Errorf("failed to read AVIF output: %w", err)
	}

	buffer.Write(avifData)

	return nil
}

func Format(buffer *bytes.Buffer, img image.Image, opts *types.ResizeOption) error {
	var errFormat error

	if slices.Contains([]string{types.TypeAVIF, types.TypeWEBP}, opts.Format) {
		if opts.Format == types.TypeAVIF {
			quality := 60
			if opts.Quality != 0 {
				quality = opts.Quality
			}
			errFormat = encodeAVIFExternal(buffer, img, quality)
		} else if opts.Format == types.TypeWEBP {
			errFormat = webp.Encode(buffer, img, nil)
		}

	} else if slices.Contains([]string{types.TypeJPEG, types.TypePNG}, opts.Format) {
		format, errFindFormat := imaging.FormatFromExtension(opts.OriginFormat)
		if errFindFormat != nil {
			return fmt.Errorf("failed to find format from %s: %w", opts.Source, errFindFormat)
		}

		if opts.OriginFormat == types.TypeJPEG && opts.Quality != 0 {
			optsEncode := imaging.JPEGQuality(opts.Quality)
			errFormat = imaging.Encode(buffer, img, format, optsEncode)
		} else {
			errFormat = imaging.Encode(buffer, img, format)
		}
	} else {
		return fmt.Errorf("unsupported format: %s", opts.Format)
	}

	return errFormat
}
