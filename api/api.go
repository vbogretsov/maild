package api

import (
	"fmt"
	// "io"
	"net/http"
	// "strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	// "github.com/vbogretsov/jsonerr"
	"github.com/vbogretsov/go-validation"

	"github.com/vbogretsov/maild/app"
	"github.com/vbogretsov/maild/model"
)

func New(a *app.App) (*echo.Echo, error) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if _, ok := err.(validation.Error); ok {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusInternalServerError, err)
		}
	}

	e.POST("/v1/sendmail", func(c echo.Context) error {
		req := model.Request{}

		if err := c.Bind(&req); err != nil {
			return err
		}

		return a.SendMail(req)
	})

	return e, nil
}
