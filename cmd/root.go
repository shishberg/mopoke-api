package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/juju/errors"
	"github.com/shishberg/mopoke-api/server"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	db, err := mongo.Connect(cmd.Context(), options.Client().ApplyURI(args.mongoURL))
	if err != nil {
		log.Fatal(errors.ErrorStack(err))
	}
	func() {
		pingCtx, cancel := context.WithTimeout(cmd.Context(), 10*time.Second)
		defer cancel()
		if err := db.Ping(pingCtx, nil); err != nil {
			log.Fatal(err)
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/ok", server.JSONHandler(
		func(w http.ResponseWriter, r *http.Request) (any, error) {
			return struct{}{}, nil
		},
	))

	handler := http.StripPrefix("/e", server.HandleEntity(db))
	v1mux := http.NewServeMux()
	v1mux.Handle("/e/", handler)
	v1mux.Handle("/e", handler)

	mux.Handle("/mopoke/v1/", http.StripPrefix("/mopoke/v1", v1mux))

	addr := fmt.Sprintf("0.0.0.0:%d", args.port)
	fmt.Println("Listening on", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
