package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"syscall"

	"github.com/aidtechnology/did-method/agent"
	"github.com/aidtechnology/did-method/agent/storage"
	"github.com/aidtechnology/did-method/info"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.bryk.io/pkg/cli"
	"go.bryk.io/pkg/net/rpc"
	"go.bryk.io/pkg/otel"
)

var agentCmd = &cobra.Command{
	Use:     "agent",
	Short:   "Start a network agent supporting the DID method requirements",
	Example: "didctl agent --port 8080",
	Aliases: []string{"server", "node"},
	RunE:    runMethodServer,
}

func init() {
	params := []cli.Param{
		{
			Name:      "port",
			Usage:     "TCP port to use for the server",
			FlagKey:   "server.port",
			ByDefault: 9090,
			Short:     "p",
		},
		{
			Name:      "pow",
			Usage:     "set the required request ticket difficulty level",
			FlagKey:   "server.pow",
			ByDefault: 24,
		},
		{
			Name:      "http",
			Usage:     "enable the HTTP interface",
			FlagKey:   "server.http",
			ByDefault: false,
		},
		{
			Name:      "monitoring",
			Usage:     "publish metrics for instrumentation and monitoring",
			FlagKey:   "server.monitoring",
			ByDefault: false,
		},
		{
			Name:      "debug",
			Usage:     "run agent in debug mode to generate profiling information",
			FlagKey:   "server.debug",
			ByDefault: false,
		},
		{
			Name:      "tls",
			Usage:     "enable secure communications using TLS with provided credentials",
			FlagKey:   "server.tls.enabled",
			ByDefault: false,
		},
		{
			Name:      "tls-ca",
			Usage:     "TLS custom certificate authority (path to PEM file)",
			FlagKey:   "server.tls.ca",
			ByDefault: "",
		},
		{
			Name:      "tls-cert",
			Usage:     "TLS certificate (path to PEM file)",
			FlagKey:   "server.tls.cert",
			ByDefault: "/etc/didctl/agent/tls.crt",
		},
		{
			Name:      "tls-key",
			Usage:     "TLS private key (path to PEM file)",
			FlagKey:   "server.tls.key",
			ByDefault: "/etc/didctl/agent/tls.key",
		},
		{
			Name:      "method",
			Usage:     "specify a supported DID method (can be provided multiple times)",
			FlagKey:   "server.method",
			ByDefault: []string{"bryk"},
			Short:     "m",
		},
		{
			Name:      "storage",
			Usage:     "specify storage mechanism to use",
			FlagKey:   "server.storage",
			ByDefault: "",
			Short:     "s",
		},
	}
	if err := cli.SetupCommandParams(agentCmd, params, viper.GetViper()); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(agentCmd)
}

func runMethodServer(_ *cobra.Command, _ []string) error {
	// CPU profile
	if viper.GetBool("server.debug") {
		cpu, err := ioutil.TempFile("", "didctl_cpu_")
		if err != nil {
			return err
		}
		if err := pprof.StartCPUProfile(cpu); err != nil {
			return err
		}
		// nolint: gosec
		defer func() {
			log.Infof("CPU profile saved at %s", cpu.Name())
			pprof.StopCPUProfile()
			_ = cpu.Close()
		}()
	}

	// Observability operator
	oop, err := otel.NewOperator([]otel.OperatorOption{
		otel.WithLogger(log),
		otel.WithServiceName("didctl"),
		otel.WithServiceVersion(info.CoreVersion),
	}...)
	if err != nil {
		return err
	}

	// Prepare API handler
	handler, err := getAgentHandler(oop)
	if err != nil {
		return err
	}

	// Base server configuration
	opts := []rpc.ServerOption{
		rpc.WithPanicRecovery(),
		rpc.WithPort(viper.GetInt("server.port")),
		rpc.WithNetworkInterface(rpc.NetworkInterfaceAll),
		rpc.WithServiceProvider(handler),
		rpc.WithObservability(oop),
	}

	// TLS configuration
	if viper.GetBool("server.tls.enabled") {
		log.Info("TLS enabled")
		opt, err := loadAgentCredentials()
		if err != nil {
			return err
		}
		opts = append(opts, opt)
	}

	// Initialize HTTP gateway
	if viper.GetBool("server.http") {
		log.Info("HTTP gateway available")
		gw, err := getAgentGateway(handler)
		if err != nil {
			return err
		}
		opts = append(opts, rpc.WithHTTPGateway(gw))
	}

	// Start server and wait for it to be ready
	log.Infof("difficulty level: %d", viper.GetInt("server.pow"))
	log.Infof("TCP port: %d", viper.GetInt("server.port"))
	log.Info("starting network agent")
	if viper.GetBool("server.tls.enabled") {
		log.Infof("certificate: %s", viper.GetString("server.tls.cert"))
		log.Infof("private key: %s", viper.GetString("server.tls.key"))
	}
	server, err := rpc.NewServer(opts...)
	if err != nil {
		return fmt.Errorf("failed to start node: %w", err)
	}
	ready := make(chan bool)
	go func() {
		_ = server.Start(ready)
	}()
	<-ready

	// Wait for system signals
	log.Info("waiting for incoming requests")
	<-cli.SignalsHandler([]os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	})

	// Close handler
	log.Info("preparing to exit")
	if err = handler.Close(); err != nil && !strings.Contains(err.Error(), "closed network connection") {
		return err
	}

	// Dump memory profile and exit
	if viper.GetBool("server.debug") {
		// Memory profile
		mem, err := ioutil.TempFile("", "didctl_mem_")
		if err != nil {
			return err
		}
		runtime.GC()
		if err := pprof.WriteHeapProfile(mem); err != nil {
			return err
		}
		log.Infof("memory profile saved at %s", mem.Name())
		_ = mem.Close()
	}
	return nil
}

