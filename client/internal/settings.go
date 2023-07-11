package internal

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/aidtechnology/did-method/info"
	"github.com/aidtechnology/did-method/resolver"
	"github.com/spf13/viper"
	"go.bryk.io/pkg/cli"
	xlog "go.bryk.io/pkg/log"
	mwHeaders "go.bryk.io/pkg/net/middleware/headers"
	"go.bryk.io/pkg/net/rpc"
	otelSdk "go.bryk.io/pkg/otel/sdk"
	"go.bryk.io/pkg/otel/sentry"
)

// Default service name; used on telemetry and reporting data.
const serviceName string = "did-bryk"

// Settings provide utilities to manage configuration options available
// when utilizing the different components available through the CLI `didctl`
// application.
type Settings struct {
	// Agent holds the configuration required to operate a "service provider" RPC node.
	Agent *agent `json:"agent" yaml:"agent" mapstructure:"agent"`

	// Client holds the configuration required to connect to an agent node.
	Client *client `json:"client" yaml:"client" mapstructure:"client"`

	// Resolver determines the registered DID providers to be used on a "resolve"
	// operations.
	Resolver []*resolver.Provider `json:"resolver" yaml:"resolver" mapstructure:"resolver"`
}

// Overrides return the available flag overrides for the command specified.
// Specific settings can be provided via: configuration file, ENV variable
// and command flags.
func (s *Settings) Overrides(cmd string) []cli.Param {
	switch cmd {
	case "agent":
		return []cli.Param{
			{
				Name:      "pow",
				Usage:     "set the required request ticket difficulty level",
				FlagKey:   "agent.pow",
				ByDefault: 16,
			},
			{
				Name:      "storage",
				Usage:     "specify storage mechanism to use",
				FlagKey:   "agent.storage",
				ByDefault: "ephemeral",
				Short:     "s",
			},
			{
				Name:      "port",
				Usage:     "TCP port to use for the server",
				FlagKey:   "agent.rpc.port",
				ByDefault: 9090,
				Short:     "p",
			},
			{
				Name:      "http",
				Usage:     "enable the HTTP interface",
				FlagKey:   "agent.rpc.http",
				ByDefault: false,
			},
			{
				Name:      "tls",
				Usage:     "enable secure communications using TLS with provided credentials",
				FlagKey:   "agent.rpc.tls.enabled",
				ByDefault: false,
			},
			{
				Name:      "tls-ca",
				Usage:     "TLS custom certificate authority (path to PEM file)",
				FlagKey:   "agent.rpc.tls.ca",
				ByDefault: "",
			},
			{
				Name:      "tls-cert",
				Usage:     "TLS certificate (path to PEM file)",
				FlagKey:   "agent.rpc.tls.cert",
				ByDefault: "/etc/didctl/tls/tls.crt",
			},
			{
				Name:      "tls-key",
				Usage:     "TLS private key (path to PEM file)",
				FlagKey:   "agent.rpc.tls.key",
				ByDefault: "/etc/didctl/tls/tls.key",
			},
			{
				Name:      "method",
				Usage:     "specify a supported DID method (can be provided multiple times)",
				FlagKey:   "agent.methods",
				ByDefault: []string{"bryk"},
				Short:     "m",
			},
		}
	case "client":
		return []cli.Param{
			{
				Name:      "key",
				Usage:     "cryptographic key to use for the sync operation",
				FlagKey:   "sync.key",
				ByDefault: "master",
				Short:     "k",
			},
			{
				Name:      "deactivate",
				Usage:     "instruct the network agent to deactivate the identifier",
				FlagKey:   "sync.deactivate",
				ByDefault: false,
				Short:     "d",
			},
			{
				Name:      "pow",
				Usage:     "set the required request ticket difficulty level",
				FlagKey:   "client.pow",
				ByDefault: 16,
				Short:     "p",
			},
			{
				Name:      "timeout",
				Usage:     "max time (in seconds) to wait for the agent to respond",
				FlagKey:   "client.timeout",
				ByDefault: 5,
				Short:     "t",
			},
		}
	default:
		return []cli.Param{}
	}
}

