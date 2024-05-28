package batch

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/dreamPathsProjekt/gosdelnet/pkg/client"
	"github.com/rs/zerolog/log"
)

func SearchByISBNFromCSV(file string, price bool, opts *client.Opts, verbose bool) ([]client.Book, error) {
	var (
		isbns []string
		books []client.Book
	)

	inFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer inFile.Close()

	reader := csv.NewReader(inFile)
	inputData, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Print the CSV data
	for i, row := range inputData {
		if i == 0 {
			continue
		}

		isbn, priceTracking := row[0], row[1]
		track, err := strconv.ParseBool(priceTracking)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to parse price tracking value as boolean: %s", priceTracking)
			continue
		}

		if price {
			if track {
				isbns = append(isbns, isbn)
			} else {
				continue
			}
		} else {
			isbns = append(isbns, isbn)
		}
	}

	if len(isbns) > 0 {
		for _, isbn13 := range isbns {
			log.Info().Msgf("Searching for ISBN-13: %s", isbn13)
			q := fmt.Sprintf("isbn13_search:(%s)", isbn13)
			r := int64(1)

			opts.Query = &q
			opts.Rows = &r

			client, err := client.New(*opts)
			if err != nil {
				return nil, err
			}

			result, err := client.Do(context.TODO(), verbose)
			if err != nil {
				return nil, err
			}

			books = append(books, result.Response.Docs...)
		}
	}

	return books, nil
}

func WriteResultsToCSV(file string, books []client.Book, price bool) error {
	outFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer outFile.Close()

	client.CSVHeader(outFile, price)

	for _, book := range books {
		book.CSVRow(outFile, price)
	}

	return nil
}
