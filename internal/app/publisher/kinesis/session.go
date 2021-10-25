package kinesis

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
)

type SessionConfig struct {
	Endpoint   string
	Region     string
	DisableSSL bool
}

// NewSession creates a new AWS session.
func NewSession(cfg SessionConfig) (client.ConfigProvider, error) {
	awscfg := aws.NewConfig().
		WithCredentialsChainVerboseErrors(true).
		WithRegion(cfg.Region).
		WithDisableSSL(cfg.DisableSSL)

	if cfg.Endpoint != "" {
		awscfg = awscfg.WithEndpoint(cfg.Endpoint)
	}

	return session.NewSession(awscfg)
}
