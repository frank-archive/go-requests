package client

import (
	"net/http/cookiejar"

	"github.com/frankli0324/go-requests/internal/client"
)

func WithSession(options *cookiejar.Options) client.Option {
	return func(c *client.Client) error {
		jar, err := cookiejar.New(options)
		c.Client.Jar = jar
		return err
	}
}
