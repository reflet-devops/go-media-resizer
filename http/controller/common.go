package controller

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/format"
	"github.com/reflet-devops/go-media-resizer/hash"
	"github.com/reflet-devops/go-media-resizer/resize"
	"github.com/reflet-devops/go-media-resizer/types"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

var TimeLocationGMT *time.Location

func init() {
	var err error
	TimeLocationGMT, err = time.LoadLocation("GMT")
	if err != nil {
		panic(err)
	}
}

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

func SendStream(ctx *context.Context, c echo.Context, opts *types.ResizeOption, content io.Reader) error {
	vary := []string{echo.HeaderAcceptEncoding}
	DetectFormatFromHeaderAccept(c.Request().Header.Get(echo.HeaderAccept), opts)
	if opts.NeedResize() {
		var err error
		content, err = resize.Resize(content, opts)
		if err != nil {
			ctx.Logger.Error(fmt.Sprintf("failed to resize(%v): %v", opts, err))
			return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to resize %s", opts.Source))
		}
	}

	if opts.OriginFormat != opts.Format {
		var err error
		content, err = format.Format(content, opts)
		if err != nil {
			ctx.Logger.Error(fmt.Sprintf("failed to format(%s from %s) %s: %v", opts.Format, opts.OriginFormat, opts.Source, err))
			return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to format %s", opts.Source))
		}
		vary = append(vary, echo.HeaderAccept)
	}
	data, err := io.ReadAll(content)
	if err != nil {
		ctx.Logger.Error(fmt.Sprintf("failed to read data %s: %v", opts.Source, err))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to read data %s", opts.Source))
	}
	contentHash, _ := hash.GenerateMD5(bytes.NewReader(data))

	c.Response().Header().Add(echo.HeaderContentLength, strconv.Itoa(len(data)))
	c.Response().Header().Add("Date", time.Now().In(TimeLocationGMT).Format(time.RFC1123))
	c.Response().Header().Add("ETag", contentHash)
	c.Response().Header().Add(echo.HeaderVary, strings.Join(vary, ", "))

	return c.Blob(http.StatusOK, types.GetMimeType(opts.Format), data)
}
