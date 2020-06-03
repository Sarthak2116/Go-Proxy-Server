package main

import (
	"net/http/httputil"
	"net/url"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	proxy1 := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:1324",
	})
	e.GET("/auth", echo.WrapHandler(proxy1))
	proxy2 := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:1325",
	})
	e.GET("/service/name", echo.WrapHandler(proxy2))
	e.GET("/user/profile", echo.WrapHandler(proxy2))
	e.Logger.Fatal(e.Start(":1323"))

}
