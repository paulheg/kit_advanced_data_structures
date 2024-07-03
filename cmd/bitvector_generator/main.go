package main

import (
	"log"
	"os"
	"path"

	"github.com/paulheg/kit_advanced_data_structures/internal/bitvector"
	"github.com/urfave/cli/v2"
)

func main() {

	app := cli.App{
		Name: "generator",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name: "output-dir",
				Aliases: []string{
					"o", "output",
				},
				Value: ".",
			},
			&cli.Uint64Flag{
				Name: "vector-length",
				Aliases: []string{
					"l", "length",
				},
				Value: 10,
			},
			&cli.Uint64Flag{
				Name: "commands",
				Aliases: []string{
					"n", "c",
				},
				Value: 10,
			},
		},
		Action: func(ctx *cli.Context) error {

			outputPath := ctx.Path("output-dir")

			fExpected, err := os.Create(path.Join(outputPath, "expected.txt"))
			if err != nil {
				return err
			}
			defer fExpected.Close()

			fCommands, err := os.Create(path.Join(outputPath, "commands.txt"))
			if err != nil {
				return err
			}
			defer fCommands.Close()

			err = bitvector.GenerateRandomTestCase(ctx.Uint64("vector-length"), ctx.Uint64("commands"), fCommands, fExpected)
			if err != nil {
				return err
			}

			log.Println("finished...")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
