package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/aidtechnology/did-method/agent"
	"github.com/aidtechnology/did-method/agent/storage"
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
	if err := cli.SetupCommandParams(agentCmd, conf.Overrides("agent"), viper.GetViper()); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(agentCmd)
}

func runMethodServer(_ *cobra.Command, _ []string) error {
	// Observability operator
	oop, err := otel.NewOperator(conf.OTEL(log)...)
	if err != nil {
		return err
	}

	// Prepare API handler
	handler, err := getAgentHandler(oop)
	if err != nil {
		return err
	}

	// Base server configuration
	opts, err := conf.Server(oop)
	if err != nil {
		return err
	}

	// Register handler as RPC service
	opts = append(opts, rpc.WithServiceProvider(handler))

	// Initialize HTTP gateway
	if conf.Agent.RPC.HTTP {
		log.Info("HTTP gateway enabled")
		gwOpts := conf.Gateway(oop)
		gwOpts = append(gwOpts, handler.CustomGatewayOptions()...)
		gw, err := rpc.NewGateway(gwOpts...)
		if err != nil {
			return err
		}
		opts = append(opts, rpc.WithHTTPGateway(gw))
	}

	// Start server and wait for it to be ready
	log.Infof("difficulty level: %d", conf.Agent.PoW)
	log.Infof("TCP port: %d", conf.Agent.RPC.Port)
	log.Info("starting network agent")
	if conf.Agent.RPC.TLS.Enabled {
		log.Infof("certificate: %s", conf.Agent.RPC.TLS.Cert)
		log.Infof("private key: %s", conf.Agent.RPC.TLS.Key)
	}
	server, err := rpc.NewServer(opts...)
	if err != nil {
		return err
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
	return handler.Close()
}

// Return an API handler instance.
func getAgentHandler(oop *otel.Operator) (*agent.Handler, error) {
	// storage handler
	store, err := getStorage(conf.Agent.Storage)
	if err != nil {
		return nil, err
	}

	// API handler
	methods := conf.Agent.Methods
	pow := conf.Agent.PoW
	handler, err := agent.NewHandler(methods, pow, store, oop)
	if err != nil {
		return nil, fmt.Errorf("failed to start method handler: %w", err)
	}
	log.Infof("storage: %s", store.Description())
	return handler, nil
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
