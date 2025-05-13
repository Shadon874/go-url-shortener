package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"url-shortener/internal/api"
	"url-shortener/internal/db"
	"url-shortener/internal/logic"

	"github.com/urfave/cli/v3"
)

func main() {
	database := db.InitDB()
	defer database.Close()

	cmd := &cli.Command{
		Name:  "cuturl",
		Usage: "URL shortener CLI tool and server",
		Commands: []*cli.Command{
			{
				Name:  "cut-url",
				Usage: "Shorten a URL",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "no-url",
						Usage: "Do not add a URL to shorten",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Bool("no-url") {
						return nil
					}
					url := cmd.Args().Get(0)
					if url == "" {
						fmt.Println("Please provide a URL.")
						return fmt.Errorf("URL is required")
					} else {
						logic.AddUrlToDb(url, database)
					}
					return nil
				},
			},
			{
				Name:  "serve",
				Usage: "Start the server",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					println("Server on :8080...")
					http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
					http.HandleFunc("/", api.HandleReq(database))
					http.HandleFunc("/api/shorten", api.HandleApi(database))
					http.ListenAndServe(":8080", nil)
					return nil
				},
			},
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
