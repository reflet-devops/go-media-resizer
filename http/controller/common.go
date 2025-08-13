package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

var HTTPErrorFileTypeNotAccepted = echo.NewHTTPError(http.StatusForbidden, "file type not accepted")
