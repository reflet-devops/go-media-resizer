package controller

import (
	"bytes"
	"fmt"
	buildinHttp "net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/route"
	"github.com/reflet-devops/go-media-resizer/http/urltools"
	"github.com/reflet-devops/go-media-resizer/logger"
	"github.com/reflet-devops/go-media-resizer/mapstructure"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/valyala/fasthttp"
)

func GetMediaCGI(ctx *context.Context) func(c echo.Context) error {
	return func(c echo.Context) error {
		source := c.Param("source")
		opts := ctx.OptsResizePool.Get().(*types.ResizeOption)
		opts.ResetToDefaults(&ctx.Config.ResizeCGI.DefaultResizeOpts)
		opts.Source = source

		fileExtension := urltools.GetExtension(source)
		fileType := types.GetType(fileExtension)
		opts.OriginFormat = fileType
		optMap := parseOption(c.Param("options"))

		fileTypeIsValid := types.ValidateType(fileType, ctx.Config.AcceptTypeFiles)
		if !fileTypeIsValid {
			ctx.Logger.Error(fmt.Sprintf("GetMediaCGI: file type not accepted: %s", source), addLogAttr(c)...)
			return c.String(buildinHttp.StatusInternalServerError, "file type not accepted")
		}

		err := mapstructure.Decode(optMap, &opts)
		if err != nil {
			return c.String(buildinHttp.StatusInternalServerError, err.Error())
		}

		buffer := ctx.BufferPool.Get().(*bytes.Buffer)
		projectIdHeader, errFetch := fetchCGIResource(ctx, c.Request().Header.Get(echo.HeaderXRequestID), source, buffer)
		if errFetch != nil {
			resetBuffer(ctx, buffer)
			return c.String(buildinHttp.StatusInternalServerError, errFetch.Error())
		}

		opts.AddTag(types.GetTagSourcePathHash(types.FormatProjectPathHash(projectIdHeader, urltools.GetUri(opts.Source))))
		for k, v := range ctx.Config.Headers {
			opts.AddHeader(k, v)
		}
		return SendStream(ctx, c, opts, buffer)
	}
}

func parseOption(optsHeader string) map[string]interface{} {
	optRaw := strings.Split(optsHeader, ",")
	optMap := map[string]interface{}{}
	for _, optStr := range optRaw {
		optSplit := strings.Split(optStr, "=")
		if len(optSplit) == 2 {
			key := optSplit[0]
			value := optSplit[1]

			key = strings.Trim(key, " ")
			value = strings.Trim(value, " ")

			optMap[key] = value
		}
	}
	return optMap
}

func fetchCGIResource(ctx *context.Context, requestId string, source string, buffer *bytes.Buffer) (string, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		if req != nil {
			fasthttp.ReleaseRequest(req)
		}
		if resp != nil {
			fasthttp.ReleaseResponse(resp)
		}
	}()
	req.Header.SetMethod(buildinHttp.MethodGet)
	req.Header.Add(echo.HeaderXRequestID, requestId)
	req.SetRequestURI(source)
	ctx.Logger.Debug(fmt.Sprintf("fetchCGIResource: GET %s", source), logger.RequestIDKey, requestId)
	err := ctx.HttpClient.DoTimeout(req, resp, ctx.Config.RequestTimeout)

	if err != nil {
		return "", fmt.Errorf("fetchCGIResource: GET %s: error with request: %v", source, err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return "", fmt.Errorf("fetchCGIResource: GET %s: invalid status code status code: %d", source, resp.StatusCode())
	}
	_, err = buffer.Write(resp.Body())
	return string(resp.Header.Peek(route.ProjectIdHeader)), err
}
