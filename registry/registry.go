package registry

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Registry struct {
	URL    string
	Client *http.Client
	Quiet  bool
}

func New(registryUrl, username, password string) (*Registry, error) {
	transport := http.DefaultTransport

	return newFromTransport(registryUrl, username, password, transport, false)
}

func NewQuiet(registryUrl, username, password string) (*Registry, error) {
	transport := http.DefaultTransport

	return newFromTransport(registryUrl, username, password, transport, true)
}

func NewInsecure(registryUrl, username, password string) (*Registry, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return newFromTransport(registryUrl, username, password, transport, false)
}

func newFromTransport(registryUrl, username, password string, transport http.RoundTripper, quiet bool) (*Registry, error) {
	url := strings.TrimSuffix(registryUrl, "/")
	transport = &TokenTransport{
		Transport: transport,
		Username:  username,
		Password:  password,
	}
	transport = &BasicTransport{
		Transport: transport,
		URL:       url,
		Username:  username,
		Password:  password,
	}
	transport = &ErrorTransport{
		Transport: transport,
	}

	registry := &Registry{
		URL: url,
		Client: &http.Client{
			Transport: transport,
		},
		Quiet: quiet,
	}

	if err := registry.Ping(); err != nil {
		return nil, err
	}

	return registry, nil
}

func (r *Registry) url(pathTemplate string, args ...interface{}) string {
	pathSuffix := fmt.Sprintf(pathTemplate, args...)
	url := fmt.Sprintf("%s%s", r.URL, pathSuffix)
	return url
}

func (r *Registry) Ping() error {
	url := r.url("/v2/")
	if !r.Quiet {
		log.Printf("registry.ping url=%s", url)
	}
	resp, err := r.Client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {

	}
	return err
}
