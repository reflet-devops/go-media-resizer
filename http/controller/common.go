package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/hash"
	"github.com/reflet-devops/go-media-resizer/transform"
	"github.com/reflet-devops/go-media-resizer/types"
)

var TimeLocationGMT *time.Location

func init() {
	var err error
	TimeLocationGMT, err = time.LoadLocation("GMT")
	if err != nil {
		panic(err)
	}
}

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

func SendStream(ctx *context.Context, c echo.Context, opts *types.ResizeOption, content *bytes.Buffer) error {
	defer func() {
		resetBuffer(ctx, content)
	}()
	vary := []string{echo.HeaderAccept}
	acceptHeaderValue := c.Request().Header.Get(echo.HeaderAccept)
	DetectFormatFromHeaderAccept(acceptHeaderValue, opts)

	if opts.NeedTransform() {
		var errTransform error
		errTransform = transform.Transform(content, opts)
		if errTransform != nil {
			ctx.Logger.Error(fmt.Sprintf("failed to read data %s: %v", opts.Source, errTransform))
			return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to transform image %s", opts.Source))
		}
	}

	contentHash, _ := hash.GenerateMD5FromBytes(content.Bytes())

	c.Response().Header().Add(echo.HeaderContentLength, strconv.Itoa(content.Len()))
	c.Response().Header().Add("Date", time.Now().In(TimeLocationGMT).Format(time.RFC1123))
	c.Response().Header().Add("ETag", contentHash)
	c.Response().Header().Add(echo.HeaderVary, strings.Join(vary, ", "))
	if opts.HasTags() {
		c.Response().Header().Add("Cache-Tag", opts.TagsString())
	}

	for k, v := range opts.Headers {
		c.Response().Header().Add(k, v)
	}

	return c.Stream(http.StatusOK, types.GetMimeType(opts.Format), content)
}

func resetBuffer(ctx *context.Context, content *bytes.Buffer) {
	if content != nil {
		content.Reset()
		ctx.BufferPool.Put(content)
	}
}
