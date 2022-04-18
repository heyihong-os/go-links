package main

import (
	"context"
	"log"
	"net/http"
	"regexp"

	"github.com/heyihong-os/go-links/storage"
	"github.com/jessevdk/go-flags"
	"github.com/julienschmidt/httprouter"

	"cloud.google.com/go/bigquery"
)

var opts Options

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		return
	}

	allowedEmailRegex, err := regexp.Compile(opts.AllowedEmailRegex)
	if err != nil {
		log.Fatalf("Failed to compile regex %s", opts.AllowedEmailRegex)
	}

	ctx := context.Background()
	bqClient, err := bigquery.NewClient(ctx, opts.BigQuery.ProjectID)
	if err != nil {
		log.Fatalf("Failed to initialize big query client for project %s", opts.BigQuery.ProjectID)
	}
	defer bqClient.Close()

	ls, err := NewLinkStore(ctx, storage.NewBigQueryStorage(bqClient, opts.BigQuery.DatasetName, opts.BigQuery.TableName))
	if err != nil {
		log.Fatalf("Failed to initialize link store: %v", err)
	}

	rh := NewRequestHandler(ls)

	log.Print("Starting server...")
	router := httprouter.New()
	router.GET("/", rh.handleHome)
	router.POST("/_api/update", rh.handleApiUpdate)

	router.HandleMethodNotAllowed = false
	router.NotFound = http.HandlerFunc(rh.handleNotFound)

	// Start HTTP server.
	log.Printf("Listening on port %s", opts.Port)
	if err := http.ListenAndServe(":"+opts.Port, AuthzHandler(allowedEmailRegex, router)); err != nil {
		log.Fatal(err)
	}
}
