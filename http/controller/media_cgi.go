package controller

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/urltools"
	"github.com/reflet-devops/go-media-resizer/mapstructure"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/valyala/fasthttp"
	buildinHttp "net/http"
	"strings"
)

func GetMediaCGI(ctx *context.Context) func(c echo.Context) error {
	return func(c echo.Context) error {
		source := c.Param("source")
		opts := &types.ResizeOption{Source: source, Headers: types.Headers{}}

		fileExtension := urltools.GetExtension(source)
		fileType := types.GetType(fileExtension)
		opts.OriginFormat = fileType

		optMap := parseOption(c.Param("options"))

		fileTypeIsValid := types.ValidateType(fileType, ctx.Config.AcceptTypeFiles)
		if !fileTypeIsValid {
			ctx.Logger.Error(fmt.Sprintf("GetMediaCGI: file type not accepted: %s", source))
			return HTTPErrorFileTypeNotAccepted
		}

		err := mapstructure.Decode(optMap, opts)
		if err != nil {
			return c.JSON(buildinHttp.StatusInternalServerError, err.Error())
		}
		if opts.Format == "" {
			opts.Format = opts.OriginFormat
		}

		resource, err := fetchCGIResource(ctx, source)
		if err != nil {
			return c.JSON(buildinHttp.StatusInternalServerError, err.Error())
		}
		buffer := bytes.NewBuffer(resource)

		opts.AddTag(types.GetTagSourcePathHash(urltools.GetUri(opts.Source)))

		for k, v := range ctx.Config.Headers {
			opts.Headers[k] = v
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

func fetchCGIResource(ctx *context.Context, source string) ([]byte, error) {
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
	req.SetRequestURI(source)

	ctx.Logger.Debug(fmt.Sprintf("fetchCGIResource: GET %s", source))
	err := ctx.HttpClient.DoTimeout(req, resp, ctx.Config.RequestTimeout)

	if err != nil {
		return nil, fmt.Errorf("fetchCGIResource: GET %s: error with request: %v", source, err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, fmt.Errorf("fetchCGIResource: GET %s: invalid status code status code: %d", source, resp.StatusCode())
	}
	return resp.Body(), nil
}
