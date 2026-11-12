// Command api serves the notebook. Without -listen it is the Lambda function behind the function URL.
// With -listen it is a local HTTP server.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambdaurl"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rustatian/redlamp.me/backend/internal/db"
	"github.com/rustatian/redlamp.me/backend/internal/notebook"
)

func main() {
	listen := flag.String("listen", "", "address of a local HTTP server, for example localhost:8787")
	url := flag.String("database-url", "", "PostgreSQL connection string, instead of the DSQL cluster of DSQL_ENDPOINT")
	flag.Parse()

	ctx := context.Background()
	var pool *pgxpool.Pool
	var err error
	if *url != "" {
		pool, err = db.OpenURL(ctx, *url)
	} else {
		cfg, cfgErr := config.LoadDefaultConfig(ctx)
		if cfgErr != nil {
			log.Fatal(cfgErr)
		}
		pool, err = db.Open(ctx, db.Config{Endpoint: os.Getenv("DSQL_ENDPOINT"), Region: cfg.Region, User: "app", Credentials: cfg.Credentials})
	}
	if err != nil {
		log.Fatal(err)
	}
	srv, err := notebook.New(notebook.Config{DB: pool, PublicOrigin: os.Getenv("PUBLIC_ORIGIN"), Accounts: os.Getenv("LEETCODE_ACCOUNTS"), LeetCode: notebook.GraphQL{URL: "https://leetcode.com/graphql"}})
	if err != nil {
		log.Fatalf("PUBLIC_ORIGIN: %v", err)
	}
	if *listen == "" {
		// The function URL must have the invoke mode RESPONSE_STREAM.
		// https://pkg.go.dev/github.com/aws/aws-lambda-go/lambdaurl#Start
		lambdaurl.Start(srv)
		return
	}
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("database: %v", err)
	}
	log.Printf("listening on %s", *listen)
	log.Fatal(http.ListenAndServe(*listen, srv))
}
