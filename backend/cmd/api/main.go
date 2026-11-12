// Command api serves the notebook. Without -listen it is the Lambda function behind the function URL.
// With -listen it is a local HTTP server. With -assets the local server also serves the built site.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/aws/aws-lambda-go/lambdaurl"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rustatian/redlamp.me/backend/internal/db"
	"github.com/rustatian/redlamp.me/backend/internal/notebook"
	"github.com/rustatian/redlamp.me/backend/internal/sets/content"
	"github.com/rustatian/redlamp.me/backend/internal/transcribe"
)

func main() {
	listen := flag.String("listen", "", "address of a local HTTP server, for example localhost:8787")
	url := flag.String("database-url", "", "PostgreSQL connection string, instead of the DSQL cluster of DSQL_ENDPOINT")
	assets := flag.String("assets", "", "folder of the built site for a local server, for example ../dist")
	flag.Parse()

	ctx := context.Background()
	var pool *pgxpool.Pool
	// A server with -database-url has no speech service, so its dictation is not available.
	var speech notebook.Transcriber
	var err error
	if *url != "" {
		pool, err = db.OpenURL(ctx, *url)
	} else {
		cfg, cfgErr := config.LoadDefaultConfig(ctx)
		if cfgErr != nil {
			log.Fatal(cfgErr)
		}
		pool, err = db.Open(ctx, db.Config{Endpoint: os.Getenv("DSQL_ENDPOINT"), Region: cfg.Region, User: "app", Credentials: cfg.Credentials})
		speech = transcribe.New(cfg)
	}
	if err != nil {
		log.Fatal(err)
	}
	var files http.Handler
	if *assets != "" {
		files = site(*assets)
	}
	srv, err := notebook.New(notebook.Config{DB: pool, Assets: files, PublicOrigin: os.Getenv("PUBLIC_ORIGIN"), Accounts: secret(ctx, "LEETCODE_ACCOUNTS"), LeetCode: notebook.GraphQL{URL: "https://leetcode.com/graphql"}, Sets: content.Sets, References: content.References, Generators: content.Generators, Anthropic: tutorClient(ctx), Transcriber: speech})
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

// secret returns the value of the variable name, or the value of the SecureString parameter of SSM Parameter Store that the variable name_PARAMETER names.
// The Lambda function has the parameter names. A local server has the values.
func secret(ctx context.Context, name string) string {
	value := os.Getenv(name)
	if parameter := os.Getenv(name + "_PARAMETER"); value == "" && parameter != "" {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err == nil {
			var found *ssm.GetParameterOutput
			if found, err = ssm.NewFromConfig(cfg).GetParameter(ctx, &ssm.GetParameterInput{Name: &parameter, WithDecryption: aws.Bool(true)}); err == nil {
				value = aws.ToString(found.Parameter.Value)
			}
		}
		if err != nil {
			log.Printf("the parameter of %s is not available: %v", name, err)
		}
	}
	return value
}

// tutorClient returns the client of the Anthropic API, or nil when the server has no API key.
func tutorClient(ctx context.Context) *anthropic.Client {
	key := secret(ctx, "ANTHROPIC_API_KEY")
	if key == "" {
		return nil
	}
	client := anthropic.NewClient(option.WithAPIKey(key))
	return &client
}

// site serves the built site for a local server. An address without a file suffix is the HTML file with that name,
// so the public addresses work as behind CloudFront, for example /examples and /mathtools/primefactorization.
func site(dir string) http.Handler {
	files := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/") && path.Ext(r.URL.Path) == "" {
			r = r.Clone(r.Context())
			r.URL.Path += ".html"
		}
		files.ServeHTTP(w, r)
	})
}
