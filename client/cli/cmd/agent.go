package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"syscall"

	"github.com/bryk-io/did-method/agent"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.bryk.io/x/cli"
	"go.bryk.io/x/net/rpc"
)

var agentCmd = &cobra.Command{
	Use:           "agent",
	Short:         "Starts a new network agent supporting the DID method requirements",
	Example:       "didctl agent --storage /var/run/didctl --port 8080",
	Aliases:       []string{"server", "node"},
	RunE:          runMethodServer,
}

func init() {
	params := []cli.Param{
		{
			Name:      "port",
			Usage:     "TCP port to use for the server",
			FlagKey:   "server.port",
			ByDefault: 9090,
		},
		{
			Name:      "storage",
			Usage:     "specify the directory to use for data storage",
			FlagKey:   "server.storage",
			ByDefault: "/etc/didctl/agent",
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
			FlagKey:   "server.tls.enable",
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
			ByDefault: "",
		},
		{
			Name:      "tls-key",
			Usage:     "TLS private key (path to PEM file)",
			FlagKey:   "server.tls.key",
			ByDefault: "",
		},
	}
	if err := cli.SetupCommandParams(agentCmd, params); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(agentCmd)
}

func runMethodServer(_ *cobra.Command, _ []string) error {
	// Prepare API handler
	storage := viper.GetString("server.storage")
	handler, err := agent.NewHandler(storage, uint(viper.GetInt("server.pow")))
	if err != nil {
		return fmt.Errorf("failed to start method handler: %s", err)
	}
	handler.Log(fmt.Sprintf("storage directory: %s", storage))

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
			handler.Log(fmt.Sprintf("CPU profile saved at %s", cpu.Name()))
			pprof.StopCPUProfile()
			_ = cpu.Close()
		}()
	}

	// Base server configuration
	opts := []rpc.ServerOption{
		rpc.WithPanicRecovery(),
		rpc.WithPort(viper.GetInt("server.port")),
		rpc.WithNetworkInterface(rpc.NetworkInterfaceAll),
	}

	// TLS configuration
	if viper.GetBool("server.tls.enable") {
		handler.Log("TLS enabled")
		opt, err := loadAgentCredentials()
		if err != nil {
			return err
		}
		opts = append(opts, opt)
	}

	// Initialize HTTP gateway
	if viper.GetBool("server.http") {
		handler.Log("HTTP gateway available")
		gw, err := getAgentGateway()
		if err != nil {
			return err
		}
		opts = append(opts, rpc.WithHTTPGateway(gw))
	}

	// Monitoring
	if viper.GetBool("server.http") && viper.GetBool("server.monitoring") {
		handler.Log("monitoring enabled")
		opts = append(opts, rpc.WithMonitoring(rpc.MonitoringOptions{
			IncludeHistograms:   true,
			UseProcessCollector: true,
			UseGoCollector:      true,
		}))
	}

	// Start server and wait for it to be ready
	handler.Log(fmt.Sprintf("difficulty level: %d", viper.GetInt("server.pow")))
	handler.Log(fmt.Sprintf("TCP port: %d", viper.GetInt("server.port")))
	handler.Log("starting network agent")
	if viper.GetBool("server.tls.enable") {
		handler.Log(fmt.Sprintf("certificate: %s", viper.GetString("server.tls.cert")))
		handler.Log(fmt.Sprintf("private key: %s", viper.GetString("server.tls.key")))
	}
	server, err := handler.GetServer(opts...)
	if err != nil {
		return fmt.Errorf("failed to start node: %s", err)
	}
	ready := make(chan bool)
	go func() {
		_ = server.Start(ready)
	}()
	<-ready

	// Wait for system signals
	handler.Log("waiting for incoming requests")
	<-cli.SignalsHandler([]os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	})

	// Close handler
	handler.Log("preparing to exit")
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
		handler.Log(fmt.Sprintf("memory profile saved at %s", mem.Name()))
		_ = mem.Close()
	}
	return nil
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

func getAgentGateway() (*rpc.HTTPGateway, error) {
	gwCl := []rpc.ClientOption{rpc.WaitForReady()}
	if viper.GetBool("server.tls.enable") {
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
		if strings.HasPrefix(strings.ToLower(h),"x-") {
			return h, true
		} else {
			return "x-rpc-metadata-" + h, true
		}
	}

	gwOpts := []rpc.HTTPGatewayOption{
		rpc.WithEncoder("application/json", rpc.MarshalerStandard(false)),
		rpc.WithEncoder("application/json+pretty", rpc.MarshalerStandard(true)),
		rpc.WithEncoder("*", rpc.MarshalerStandard(false)),
		rpc.WithClientOptions(gwCl),
		rpc.WithOutgoingHeaderMatcher(headersMatcher),
	}
	gw, err := rpc.NewHTTPGateway(gwOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize HTTP gateway: %s", err)
	}
	return gw, nil
}
