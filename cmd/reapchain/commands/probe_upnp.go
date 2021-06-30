package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	tmjson "gitlab.reappay.net/reapchain/reapchain-core/libs/json"
	"gitlab.reappay.net/reapchain/reapchain-core/p2p/upnp"
)

// ProbeUpnpCmd adds capabilities to test the UPnP functionality.
var ProbeUpnpCmd = &cobra.Command{
	Use:     "probe-upnp",
	Aliases: []string{"probe_upnp"},
	Short:   "Test UPnP functionality",
	RunE:    probeUpnp,
	PreRun:  deprecateSnakeCase,
}

func probeUpnp(cmd *cobra.Command, args []string) error {
	capabilities, err := upnp.Probe(logger)
	if err != nil {
		fmt.Println("Probe failed: ", err)
	} else {
		fmt.Println("Probe success!")
		jsonBytes, err := tmjson.Marshal(capabilities)
		if err != nil {
			return err
		}
		fmt.Println(string(jsonBytes))
	}
	return nil
}
