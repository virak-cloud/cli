package presenter

import (
	"fmt"
	"github.com/virak-cloud/cli/pkg/http/responses"
	"os"

	"github.com/olekukonko/tablewriter"
)

// RenderBucketList renders a table of object storage buckets.
func RenderBucketList(buckets []responses.ObjectStorageBucket) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "URL", "Status", "Policy", "Size"})

	for _, bucket := range buckets {
		table.Append([]string{bucket.ID, bucket.Name, bucket.URL, bucket.Status, bucket.Policy, fmt.Sprintf("%d", bucket.Size)})
	}
	table.Render()
}

// RenderBucketDetail renders a table with the details of an object storage bucket.
func RenderBucketDetail(bucket responses.ObjectStorageBucket) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})
	table.Append([]string{"ID", bucket.ID})
	table.Append([]string{"Name", bucket.Name})
	table.Append([]string{"URL", bucket.URL})
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

// RenderBucketEvents renders a table of object storage events.
func RenderBucketEvents(events []responses.ObjectStorageEvent) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Bucket ID", "Source", "Type", "Content", "Created At"})

	for _, event := range events {
		table.Append([]string{event.ProductID, event.ProductSource, event.Type, event.Content, fmt.Sprintf("%d", event.CreatedAt)})
	}
	table.Render()
}

