package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/buaazp/fasthttprouter"
	fhm "github.com/nlapshin/fasthttpmiddleware"
	"github.com/valyala/fasthttp"
)

const PORT = ":8080"

func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Welcome!\n")
}

func Hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
}

func HelloAuth(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello, %s! Access Token = %s.", ctx.UserValue("name"), ctx.UserValue("access-token"))
}

func main() {
	router := fasthttprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	router.GET("/auth/:name", HelloAuth)

	middleware := fhm.New()
	middleware.Use(loggerMiddleware)
	middleware.Use(authMiddleware)

	fmt.Println("Server start on " + PORT + " address")

	log.Fatal(fasthttp.ListenAndServe(PORT, middleware.Handler(router.Handler)))
}

func loggerMiddleware(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		startTime := time.Now()

		handler(ctx)

		duration := time.Since(startTime)

		fmt.Println("Request URI: ", ctx.URI())
		fmt.Println("Request duration: ", duration)
	}
}

func authMiddleware(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if accessIsAllowed(ctx) {
			ctx.SetUserValue("access-token", "admin-token")

			handler(ctx)
		} else {
			ctx.SetStatusCode(fasthttp.StatusForbidden)
		}
	}
}

func accessIsAllowed(ctx *fasthttp.RequestCtx) bool {
	path := string(ctx.Path())
	needAuth := strings.HasPrefix(path, "/auth")

	if !needAuth {
		return true
	}

	return path == "/auth/admin"
}
