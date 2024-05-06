package discord

import (
	"fmt"
	"log/slog"

	"github.com/aqyuki/discord-bot-template/pkg/logging"
	"github.com/bwmarrin/discordgo"
)

type ClientError struct {
	raw error
}

func NewClientError(err error) error {
	return &ClientError{raw: err}
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("discord client error: %v", e.raw)
}

func (e *ClientError) Unwrap() error {
	return e.raw
}

// Client is a struct to provide features to interact with Discord API
type Client struct {
	// logger is a logger to use logging
	logger *slog.Logger

	// config is the configuration for the Client
	config Config

	// session is the Discord Session
	session *discordgo.Session
}

// Config holds the configuration for the Client
type Config struct {
	// Token is the Discord Bot Token
	Token string
}

// DiscordConfigProvider is an interface to get the configuration for the Client
type DiscordConfigProvider interface {
	Config() Config
}

// Option ...
type Option func(*Client)

// WithLogger creates an Option to set the logger of the Client
func WithLogger(l *slog.Logger) Option {
	return func(c *Client) {
		c.logger = l
	}
}

// NewClient creates a new Client with the given options
func NewClient(provider DiscordConfigProvider, options ...Option) (*Client, error) {
	c := defaultClient()
	c.config = provider.Config()

	// create session
	session, err := discordgo.New("Bot " + c.config.Token)
	if err != nil {
		return nil, NewClientError(err)
	}
	c.session = session

	for _, f := range options {
		f(c)
	}
	return c, nil
}

// defaultClient is a helper function to create a new Client with default options
func defaultClient() *Client {
	return &Client{
		logger: logging.DefaultLogger(),
	}
}

func (c *Client) Open() error {
	if c.session == nil {
		panic("discord/session: session is nil. Did you created other way than NewClient?")
	}

	if err := c.session.Open(); err != nil {
		c.purge()
		return NewClientError(err)
	}
	return nil
}

func (c *Client) purge() {
	if c.session != nil {
		c.session.Close()
		c.session = nil
	}
}

func (c *Client) Session() *discordgo.Session {
	return c.session
}
