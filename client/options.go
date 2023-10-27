package client

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/frankli0324/go-requests/internal/client"
	"github.com/frankli0324/go-requests/utils"
)

func WithSession(options *cookiejar.Options) client.Option {
	return func(c *client.Client) error {
		jar, err := cookiejar.New(options)
		c.Client.Jar = jar
		return err
	}
}

// WithDisableH2 disables HTTP2 on the transport, only works
// with [http.Transport]
func WithDisableH2() client.Option {
	return func(c *client.Client) error {
		htr, ok := utils.GetHttpTransport(c.Client.Transport)
		if !ok {
			return errors.New("unsupport roundtripper")
		}
		if htr.TLSNextProto == nil {
			htr.TLSNextProto = make(map[string]func(string, *tls.Conn) http.RoundTripper)
		} else {
			delete(htr.TLSNextProto, "h2")
		}
		return nil
	}
}

// WithProxy sets the proxy for a client
func WithProxy(proxy string) client.Option {
	return func(c *client.Client) error {
		p, err := url.Parse(proxy)
		if err != nil {
			return err
		}
		htr, ok := utils.GetHttpTransport(c.Client.Transport)
		if !ok {
			return errors.New("unsupport roundtripper")
		}
		htr.Proxy = http.ProxyURL(p)
		return nil
	}
}

// WithProxyFunc sets the proxy getter for a client
func WithProxyFunc(proxyFunc func(*http.Request) (*url.URL, error)) client.Option {
	return func(c *client.Client) error {
		htr, ok := utils.GetHttpTransport(c.Client.Transport)
		if !ok {
			return errors.New("unsupport roundtripper")
		}
		htr.Proxy = proxyFunc
		return nil
	}
}
