package cmd

import (
	"context"
	"log"
	"os"

	"github.com/katallaxie/autobot/internal/adapters/cmd"
	"github.com/katallaxie/autobot/internal/analyzer"
	"github.com/katallaxie/autobot/internal/controllers"
	"github.com/katallaxie/autobot/internal/locales"
	"github.com/katallaxie/autobot/internal/network"
	"github.com/katallaxie/autobot/internal/services/autobot"
	"github.com/katallaxie/autobot/internal/services/status"
	"github.com/katallaxie/autobot/internal/training"
	"github.com/katallaxie/autobot/pkg/adapter"
	"github.com/katallaxie/autobot/util"

	"github.com/katallaxie/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type flags struct {
	Addr       string
	Host       string
	Token      string
	StatusAddr string
	Adapter    string
	Locales    []string
	Verbose    bool
}

var f = &flags{}

var RootCommand = &cobra.Command{
	Use: "autobot",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(cmd.Context())
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCommand.PersistentFlags().BoolVarP(&f.Verbose, "verbose", "v", f.Verbose, "verbose output")
	RootCommand.PersistentFlags().StringSliceVar(&f.Locales, "locale", f.Locales, "locales")
	RootCommand.Flags().StringVar(&f.Adapter, "adapter", f.Adapter, "adapter")
	RootCommand.Flags().StringVar(&f.Addr, "addr", "0.0.0.0:8080", "address")
	RootCommand.Flags().StringVar(&f.StatusAddr, "status-addr", "0.0.0.0:8081", "status address")

	RootCommand.SilenceUsage = true
}

func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}

func run(ctx context.Context) error {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	nn := make(map[string]network.Network)

	for _, locale := range locales.Locales {
		util.SerializeMessages(locale.Tag)

		nn[locale.Tag] = training.CreateNeuralNetwork(locale.Tag, false)
	}

	meta := &adapter.Meta{
		Path: f.Adapter,
	}
	a := meta.Factory(ctx)

	adapter, err := a()
	if err != nil {
		return err
	}
	defer func() { _ = adapter.Close() }()

	analyzer := analyzer.NewAnalyzer()
	cmds := cmd.New()
	cmd.Register(2, cmd.NewDebug())

	m := controllers.NewMessageController(adapter, cmds, analyzer, nn)

	srv, _ := server.WithContext(ctx)

	s := status.New(status.WithAddr(f.StatusAddr))
	srv.Listen(s, false)

	r := autobot.New(m)
	srv.Listen(r, false)

	if err := srv.Wait(); err != nil {
		return err
	}

	return nil
}
