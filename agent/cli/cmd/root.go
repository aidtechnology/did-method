package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bryk-io/did-method/agent"
	"github.com/bryk-io/x/net/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "bryk-id-node",
	Short: "Starts a new node supporting the DID method requirements",
	RunE:  runMethodServer,
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	params := []cParam{
		{
			name:      "port",
			usage:     "TCP port to use for the server",
			flagKey:   "server.port",
			byDefault: 9090,
		},
		{
			name:      "storage",
			usage:     "specify the directory to use for data storage",
			flagKey:   "server.storage",
			byDefault: "/etc/bryk-did/agent",
		},
	}
	if err := setupCommandParams(rootCmd, params); err != nil {
		log.Fatal(err)
	}
}

func runMethodServer(_ *cobra.Command, _ []string) error {
	port := viper.GetInt("server.port")
	storage := viper.GetString("server.storage")
	fmt.Printf("Starting node. Port: %d - Storage: %s\n", port, storage)
	handler, err := agent.NewHandler(storage)
	if err != nil {
		return fmt.Errorf("failed to start method handler: %s", err)
	}

	var opts []rpc.ServerOption
	opts = append(opts, rpc.WithPort(port))
	opts = append(opts, rpc.WithHTTPGateway(rpc.HTTPGatewayOptions{
		Port:         port,
		EmitDefaults: false,
	}))
	server, err := handler.GetServer(opts...)
	if err != nil {
		return fmt.Errorf("failed to start node: %s", err)
	}
	go server.Start()

	fmt.Println("Waiting for requests...")
	<-signalsHandler()
	fmt.Println("Preparing to exit")
	err = handler.Close()
	if !strings.Contains(err.Error(), "closed network connection") {
		return err
	}
	return nil
}

// Custom OS signals handler
func signalsHandler() chan os.Signal {
	signalsCh := make(chan os.Signal, 1)
	signalList := []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt}
	signal.Reset(signalList...)
	signal.Notify(signalsCh, signalList...)
	return signalsCh
}