// Load available config values from Viper into the settings instance.
func (s *Settings) Load(v *viper.Viper) error {
	return v.Unmarshal(s)
}

// OTEL returns the configuration options available to set up an OTEL operator.
func (s *Settings) OTEL(log xlog.Logger) []otelSdk.Option {
	opts := []otelSdk.Option{
		otelSdk.WithBaseLogger(log),
		otelSdk.WithServiceName(serviceName),
		otelSdk.WithServiceVersion(info.CoreVersion),
		otelSdk.WithResourceAttributes(s.Agent.OTEL.Attributes),
	}
	collector := s.Agent.OTEL.Collector
	if collector != "" {
		opts = append(opts, otelSdk.WithExporterOTLP(collector, true, nil)...)
	}

	// Error reporter
	if sentryOpts := s.Agent.OTEL.Sentry; sentryOpts.DSN != "" {
		if sentryOpts.Release == "" {
			sentryOpts.Release = s.ReleaseCode()
		}
		rep, err := sentry.NewReporter(sentryOpts)
		if err == nil {
			opts = append(opts,
				otelSdk.WithPropagator(rep.Propagator()),
				otelSdk.WithSpanProcessor(rep.SpanProcessor()),
			)
		}
	}
	return opts
}

// Server returns the configuration options available to set up an RPC server.
func (s *Settings) Server() ([]rpc.ServerOption, error) {
	opts := []rpc.ServerOption{
		rpc.WithPanicRecovery(),
		rpc.WithInputValidation(),
		rpc.WithReflection(),
		rpc.WithPort(s.Agent.RPC.Port),
		rpc.WithNetworkInterface(s.Agent.RPC.NetInt),
		rpc.WithResourceLimits(s.Agent.RPC.Limits),
	}
	tlsConf := s.Agent.RPC.TLS
	if tlsConf.Enabled {
		if err := expandTLS(tlsConf); err != nil {
			return nil, err
		}
		opts = append(opts, rpc.WithTLS(rpc.ServerTLSConfig{
			Cert:             tlsConf.cert,
			PrivateKey:       tlsConf.key,
			IncludeSystemCAs: tlsConf.SystemCA,
			CustomCAs:        tlsConf.customCAs,
		}))
	}
	return opts, nil
}

// Gateway returns the configuration options available to set up an HTTP gateway.
func (s *Settings) Gateway() []rpc.GatewayOption {
	// gateway internal client options
	clOpts := []rpc.ClientOption{}
	if s.Agent.RPC.TLS.Enabled {
		clOpts = append(clOpts, rpc.WithClientTLS(rpc.ClientTLSConfig{
			IncludeSystemCAs: s.Agent.RPC.TLS.SystemCA,
			CustomCAs:        s.Agent.RPC.TLS.customCAs,
		}))
	}

	return []rpc.GatewayOption{
		rpc.WithClientOptions(clOpts...),
		rpc.WithHandlerName("http-gateway"),
		rpc.WithPrettyJSON("json+pretty"),
		rpc.WithGatewayMiddleware(mwHeaders.Handler(map[string]string{
			"x-didctl-version": info.CoreVersion,
			"x-didctl-build":   info.BuildCode,
		})),
	}
}

// ClientRPC returns the configuration options available to establish a connection
// to an agent RPC server.
func (s *Settings) ClientRPC() ([]rpc.ClientOption, error) {
	timeout := viper.GetInt("client.timeout")
	opts := []rpc.ClientOption{
		rpc.WaitForReady(),
		rpc.WithCompression(),
		rpc.WithTimeout(time.Duration(timeout) * time.Second),
		rpc.WithUserAgent(fmt.Sprintf("didctl-client/%s", info.CoreVersion)),
	}
	// apply server name override (dev/testing only)
	if override := viper.GetString("client.override"); override != "" {
		opts = append(opts, rpc.WithServerNameOverride(override))
	}

	// apply TLS settings
	tlsConf := s.Client.TLS
	if tlsConf.Enabled {
		if err := expandTLS(tlsConf); err != nil {
			return nil, err
		}
		// client TLS configuration
		opts = append(opts, rpc.WithClientTLS(rpc.ClientTLSConfig{
			IncludeSystemCAs: tlsConf.SystemCA,
			CustomCAs:        tlsConf.customCAs,
		}))
		// client TLS credentials
		if tlsConf.cert != nil && tlsConf.key != nil {
			opts = append(opts, rpc.WithAuthCertificate(tlsConf.cert, tlsConf.key))
		}
	}
	return opts, nil
}

