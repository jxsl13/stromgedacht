package client

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jxsl13/stromgedacht/api"
)

var (
	// DefaultDateTimeFormat may be changed for the whole package to use a different format
	// look in the time package for further information
	DefaultDateTimeFormat = "2006-01-02T15:04:05-07:00"

	// DefaultZipCodeRegex is used to validate German zip codes
	DefaultZipCodeRegex = regexp.MustCompile(`^((?:0[1-46-9]\d{3})|(?:[1-357-9]\d{4})|(?:[4][0-24-9]\d{3})|(?:[6][013-9]\d{3}))$`)

	// DefaultUserAgent is the default user agent of this library
	DefaultUserAgent = "github.com/jxsl13/stromgedacht"
)

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	client       *resty.Client
	ctx          context.Context
	timeFormat   string
	zipCodeRegex *regexp.Regexp
}

// Creates a new Client, with reasonable defaults
func New(opts ...Option) (*Client, error) {
	options := clientOptions{
		apiUrl:        "https://api.stromgedacht.de/v1/",
		httpClient:    &http.Client{},
		timeFormat:    DefaultDateTimeFormat,
		clientContext: context.Background(),
		zipCodeRegex:  DefaultZipCodeRegex,
		userAgent:     DefaultUserAgent,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&options); err != nil {
			return nil, err
		}
	}

	client := &Client{
		client: resty.NewWithClient(options.httpClient).
			SetBaseURL(options.apiUrl).
			SetHeader("User-Agent", options.userAgent),
		timeFormat:   options.timeFormat,
		zipCodeRegex: options.zipCodeRegex,
	}

	return client, nil
}

// GetNow Retrieval of the StromGedacht status in the TransnetBW control area
// for the current time at a specified location (postal code).
func (c *Client) GetNow(zip string) (*api.RegionStateNowViewModel, error) {
	return c.GetNowContext(c.ctx, zip)
}

// GetNowContext Retrieval of the StromGedacht status in the TransnetBW control area
// for the current time at a specified location (postal code).
func (c *Client) GetNowContext(ctx context.Context, zip string) (*api.RegionStateNowViewModel, error) {
	var result api.RegionStateNowViewModel
	if !c.zipCodeRegex.MatchString(zip) {
		return nil, fmt.Errorf("invalid zip code: %s", zip)
	}

	resp, err := c.client.R().
		SetHeader("Accept", "application/json").
		SetQueryParam("zip", zip).
		SetContext(ctx).
		SetResult(&result).
		Get("/now")
	if err != nil {
		return nil, err
	}

	if resp.IsSuccess() {
		return &result, nil
	}
	return nil, fmt.Errorf("failed to get /now state: %s", string(resp.Body()))
}

// GetStates Retrieval of the StromGedacht statuses in the TransnetBW control area for a specified
// time period at a specified location (zip code), whereby the status can be retrieved
// a maximum of 4 days into the past and a maximum of 2 days into the future.
func (c *Client) GetStates(zip string, from, to time.Time) (*api.RegionStateRangeViewModel, error) {
	return c.GetStatesContext(c.ctx, zip, from, to)
}

// GetStatesContext Retrieval of the StromGedacht statuses in the TransnetBW control area for a specified
// time period at a specified location (zip code), whereby the status can be retrieved
// a maximum of 4 days into the past and a maximum of 2 days into the future.
func (c *Client) GetStatesContext(ctx context.Context, zip string, from, to time.Time) (*api.RegionStateRangeViewModel, error) {
	if !c.zipCodeRegex.MatchString(zip) {
		return nil, fmt.Errorf("invalid zip code: %s", zip)
	}

	if to.Before(from) {
		return nil, fmt.Errorf("from date time %s is before to date time %s", from, to)
	}

	var (
		now            = time.Now().Truncate(time.Second)
		hour, min, sec = now.Clock()
		offset         = time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute + time.Duration(sec)*time.Second
		lb             = now.AddDate(0, 0, -4).Add(-offset)
		ub             = now.AddDate(0, 0, 3).Add(-offset - time.Second)
	)

	if from.Before(lb) {
		return nil, fmt.Errorf("minimum allowed from date time is: %s", lb)
	}

	if to.After(ub) {
		return nil, fmt.Errorf("maximum allowed to date time is: %s", ub)
	}

	var result api.RegionStateRangeViewModel
	resp, err := c.client.R().
		SetHeader("Accept", "application/json").
		SetQueryParams(map[string]string{
			"zip":  zip,
			"from": from.Format(c.timeFormat),
			"to":   to.Format(c.timeFormat),
		}).
		SetContext(ctx).
		SetResult(&result).
		Get("/states")
	if err != nil {
		return nil, err
	}

	if resp.IsSuccess() {
		return &result, nil
	}
	return nil, fmt.Errorf("failed to get /states state: %s", string(resp.Body()))
}
