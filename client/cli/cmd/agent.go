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

	"github.com/bryk-io/did-method/agent"
	"github.com/bryk-io/did-method/agent/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.bryk.io/x/cli"
	"go.bryk.io/x/net/rpc"
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
	if err := cli.SetupCommandParams(agentCmd, params); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(agentCmd)
}

func runMethodServer(_ *cobra.Command, _ []string) error {
	// Prepare API handler
	handler, err := getAgentHandler()
	if err != nil {
		return err
	}

	// CPU profile
	if viper.GetBool("server.debug") {
		cpu, err := ioutil.TempFile("", "didctl_cpu_")
		if err != nil {
			return err
		}
		if err := pprof.StartCPUProfile(cpu); err != nil {
			return err
		}
		defer func() {
			log.Infof("CPU profile saved at %s", cpu.Name())
			pprof.StopCPUProfile()
			_ = cpu.Close()
		}()
	}

	// Base server configuration
	opts := []rpc.ServerOption{
		rpc.WithPanicRecovery(),
		rpc.WithPort(viper.GetInt("server.port")),
		rpc.WithNetworkInterface(rpc.NetworkInterfaceAll),
		rpc.WithService(handler.ServiceDefinition()),
		rpc.WithLogger(rpc.LoggingOptions{
		 	Logger: log,
		 	IncludePayload: false,
		 	FilterMethods: []string{
		 		"bryk.did.proto.v1.AgentAPI/Ping",
		  },
		}),
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

	// Monitoring
	if viper.GetBool("server.http") && viper.GetBool("server.monitoring") {
		log.Info("monitoring enabled")
		opts = append(opts, rpc.WithMonitoring(rpc.MonitoringOptions{
			IncludeHistograms:   true,
			UseProcessCollector: true,
			UseGoCollector:      true,
		}))
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
		return fmt.Errorf("failed to start node: %s", err)
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

func getAgentHandler() (*agent.Handler, error) {
	// Get handler settings
	methods := viper.GetStringSlice("server.method")
	pow := uint(viper.GetInt("server.pow"))
	store, err := getStorage(viper.GetString("server.storage"))
	if err != nil {
		return nil, err
	}

	// Prepare API handler
	handler, err := agent.NewHandler(log, methods, pow, store)
	if err != nil {
		return nil, fmt.Errorf("failed to start method handler: %s", err)
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
		return nil, fmt.Errorf("failed to load certificate file: %s", err)
	}
	tlsConf.PrivateKey, err = ioutil.ReadFile(viper.GetString("server.tls.key"))
	if err != nil {
		return nil, fmt.Errorf("failed to load private key file: %s", err)
	}
	if viper.GetString("server.tls.ca") != "" {
		caPEM, err := ioutil.ReadFile(viper.GetString("server.tls.ca"))
		if err != nil {
			return nil, fmt.Errorf("failed to load CA file: %s", err)
		}
		tlsConf.CustomCAs = append(tlsConf.CustomCAs, caPEM)
	}
	return rpc.WithTLS(tlsConf), nil
}

func getAgentGateway(handler *agent.Handler) (*rpc.HTTPGateway, error) {
	gwCl := []rpc.ClientOption{rpc.WaitForReady()}
	if viper.GetBool("server.tls.enabled") {
		tlsConf := rpc.ClientTLSConfig{IncludeSystemCAs: true}
		if viper.GetString("server.tls.ca") != "" {
			caPEM, err := ioutil.ReadFile(viper.GetString("server.tls.ca"))
			if err != nil {
				return nil, fmt.Errorf("failed to load CA file: %s", err)
			}
			tlsConf.CustomCAs = append(tlsConf.CustomCAs, caPEM)
		}
		gwCl = append(gwCl, rpc.WithClientTLS(tlsConf))
		gwCl = append(gwCl, rpc.WithInsecureSkipVerify()) // Internally the gateway proxy accept any certificate
	}

	// Properly adjust outgoing headers
	headersMatcher := func(h string) (string, bool) {
		if strings.HasPrefix(strings.ToLower(h), "x-") {
			return h, true
		}
		return "x-rpc-" + h, true
	}

	gwOpts := []rpc.HTTPGatewayOption{
		rpc.WithEncoder("application/json", rpc.MarshalerStandard(false)),
		rpc.WithEncoder("application/json+pretty", rpc.MarshalerStandard(true)),
		rpc.WithEncoder("*", rpc.MarshalerStandard(false)),
		rpc.WithClientOptions(gwCl),
		rpc.WithOutgoingHeaderMatcher(headersMatcher),
		rpc.WithFilter(handler.QueryResponseFilter()),
	}
	gw, err := rpc.NewHTTPGateway(gwOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize HTTP gateway: %s", err)
	}
	return gw, nil
}

// Return the proper storage handler instance based on the connection
// details provided.
func getStorage(info string) (agent.Storage, error) {
	switch {
	case info == "":
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
