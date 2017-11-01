package main

import (
	"strings"
	"flag"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type LinkHeader struct {
	headers []string
}

func NewLinkHeader(headers string) *LinkHeader {
	return &LinkHeader{strings.Split(headers, ",")}
}

func (l *LinkHeader) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		url := c.Request().URL.Path
		if url == "/" {
			for _, file := range(l.headers) {
				params=[]string{
					"<"+ file +">",
					"rel=preload",
					defineLinkType(file),
				}
				c.Response().Header().Add("link", strings.Join(params, "; "))
			}
		}
		return next(c)
	}
}

func defineLinkType(file string) {
	types:=map[string]string{
		"js": "script",
		"css": "style",
		"jpg": "image",
		"gif": "image",
		"png": "image",
	}
}

func main() {
	port := flag.String("port", "3003", "port number")
	folder := flag.String("folder", "public", "folder")
	headers := flag.String("indexHeaders", "", "index headers separated by , (comma)")
	flag.Parse()

	e := echo.New()
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
	  Level: 8,
	}))

	e.Static("/", *folder)

	linkHeader := NewLinkHeader(*headers)
	e.Pre(linkHeader.Middleware)

	e.Logger.Fatal(e.Start(":" + *port))
}