func getAgentHandler(oop *otel.Operator) (*agent.Handler, error) {
	// Get handler settings
	methods := viper.GetStringSlice("server.method")
	pow := uint(viper.GetInt("server.pow"))
	store, err := getStorage(viper.GetString("server.storage"))
	if err != nil {
		return nil, err
	}

	// Prepare API handler
	handler, err := agent.NewHandler(methods, pow, store, oop)
	if err != nil {
		return nil, fmt.Errorf("failed to start method handler: %w", err)
	}
	log.Infof("storage: %s", store.Description())
	return handler, nil
}

func loadAgentCredentials() (rpc.ServerOption, error) {
	var err error
	tlsConf := rpc.ServerTLSConfig{
		IncludeSystemCAs: true,
	}
	tlsConf.Cert, err = ioutil.ReadFile(viper.GetString("server.tls.cert"))
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate file: %w", err)
	}
	tlsConf.PrivateKey, err = ioutil.ReadFile(viper.GetString("server.tls.key"))
	if err != nil {
		return nil, fmt.Errorf("failed to load private key file: %w", err)
	}
	if viper.GetString("server.tls.ca") != "" {
		caPEM, err := ioutil.ReadFile(viper.GetString("server.tls.ca"))
		if err != nil {
			return nil, fmt.Errorf("failed to load CA file: %w", err)
		}
		tlsConf.CustomCAs = append(tlsConf.CustomCAs, caPEM)
	}
	return rpc.WithTLS(tlsConf), nil
}

func getAgentGateway(handler *agent.Handler) (*rpc.Gateway, error) {
	gwCl := []rpc.ClientOption{}
	if viper.GetBool("server.tls.enabled") {
		tlsConf := rpc.ClientTLSConfig{IncludeSystemCAs: true}
		if viper.GetString("server.tls.ca") != "" {
			caPEM, err := ioutil.ReadFile(viper.GetString("server.tls.ca"))
			if err != nil {
				return nil, fmt.Errorf("failed to load CA file: %w", err)
			}
			tlsConf.CustomCAs = append(tlsConf.CustomCAs, caPEM)
		}
		gwCl = append(gwCl, rpc.WithClientTLS(tlsConf))
		gwCl = append(gwCl, rpc.WithInsecureSkipVerify()) // Internally the gateway proxy accept any certificate
	}

	gwOpts := []rpc.GatewayOption{
		rpc.WithClientOptions(gwCl...),
		rpc.WithInterceptor(handler.QueryResponseFilter()),
	}
	gw, err := rpc.NewGateway(gwOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize HTTP gateway: %w", err)
	}
	return gw, nil
}

// Return the proper storage handler instance based on the connection
// details provided.
func getStorage(info string) (agent.Storage, error) {
	switch {
	case info == "ephemeral":
		store := &storage.Ephemeral{}
		_ = store.Open("no-op")
		return store, nil
	case strings.HasPrefix(info, "mongodb://"):
		store := &storage.MongoStore{}
		return store, store.Open(info)
	default:
		return nil, errors.New("non supported storage")
	}
}
