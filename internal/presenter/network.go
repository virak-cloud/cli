package presenter

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/virak-cloud/cli/pkg/http/responses"
)

func RenderNetworkDetail(network responses.Network) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})
	table.Append([]string{"ID", network.ID})
	table.Append([]string{"Name", network.Name})
	table.Append([]string{"Status", network.Status})
	table.Append([]string{"Network Offering ID", network.NetworkOffering.ID})
	table.Append([]string{"Network Offering Name", network.NetworkOffering.Name})
	table.Render()
}

func RenderNetworkList(networks []responses.Network) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Status", "Offering Name"})

	for _, network := range networks {
		table.Append([]string{network.ID, network.Name, network.Status, network.NetworkOffering.Name})
	}
	table.Render()
}

func RenderInstanceNetworkList(instances []responses.InstanceNetwork) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Instance Network ID", "Instance ID", "IP Address", "Network Name", "Is Default"})

	for _, instance := range instances {
		isDefault := "No"
		if instance.IsDefault {
			isDefault = "Yes"
		}
		table.Append([]string{instance.ID, instance.InstanceID, instance.IPAddress, instance.Network.Name, isDefault})
	}
	table.Render()
}

func RenderPortForwardList(rules []responses.PortForwardRule) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Protocol", "Public Port", "Private Port", "Private IP", "Status", "Created At"})

	for _, rule := range rules {
		table.Append([]string{rule.ID, rule.Protocol, fmt.Sprintf("%d", rule.PublicPort), fmt.Sprintf("%d", rule.PrivatePort), rule.PrivateIP, rule.Status, fmt.Sprintf("%d", rule.CreatedAt)})
	}
	table.Render()
}
