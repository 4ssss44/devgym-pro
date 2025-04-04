package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/devgymbr/kong"
	"github.com/devgymbr/kong/config"
	internalhttp "github.com/devgymbr/kong/http"
	"github.com/devgymbr/kong/plugin"
)

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(l)

	plugin.Register("http_log", plugin.Log)
	plugin.Register("add_header", plugin.AddHeader)
	plugin.Register("jwt_auth", plugin.JWTAuth)
	plugin.Register("request_size_limiting", plugin.RequestSizeLimit)

	c := &kong.Config{}
	if err := config.Loader(c, "config.yaml", 5*time.Second); err != nil {
		panic(err)
	}

	server := internalhttp.NewServer(c)

	fmt.Println("Server listening on port 8080...")

	http.Handle("/", server)
	http.ListenAndServe(":8080", nil)
}
