package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type clientOptions struct {
	apiUrl        string
	httpClient    *http.Client
	timeFormat    string
	clientContext context.Context
	zipCodeRegex  *regexp.Regexp
	userAgent     string
}

type Option func(*clientOptions) error

// WithApiUrl changes the base url of the api
func WithApiUrl(apiUrl string) Option {
	return func(co *clientOptions) error {
		if !strings.HasSuffix(apiUrl, "/") {
			apiUrl += "/"
		}
		_, err := url.ParseRequestURI(apiUrl)
		if err != nil {
			return fmt.Errorf("invalid url: %s: %w", apiUrl, err)
		}
		co.apiUrl = apiUrl
		return nil
	}
}

// WithHTTPCLient changes the used http client
func WithHTTPClient(client *http.Client) Option {
	return func(co *clientOptions) error {
		if client == nil {
			return errors.New("http client is nil")
		}
		co.httpClient = client
		return nil
	}
}

// WithTimeFormat changes the format of the used time string
func WithTimeFormat(format string) Option {
	return func(co *clientOptions) error {
		co.timeFormat = format
		return nil
	}
}

// WithZipCodeRegex changes the format of the used time string
func WithZipCodeRegex(regex string) Option {
	return func(co *clientOptions) error {
		r, err := regexp.Compile(regex)
		if err != nil {
			return fmt.Errorf("invalid zip code regex: %w", err)
		}
		co.zipCodeRegex = r
		return nil
	}
}

// WithContext change sthe default context that is used for
// every method that does NOT expect a context
func WithContext(ctx context.Context) Option {
	return func(co *clientOptions) error {
		co.clientContext = ctx
		return nil
	}
}

// WithUserAgent sets the User-Agent for every client request
func WithUserAgent(userAgent string) Option {
	return func(co *clientOptions) error {
		co.userAgent = userAgent
		return nil
	}
}
