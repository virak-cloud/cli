package cmd

import (
	"bufio"
	"fmt"
	urls "github.com/virak-cloud/cli/pkg"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"
	"os"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"github.com/pkg/browser"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:     "login",
	Aliases: []string{"log-in", "auth"},
	Short:   "Login to the Virak Cloud API",
	Long:    `Login command allows you to authenticate with the Virak Cloud API.`,
	Run: func(cmd *cobra.Command, args []string) {

		token, _ := cmd.Flags().GetString("token")
		if token == "" {
			var machineID string
			if id, errMachineID := machineid.ID(); errMachineID != nil {
				slog.Error("failed to get machine-id", "error", errMachineID)
				fmt.Println("Failed to get machine-id:", errMachineID)
				machineID = "unknown-machine-id"
			} else {
				machineID = id
			}
			// create a machine-id of this machine
			url := fmt.Sprintf(urls.LoginUrl, machineID)
			err := browser.OpenURL(string(url))
			if err != nil {
				slog.Error("failed to open browser", "error", err)
				fmt.Println("Failed to open browser:", err)
			}

			fmt.Println("Please log in to the Virak Cloud API in your browser.")
			fmt.Println("After logging in, you will receive a token. Use this token to authenticate future requests.")
			fmt.Println("You can also use the --token flag to provide your token directly.")
			fmt.Println("For more information, visit the Virak Cloud documentation. \n\n  ")
			fmt.Println("Enter token: ")
			reader := bufio.NewReader(os.Stdin)
			inputToken, err := reader.ReadString('\n')
			if err != nil {
				slog.Error("failed to get token from user", "error", err)
				fmt.Println("Failed to get token from user")
				return
			}
			token = inputToken
			token = strings.TrimSpace(token)
		}

		client := http.NewClient(token)
		_, err := client.GetTokenAbilities()

		if err != nil {
			slog.Error("failed to login with the provided token", "error", err)
			fmt.Println("Failed to login with the provided token. Please check your token and try again.")
			return
		}

		viper.Set("auth.token", token)
		err = viper.SafeWriteConfig()
		if err != nil {
			err = viper.WriteConfig()
		}
		if err != nil {
			slog.Error("failed to save token to config", "error", err)
			fmt.Println("Failed to save token to config:", err)
			return
		}
		fmt.Println("Login successful. Token saved to config.")
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)

	loginCmd.PersistentFlags().String("token", "", "if you have a token, you can use it to login without opening the browser")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
