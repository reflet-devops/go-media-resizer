package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/types"
	"io"
	"net/http"
	"slices"
	"strings"
)

var HTTPErrorFileTypeNotAccepted = echo.NewHTTPError(http.StatusForbidden, "file type not accepted")

func DetectFormatFromHeaderAccept(acceptHeaderValue string, opts *types.ResizeOption) {
	acceptedFormat := strings.Split(acceptHeaderValue, ",")

	if slices.Contains([]string{types.TypeFormatAuto, types.TypeAVIF}, opts.Format) && slices.Contains(acceptedFormat, types.MimeTypeAVIF) {
		opts.Format = types.TypeAVIF
		return
	} else if slices.Contains([]string{types.TypeFormatAuto, types.TypeWEBP, types.TypeAVIF}, opts.Format) && slices.Contains(acceptedFormat, types.MimeTypeWEBP) {
		opts.Format = types.TypeWEBP
		return
	}

	opts.Format = opts.OriginFormat
}

func SendStream(c echo.Context, opts *types.ResizeOption, content io.Reader) error {
	contentType := types.GetMimeType(opts.Format)
	return c.Stream(http.StatusOK, contentType, content)
}
