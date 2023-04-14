//go:build ignore

package main

import (
	"context"
	"log"
	"os"
	"path"

	"github.com/katallaxie/autobot/internal/training"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type flags struct {
	Locales []string
}

var f = &flags{}

func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCommand.Flags().StringSliceVar(&f.Locales, "locales", f.Locales, "locales")
	RootCommand.SilenceUsage = true
}

var RootCommand = &cobra.Command{
	Use: "train",
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCommand(cmd.Context())
	},
}

func rootCommand(_ context.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, locale := range f.Locales {
		p := path.Join(cwd, "locales", locale, "training.json")

		err := os.Remove(path.Clean(p))
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		training.CreateNeuralNetwork(locale, true)
	}

	return nil
}

func main() {
	if err := RootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
