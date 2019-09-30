package main

import (
	"fmt"
	"log"
	"os"

	"github.com/namsral/flag"
	"github.com/prologic/bitcask"
)

var (
	db  *bitcask.Bitcask
	cfg Config
)

func main() {
	var (
		version    bool
		config     string
		dbpath     string
		title      string
		fqdn       string
		bind       string
		url        string
		suggestURL string
	)

	flag.BoolVar(&version, "v", false, "display version information")

	flag.StringVar(&config, "config", "", "config file")
	flag.StringVar(&dbpath, "dbpath", "search.db", "database path")
	flag.StringVar(&title, "title", "Search", "OpenSearch title")
	flag.StringVar(&bind, "bind", "0.0.0.0:8000", "[int]:<port> to bind to")
	flag.StringVar(&fqdn, "fqdn", "localhost:8000", "FQDN for public access")
	flag.StringVar(&url, "url", DefaultURL, "default URL to redirect to")
	flag.StringVar(&suggestURL, "suggest", DefaultSuggestURL,
		"default URL to retrieve search suggestions from")

	flag.Parse()

	if version {
		fmt.Println(FullVersion())
		os.Exit(0)
	}

	cfg.Title = title
	cfg.FQDN = fqdn
	cfg.URL = url
	cfg.SuggestURL = suggestURL

	var err error
	db, err = bitcask.Open(dbpath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = EnsureDefaultBookmarks()
	if err != nil {
		log.Fatal(err)
	}

	NewServer(bind, cfg).ListenAndServe()
}
