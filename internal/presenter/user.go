package presenter

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/virak-cloud/cli/pkg/http/responses"
)

func RenderTokenAbilities(abilities []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Ability"})
	for _, ability := range abilities {
		table.Append([]string{ability})
	}
	table.Render()
}

func RenderSSHKeyList(sshKeys []responses.UserSSHKey) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Key", "Created At"})
	for _, key := range sshKeys {
		keyData := key.DataValue
		if len(keyData) > 10 {
			keyData = keyData[:17] + "..."
		}
		table.Append([]string{key.ID, key.DisplayName, keyData, key.CreatedAt})
	}
	table.Render()
}

func RenderUserProfile(profile *responses.UserProfileResponse) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})
	table.SetAutoWrapText(false)

	data := profile.Data

	// Add rows for each field
	table.Append([]string{"ID", data.ID})
	table.Append([]string{"Name", data.Name})
	table.Append([]string{"Language", data.Language})
	table.Append([]string{"National Code", data.NationalCode})
	table.Append([]string{"Email", data.Email})
	table.Append([]string{"Phone", data.Phone})
	table.Append([]string{"Country", formatInterface(data.Country)})
	table.Append([]string{"State", formatInterface(data.State)})
	table.Append([]string{"City", formatInterface(data.City)})
	table.Append([]string{"Address", formatInterface(data.Address)})
	table.Append([]string{"ZIP", formatInterface(data.Zip)})
	table.Append([]string{"Website", formatInterface(data.Website)})
	table.Append([]string{"Referral Code", formatInterface(data.Extra.ReferralCode)})
	table.Append([]string{"Status", data.Status})
	table.Append([]string{"Type", data.Type})
	table.Append([]string{"Created At", data.CreatedAt})
	table.Append([]string{"Updated At", data.UpdatedAt})
	table.Append([]string{"Customer Zones Count", formatInterface(data.CustomerZonesCount)})
	table.Append([]string{"Instances Count", formatInterface(data.InstancesCount)})
	table.Append([]string{"Payments Count", formatInterface(data.PaymentsCount)})
	table.Append([]string{"Wallets Count", formatInterface(data.WalletsCount)})
	table.Append([]string{"Invite Code", data.InviteCode})
	table.Append([]string{"Invited By Me", fmt.Sprintf("%d", data.InvitedByMe)})
	table.Append([]string{"Picture", data.Picture})

	table.Render()
}

// Helper function to format interface{} values (handles null values)
func formatInterface(value interface{}) string {
	if value == nil {
		return "N/A"
	}
	return fmt.Sprintf("%v", value)
}
