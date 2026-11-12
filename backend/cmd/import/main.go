// Command import moves the notebook data from files of the D1 database into Aurora DSQL or PostgreSQL. It runs one time, at the cutover.
//
// Make one file for each table of importer.Tables in one folder, for example problems.json:
//
//	wrangler d1 execute redlamp-leetcode --remote --json --command "SELECT * FROM problems" > problems.json
//
// The database must have the schema of cmd/migrate and no rows. After the import the command compares each table with its file.
package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rustatian/redlamp.me/backend/internal/db"
	"github.com/rustatian/redlamp.me/backend/internal/importer"
)

func main() {
	dir := flag.String("dir", "", "folder with one JSON file for each table")
	endpoint := flag.String("endpoint", "", "DSQL cluster endpoint")
	region := flag.String("region", "eu-central-1", "DSQL cluster region")
	url := flag.String("database-url", "", "PostgreSQL connection string, instead of -endpoint")
	check := flag.Bool("check", false, "only compare the database with the files")
	flag.Parse()
	if *dir == "" || (*endpoint == "") == (*url == "") {
		log.Fatal("give -dir and exactly one of -endpoint and -database-url")
	}

	ctx := context.Background()
	var pool *pgxpool.Pool
	var err error
	if *url != "" {
		pool, err = db.OpenURL(ctx, *url)
	} else {
		cfg, cfgErr := config.LoadDefaultConfig(ctx, config.WithRegion(*region))
		if cfgErr != nil {
			log.Fatal(cfgErr)
		}
		pool, err = db.Open(ctx, db.Config{Endpoint: *endpoint, Region: *region, User: "admin", Credentials: cfg.Credentials})
	}
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	run := importer.Run
	if *check {
		run = importer.Verify
	}
	if err := run(ctx, pool, *dir, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
