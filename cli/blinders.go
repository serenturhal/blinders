package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Blinders",
		Usage: "CLI tools for backend development",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "development",
				Usage: "Define environment for the CLI",
			},
		},
		Commands: []*cli.Command{},
		Before: func(ctx *cli.Context) error {
			envFile := fmt.Sprintf(".env.%s", ctx.String("env"))
			if godotenv.Load(envFile) != nil {
				log.Fatal("Error loading .env file", envFile)
			}

			fmt.Println("[cli] running at", os.Getenv("ENVIRONMENT"))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
