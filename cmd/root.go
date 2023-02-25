package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

var (
	args struct {
		mongoURL string
		port     int
	}

	rootCmd = &cobra.Command{
		Use: "mopoke-api",
		Run: run,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&args.mongoURL, "db", "mongodb://localhost:27017", "mongodb address")
	rootCmd.PersistentFlags().IntVar(&args.port, "port", 3090, "service port")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(errors.ErrorStack(err))
	}
}

func run(cmd *cobra.Command, argv []string) {
	addr := fmt.Sprintf("0.0.0.0:%d", args.port)
	fmt.Println("Listening on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
