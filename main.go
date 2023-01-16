package main

import (
	"context"
	"gitlab.shoplazza.site/shoplaza/cobra/cmd"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	filename "github.com/keepeye/logrus-filename"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/migrations"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fnhook := filename.NewHook()
	logrus.AddHook(fnhook)
	ctx, cancel := context.WithCancel(context.Background())

	app := cli.NewApp()
	app.Name = "cobra"
	app.Usage = "all app interface"

	app.Commands = []cli.Command{
		{
			Name:  "web",
			Usage: "start web server",
			Action: func(c *cli.Context) error {
				cmd.StartMetricServer()
				return cmd.StartServerAndBlock(ctx)
			},
		},
		{
			Name:  "migrate",
			Usage: "db migration",
			Action: func(c *cli.Context) error {
				return migrations.Migrate(config.DB)
			},
		},
	}

	sigcomplete := make(chan struct{})
	go func() {
		defer close(sigcomplete)
		err := app.Run(os.Args)
		if err != nil {
			log.WithError(err).Fatal("app run failed")
		}
	}()

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sigterm:
		log.Info("receive stop signal")
	case <-sigcomplete:
	}

	cancel()
	config.DB.Close()
	if config.Producer != nil {
		config.Producer.Close()
	}
	//wait for other goroutines to quit
	time.Sleep(3 * time.Second)

}
