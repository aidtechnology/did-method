package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/bryk-io/x/did"
	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addServiceCmd = &cobra.Command{
	Use:     "add",
	Short:   "Register a new service entry for the DID",
	Example: "bryk-id did service add [DID reference name] --name \"service name\" --endpoint https://www.agency.com/user_id",
	RunE:    runAddServiceCmd,
}

func init() {
	params := []cParam{
		{
			name:      "name",
			usage:     "service's reference name",
			flagKey:   "service-add.name",
			byDefault: "external-service-#",
		},
		{
			name:      "type",
			usage:     "type identifier for the service handler",
			flagKey:   "service-add.type",
			byDefault: "identity.bryk.io.ExternalService",
		},
		{
			name:      "endpoint",
			usage:     "main URL to access the service",
			flagKey:   "service-add.endpoint",
			byDefault: "",
		},
	}
	if err := setupCommandParams(addServiceCmd, params); err != nil {
		panic(err)
	}
	serviceCmd.AddCommand(addServiceCmd)
}

func runAddServiceCmd(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("you must specify a DID reference name")
	}
	if strings.TrimSpace(viper.GetString("service-add.endpoint")) == "" {
		return errors.New("service endpoint is required")
	}

	// Get store handler
	st, err := getClientStore()
	if err != nil {
		return err
	}
	defer st.Close()

	// Get identifier
	ll := getLogger()
	name := sanitize.Name(args[0])
	ll.Info("adding new service")
	ll.Debugf("retrieving entry with reference name: %s", name)
	e := st.Get(name)
	if e == nil {
		return fmt.Errorf("no available record under the provided reference name: %s", name)
	}
	id := &did.Identifier{}
	if err = id.Decode(e.Contents); err != nil {
		return errors.New("failed to decode entry contents")
	}

	// Validate service data
	ll.Debug("validating parameters")
	svc := &did.ServiceEndpoint{
		ID:       viper.GetString("service-add.name"),
		Type:     viper.GetString("service-add.type"),
		Endpoint: viper.GetString("service-add.endpoint"),
	}
	if strings.Count(svc.ID, "#") > 1 {
		return errors.New("invalid service name")
	}
	if strings.Count(svc.ID, "#") == 1 {
		svc.ID = strings.Replace(svc.ID, "#", fmt.Sprintf("%d", len(id.Services())+1), 1)
	}
	svc.ID = sanitize.Name(svc.ID)
	if _, err = url.ParseRequestURI(svc.Endpoint); err != nil {
		return fmt.Errorf("invalid service enpoint: %s", svc.Endpoint)
	}

	// Add service
	ll.Debugf("registering service with id: %s", svc.ID)
	if err = id.AddService(svc); err != nil {
		return fmt.Errorf("failed to add new service: %s", err)
	}

	// Update record
	ll.Info("updating local record")
	contents, err := id.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode identifier: %s", err)
	}
	return st.Update(name, contents)
}
