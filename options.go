package pocketbook_cloud_client

type Option func(*Client)

func WithHTTPClient(client doer) Option {
	return func(c *Client) {
		c.http = client
	}
}

func WithClientID(id string) Option {
	return func(c *Client) {
		c.clientID = id
	}
}

func WithClientSecret(sec string) Option {
	return func(c *Client) {
		c.clientSecret = sec
	}
}
