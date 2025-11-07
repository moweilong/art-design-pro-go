// nolint: err113
package options

import (
	genericoptions "github.com/moweilong/milady/pkg/options"
	"github.com/spf13/pflag"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/moweilong/art-design-pro-go/internal/apiserver"
)

// ServerOptions contains the configuration options for the server.
type ServerOptions struct {
	// TLSOptions contains the TLS configuration options.
	TLSOptions *genericoptions.TLSOptions `json:"tls" mapstructure:"tls"`
	// HTTPOptions contains the HTTP configuration options.
	HTTPOptions *genericoptions.HTTPOptions `json:"http" mapstructure:"http"`
	// MySQLOptions contains the MySQL configuration options.
	MySQLOptions *genericoptions.MySQLOptions `json:"mysql" mapstructure:"mysql"`
	// OTelOptions used to specify the otel options.
	OTelOptions *genericoptions.OTelOptions `json:"otel" mapstructure:"otel"`
	// JWT options for configuring JWT related options.
	JWTOptions *genericoptions.JWTOptions `json:"jwt" mapstructure:"jwt"`
	// Redis options for configuring Redis related options.
	RedisOptions *genericoptions.RedisOptions `json:"redis" mapstructure:"redis"`
}

// NewServerOptions creates a ServerOptions instance with default values.
func NewServerOptions() *ServerOptions {
	opts := &ServerOptions{
		TLSOptions:   genericoptions.NewTLSOptions(),
		HTTPOptions:  genericoptions.NewHTTPOptions(),
		JWTOptions:   genericoptions.NewJWTOptions(),
		RedisOptions: genericoptions.NewRedisOptions(),
		MySQLOptions: genericoptions.NewMySQLOptions(),
		OTelOptions:  genericoptions.NewOTelOptions(),
	}
	opts.HTTPOptions.Addr = ":5555"

	return opts
}

// AddFlags binds the options in ServerOptions to command-line flags.
func (o *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	// Add command-line flags for sub-options.
	o.TLSOptions.AddFlags(fs)
	o.HTTPOptions.AddFlags(fs)
	o.MySQLOptions.AddFlags(fs)
	o.OTelOptions.AddFlags(fs)
	o.JWTOptions.AddFlags(fs)
	o.RedisOptions.AddFlags(fs)
}

// Complete completes all the required options.
func (o *ServerOptions) Complete() error {
	// TODO: Add the completion logic if needed.
	return nil
}

// Validate checks whether the options in ServerOptions are valid.
func (o *ServerOptions) Validate() error {
	errs := []error{}

	// Validate sub-options.
	errs = append(errs, o.TLSOptions.Validate()...)
	errs = append(errs, o.HTTPOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.OTelOptions.Validate()...)
	errs = append(errs, o.JWTOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)

	// Aggregate all errors and return them.
	return utilerrors.NewAggregate(errs)
}

// Config builds an apiserver.Config based on ServerOptions.
func (o *ServerOptions) Config() (*apiserver.Config, error) {
	return &apiserver.Config{
		TLSOptions:   o.TLSOptions,
		HTTPOptions:  o.HTTPOptions,
		MySQLOptions: o.MySQLOptions,
		JWTOptions:   o.JWTOptions,
		RedisOptions: o.RedisOptions,
	}, nil
}
