package presenter

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/virak-cloud/cli/pkg/http/responses"
)

// RenderWallet displays wallet balance information in a key-value format
func RenderWallet(wallet responses.WalletsBalanceResponse) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})
	table.Append([]string{"Name", wallet.Data.Name})
	table.Append([]string{"Track", wallet.Data.Track})
	table.Append([]string{"Type", wallet.Data.Type})
	table.Append([]string{"Balance", fmt.Sprintf("%.2f", wallet.Data.Balance)})
	table.Append([]string{"Balance Limit", fmt.Sprintf("%.2f", wallet.Data.BalanceLimit)})
	table.Append([]string{"Is Blocked", fmt.Sprintf("%t", wallet.Data.IsBlocked)})
	table.Append([]string{"Max Cost", fmt.Sprintf("%.2f", wallet.Data.MaxCost)})
	table.Append([]string{"Remaining Hours", fmt.Sprintf("%.2f", wallet.Data.RemainingHours)})
	table.Append([]string{"Updated At", wallet.Data.UpdatedAt})
	table.Render()
}

// RenderCostDocuments displays cost documents in a table format
func RenderCostDocuments(documents []responses.CostDocument) {
	if len(documents) == 0 {
		fmt.Println("No cost documents found.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Period", "Instance", "Network", "Snapshots", "Volumes", "Support", "Public IP", "Device", "Bucket Size", "Bucket Traffic", "Kubernetes"})

	for _, doc := range documents {
		period := fmt.Sprintf("%s to %s", doc.DateFrom, doc.DateTo)
		table.Append([]string{
			period,
			fmt.Sprintf("%.2f", doc.Instance),
			fmt.Sprintf("%.2f", doc.NetworkNetflow),
			fmt.Sprintf("%.2f", doc.InstanceSnapshot),
			fmt.Sprintf("%.2f", doc.InstanceDataVolumes),
			fmt.Sprintf("%.2f", doc.SupportOfferings),
			fmt.Sprintf("%.2f", doc.NetworkInternetPublicAddressV4),
			fmt.Sprintf("%.2f", doc.NetworkDevice),
			fmt.Sprintf("%.2f", doc.BucketSize),
			fmt.Sprintf("%.2f", doc.BucketDownloadTraffic+doc.BucketUploadTraffic),
			fmt.Sprintf("%.2f", doc.KubernetesNode),
		})
	}
	table.Render()
}

// RenderPayments displays payment history in a table format
func RenderPayments(payments []interface{}) {
	if len(payments) == 0 {
		fmt.Println("No payments found.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Amount", "Driver", "Status", "Reference ID", "Created At"})

	for _, payment := range payments {
		if paymentMap, ok := payment.(map[string]interface{}); ok {
			id := getStringValue(paymentMap, "id")
			amountStr := getStringValue(paymentMap, "amount")
			driver := getStringValue(paymentMap, "driver")
			status := getStringValue(paymentMap, "status")
			referenceID := getStringValue(paymentMap, "reference_id")
			createdAtStr := getStringValue(paymentMap, "created_at")

			// Format amount with commas
			var formattedAmount string
			if amount, err := strconv.ParseFloat(amountStr, 64); err == nil {
				formattedAmount = formatWithCommas(amount)
			} else {
				formattedAmount = amountStr
			}

			// Format created_at as human readable
			var formattedTime string
			if timestamp, err := strconv.ParseFloat(createdAtStr, 64); err == nil {
				t := time.Unix(int64(timestamp), 0)
				formattedTime = t.Format("2006-01-02 15:04:05")
			} else {
				formattedTime = createdAtStr
			}

			table.Append([]string{id, formattedAmount, driver, status, referenceID, formattedTime})
		} else {
			// Fallback to original format if it's not a map
			table.Append([]string{fmt.Sprintf("%+v", payment), "", "", "", "", ""})
		}
	}
	table.Render()
}

// Helper function to safely get string values from map
func getStringValue(m map[string]interface{}, key string) string {
	if val, exists := m[key]; exists {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

// Helper function to format number with comma separators
func formatWithCommas(num float64) string {
	// Convert to string without decimals first
	intPart := int64(num)
	str := strconv.FormatInt(intPart, 10)

	// Add commas from right to left
	var result []byte
	for i, j := 0, len(str); j > 0; i++ {
		if i > 0 && i%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, str[j-1])
		j--
	}

	// Reverse the result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

// RenderExpenses displays expenses in a table format
func RenderExpenses(expenses []responses.Expense) {
	if len(expenses) == 0 {
		fmt.Println("No expenses found.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Date", "Type", "Description", "Amount", "Status", "Created At"})

	for _, expense := range expenses {
		table.Append([]string{
			expense.ID,
			expense.Date,
			expense.Type,
			expense.Description,
			fmt.Sprintf("%.2f", expense.Amount),
			expense.Status,
			expense.CreatedAt,
		})
	}
	table.Render()
}

// FormatCurrency formats a float64 amount to a currency string
func FormatCurrency(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

// ParseDate parses a date string and returns it in a standardized format
func ParseDate(dateStr string) (string, error) {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return "", err
	}
	return t.Format("January 2, 2006 at 3:04pm"), nil
}

// ConvertToFloat converts a string to a float64, handling commas and other formatting
func ConvertToFloat(value string) (float64, error) {
	value = strings.ReplaceAll(value, ",", "")
	return strconv.ParseFloat(value, 64)
}
