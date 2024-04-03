package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	fmt.Print("\x1B[2J\x1B[1;1H")
	fmt.Println("🌎 Iniciando servidor fake na porta \"8765\"...")

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{"*"},
	}))

	e.Any("*", handleRequest)

	e.Logger.Fatal(e.Start(":8765"))
}

func handleRequest(c echo.Context) error {
	method := c.Request().Method
	fullURI := c.Request().RequestURI
	parts := strings.Split(fullURI, "?")
	uri := parts[0]

	var query []string
	if len(parts) > 1 {
		query = strings.Split(parts[1], "&")
	}

	path := fmt.Sprintf("paths/%s.json", uri[1:])

	if content, err := os.ReadFile(path); err == nil {
		fmt.Printf("%-7s  ✅ 200     %s\n", method, uri)
		return c.String(http.StatusOK, string(content))
	}

	if len(query) > 0 {
		for _, item := range query {
			keyValue := strings.Split(item, "=")

			path := fmt.Sprintf("paths/%s/%s.json", uri[1:], keyValue[0])
			if content, err := os.ReadFile(path); err == nil {
				fmt.Printf("%-7s  ✅ 200     %s\n", method, uri)
				return c.String(http.StatusOK, string(content))
			}

			path = fmt.Sprintf("paths/%s/%s.json", uri[1:], keyValue[1])
			if content, err := os.ReadFile(path); err == nil {
				fmt.Printf("%-7s  ✅ 200     %s\n", method, uri)
				return c.String(http.StatusOK, string(content))
			}
		}
	}

	fmt.Printf("%-7s  ❌ 404     %s\n", method, fullURI)

	return c.String(http.StatusNotFound, "Not Found")
}
