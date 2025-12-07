package presenter

import (
	"fmt"
	"os"

	"github.com/virak-cloud/cli/pkg/http/responses"

	"github.com/olekukonko/tablewriter"
)

func RenderBucketList(buckets []responses.ObjectStorageBucket) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "URL", "Region", "Status", "Policy", "Size"})

	for _, bucket := range buckets {
		table.Append([]string{
			bucket.ID,
			bucket.Name,
			bucket.URL,
			bucket.Region,
			bucket.Status,
			bucket.Policy,
			fmt.Sprintf("%d", bucket.Size),
		})
	}
	table.Render()
}

func RenderBucketDetail(bucket responses.ObjectStorageBucket) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})
	table.Append([]string{"ID", bucket.ID})
	table.Append([]string{"Name", bucket.Name})
	table.Append([]string{"URL", bucket.URL})
	table.Append([]string{"Region", bucket.Region})
	table.Append([]string{"Access Key", bucket.AccessKey})
	table.Append([]string{"Secret Key", bucket.SecretKey})
	table.Append([]string{"Status", bucket.Status})
	table.Append([]string{"Policy", bucket.Policy})
	table.Append([]string{"Size", fmt.Sprintf("%d", bucket.Size)})
	table.Append([]string{"Created At", fmt.Sprintf("%d", bucket.CreatedAt)})
	table.Append([]string{"Updated At", fmt.Sprintf("%d", bucket.UpdatedAt)})
	table.Append([]string{"Tier", bucket.Tier})
	table.Append([]string{"Is Failed", fmt.Sprintf("%t", bucket.IsFailed)})
	table.Append([]string{"Message", bucket.Message})
	table.Render()
}

func RenderBucketEvents(events []responses.ObjectStorageEvent) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Bucket ID", "Source", "Type", "Content", "Created At"})

	for _, event := range events {
		table.Append([]string{event.ProductID, event.ProductSource, event.Type, event.Content, fmt.Sprintf("%d", event.CreatedAt)})
	}
	table.Render()
}
