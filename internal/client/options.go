package client

type Option func(*Client) error

func (c *Client) Configure(opts ...Option) error {
	for _, op := range opts {
		if err := op(c); err != nil {
			return err
		}
	}
	return nil
}
