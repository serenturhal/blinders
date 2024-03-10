package main

import (
	"fmt"
	"log"
	"os"

	"blinders/cli/commands"

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
		Commands: []*cli.Command{&commands.AuthCommand},
		Before: func(ctx *cli.Context) error {
			fmt.Println("running on", ctx.String("env"))

			envFile := fmt.Sprintf(".env.%s", ctx.String("env"))
			if godotenv.Load(envFile) != nil {
				log.Fatal("Error loading .env file ", envFile)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
