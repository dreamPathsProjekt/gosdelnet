package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	gosdel "github.com/dreamPathsProjekt/gosdelnet/pkg/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	timeFmt = time.RFC3339
	baseUrl = "http://solr.osdelnet.gr/solr/index.php"
)

func main() {
	var (
		query string = "*:*"
	)

	time.Local = time.UTC
	zerolog.TimeFieldFormat = timeFmt

	user := os.Getenv("OSDELNET_USER")
	password := os.Getenv("OSDELNET_PASSWORD")

	file := flag.String("file", "", "specify an input csv file to read from")
	rows := flag.Int64("rows", 10, "number of rows to fetch")
	isbn := flag.String("isbn", "", "specify an ISBN-13 to search for, can use prefix too (e.g. 978-960-451-482 instead of full ISBN-13 978-960-451-482-3)")
	publisher := flag.String("publisher", "", "specify a publisher to search for")
	csvFile := flag.String("csv", "", "specify an output csv file to write to")
	verbose := flag.Bool("verbose", false, "enable verbose logging")

	flag.Parse()

	if *file != "" {
		log.Info().Str("file", *file).Msgf("Reading from CSV file: %s", *file)
	} else {
		log.Info().Msg("No file specified, parsing all book results")
	}

	if *isbn != "" {
		log.Info().Str("isbn", *isbn).Msgf("Searching for ISBN-13: %s", *isbn)
		query = fmt.Sprintf("isbn13_search:(%s)", *isbn)
	} else {
		log.Info().Msg("No ISBN-13 specified, searching all books")
	}

	if *publisher != "" {
		log.Info().Str("publisher", *publisher).Msgf("Searching for publisher: %s", *publisher)
		query = fmt.Sprintf("imprint_search:(%s)", *publisher)
	} else {
		log.Info().Msg("No publisher specified, searching all publishers")
	}

	client, err := gosdel.New(gosdel.Opts{
		URL:      baseUrl,
		User:     user,
		Password: password,
		Rows:     rows,
		Query:    &query,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create client")
	}

	result, err := client.Do(context.TODO(), *verbose)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if *csvFile != "" {
		log.Info().Str("csv", *csvFile).Msgf("Writing to CSV file: %s", *csvFile)
		file, err := os.Create(*csvFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		books := result.Response.Docs
		gosdel.CSVHeader(file)

		for _, book := range books {
			book.CSVRow(file)
		}
	}
}
