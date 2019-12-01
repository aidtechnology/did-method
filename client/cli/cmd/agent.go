package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/bryk-io/did-method/agent"
	"github.com/bryk-io/x/cli"
	"github.com/bryk-io/x/net/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	}
	if err := cli.SetupCommandParams(agentCmd, params); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(agentCmd)
}

func runMethodServer(_ *cobra.Command, _ []string) error {
	port := viper.GetInt("server.port")
	storage := viper.GetString("server.storage")
	handler, err := agent.NewHandler(storage, uint(viper.GetInt("server.pow")))
	if err != nil {
		return fmt.Errorf("failed to start method handler: %s", err)
	}

	handler.Log("starting network agent")
	handler.Log(fmt.Sprintf("TCP port: %d", port))
	handler.Log(fmt.Sprintf("storage directory: %s", storage))
	handler.Log(fmt.Sprintf("difficulty level: %d", viper.GetInt("server.pow")))

	gw, err := rpc.NewHTTPGateway(
		rpc.WithEncoder("application/json", rpc.MarshalerStandard(false)),
		rpc.WithEncoder("application/json+pretty", rpc.MarshalerStandard(true)),
		rpc.WithEncoder("*", rpc.MarshalerStandard(false)))
	if err != nil {
		return fmt.Errorf("failed to initialize HTTP gateway: %s", err)
	}
	var opts []rpc.ServerOption
	opts = append(opts, rpc.WithPanicRecovery())
	opts = append(opts, rpc.WithPort(port))
	opts = append(opts, rpc.WithNetworkInterface(rpc.NetworkInterfaceAll))
	opts = append(opts, rpc.WithHTTPGateway(gw))

	// Start server and wait for it to be ready
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
	handler.Log("preparing to exit")
	err = handler.Close()
	if err != nil && !strings.Contains(err.Error(), "closed network connection") {
		return err
	}
	return nil
}