// SetDefaults loads default values to the provided viper instance.
func (s *Settings) SetDefaults(v *viper.Viper, home, appID string) {
	v.SetDefault("client.timeout", 5)
	v.SetDefault("client.home", filepath.Join(home, fmt.Sprintf(".%s", appID)))
	v.SetDefault("resolver", []*resolver.Provider{
		{
			Method:   "bryk",
			Endpoint: "https://did.bryk.io/v1/retrieve/{{.Method}}/{{.Subject}}",
			Protocol: "http",
		},
	})
}

// ReleaseCode returns the release identifier for the application. A release
// identifier is of the form: `service-name@version+commit_hash`. If `version`
// or `commit_hash` are not available will be omitted.
func (s *Settings) ReleaseCode() string {
	// use service name
	release := serviceName

	// attach version tag. manually set value by default but prefer the one set
	// at build time if available
	version := s.Agent.OTEL.ServiceVersion
	if strings.Count(info.CoreVersion, ".") >= 2 {
		version = info.CoreVersion
	}
	if version != "" {
		release = fmt.Sprintf("%s@%s", release, version)
	}

	// attach commit hash if available
	if info.BuildCode != "" {
		release = fmt.Sprintf("%s+%s", release, info.BuildCode)
	}
	return release
}

// Configuration settings available when running an agent instance.
type agent struct {
	PoW     uint          `json:"pow" yaml:"pow" mapstructure:"pow"`
	Storage string        `json:"storage" yaml:"storage" mapstructure:"storage"`
	Methods []string      `json:"methods" yaml:"methods" mapstructure:"methods"`
	OTEL    *otelSettings `json:"otel" yaml:"otel" mapstructure:"otel"`
	RPC     *rpcSettings  `json:"rpc" yaml:"rpc" mapstructure:"rpc"`
}

// Configuration settings available when running a client instance.
type client struct {
	Node    string       `json:"node" yaml:"node" mapstructure:"node"`
	Timeout uint         `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	PoW     uint         `json:"pow" yaml:"pow" mapstructure:"pow"`
	TLS     *tlsSettings `json:"tls" yaml:"tls" mapstructure:"tls"`
}

type otelSettings struct {
	ServiceName    string                 `json:"service_name" yaml:"service_name" mapstructure:"service_name"`
	ServiceVersion string                 `json:"service_version" yaml:"service_version" mapstructure:"service_version"`
	Collector      string                 `json:"collector" yaml:"collector" mapstructure:"collector"`
	LogJSON        bool                   `json:"log_json" yaml:"log_json" mapstructure:"log_json"`
	Attributes     map[string]interface{} `json:"attributes" yaml:"attributes" mapstructure:"attributes"`
	Sentry         *sentry.Options        `json:"sentry" yaml:"sentry" mapstructure:"sentry"`
}

type rpcSettings struct {
	NetInt string             `json:"network_interface" yaml:"network_interface" mapstructure:"network_interface"`
	Port   int                `json:"port" yaml:"port" mapstructure:"port"`
	HTTP   bool               `json:"http" yaml:"http" mapstructure:"http"`
	Limits rpc.ResourceLimits `json:"resource_limits" yaml:"resource_limits" mapstructure:"resource_limits"`
	TLS    *tlsSettings       `json:"tls" yaml:"tls" mapstructure:"tls"`
}

type tlsSettings struct {
	Enabled  bool     `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	SystemCA bool     `json:"system_ca" yaml:"system_ca" mapstructure:"system_ca"`
	Cert     string   `json:"cert" yaml:"cert" mapstructure:"cert"`
	Key      string   `json:"key" yaml:"key" mapstructure:"key"`
	CustomCA []string `json:"custom_ca" yaml:"custom_ca" mapstructure:"custom_ca"`

	// private expanded values
	cert      []byte
	key       []byte
	customCAs [][]byte
}
