package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

type Client struct {
	client *resty.Client
}

type Opts struct {
	URL              string
	MaxRetryCount    *int           // Max retry count for the request. Default is 3.
	RetryWaitTime    *time.Duration // Wait time between retry attempts. Default is 1 second.
	MaxRetryWaitTime *time.Duration // Maximum wait time between retry attempts. Default is 3 seconds.
	User             string         // Username for the authorization of the request.
	Password         string         // Password for the authorization of the request.
	Rows             *int64         // Number of rows to fetch. Default is 10.
	Query            *string        // Query to search for. Default is "*:*".
}

func New(opts Opts) (*Client, error) {
	base, err := url.Parse(opts.URL)
	if err != nil || base.Scheme == "" || base.Host == "" {
		return nil, err
	}

	var rows int64
	var query string = "*:*"

	client := resty.New()
	client.SetBaseURL(base.String())

	if opts.MaxRetryCount != nil {
		client.SetRetryCount(*opts.MaxRetryCount)
	} else {
		client.SetRetryCount(MaxRetryCount)
	}

	if opts.RetryWaitTime != nil {
		client.SetRetryWaitTime(*opts.RetryWaitTime)
	} else {
		client.SetRetryWaitTime(RetryWaitTime)
	}

	if opts.MaxRetryWaitTime != nil {
		client.SetRetryMaxWaitTime(*opts.MaxRetryWaitTime)
	} else {
		client.SetRetryMaxWaitTime(MaxRetryWaitTime)
	}

	if opts.Rows != nil {
		rows = *opts.Rows
	} else {
		rows = Rows
	}

	if opts.Query != nil {
		query = *opts.Query
	}

	client.SetQueryParams(map[string]string{
		"user": opts.User,
		"pass": opts.Password,
		"q":    query,
		"wt":   "json",
		"rows": strconv.FormatInt(rows, 10),
	})

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Do(ctx context.Context, verbose bool) (*Response, error) {
	var result Response

	request := c.client.R().SetContext(context.TODO())
	response, err := request.Get("")
	if err != nil {
		return nil, fmt.Errorf("error during request: %v", err)
	}

	// solr.osdelnet.gr doesn't support status code responses.
	if response.IsError() {
		err = response.Error().(*resty.ResponseError)
		return nil, fmt.Errorf("error in response: %v", err)
	}

	// Response body is text/plain, cannot accept application/json.
	err = json.Unmarshal([]byte(response.String()), &result)
	if err != nil {
		// solr.osdelnet.gr returns error messages with status 200, logged in Msg.
		log.Error().Err(err).Int("status", response.StatusCode()).Msg(response.String())
		return nil, fmt.Errorf("error casting to Response: %v", err)
	}

	if verbose {
		log.Info().Interface("response", result.Response).Int("status", response.StatusCode()).Msg("Request successful")
	}
	log.Info().
		Interface("responseHeader", result.ResponseHeader).
		Int("status", response.StatusCode()).
		Int64("found", result.Response.NumFound).
		Msg("Request successful")

	return &result, nil
}
