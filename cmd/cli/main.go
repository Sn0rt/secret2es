package main

import (
	"fmt"
	"os"

	esv1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"

	"github.com/spf13/cobra"

	"github.com/Sn0rt/secret2es/pkg/converter"
)

var (
	version   string
	buildTime string
)

func main() {
	if err := setupAndExecute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func setupAndExecute() error {
	rootCmd := &cobra.Command{
		Use:   "secret2es",
		Short: "A tool to convert AVP secrets to ExternalSecrets",
	}

	rootCmd.AddCommand(extSecretGenCmd())
	rootCmd.AddCommand(versionCmd())

	return rootCmd.Execute()
}

func extSecretGenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "es-gen",
		Short: "Generate external secrets from corev1 secrets",
		RunE: func(cmd *cobra.Command, args []string) error {
			inputPath, err := cmd.Flags().GetString("input")
			if err != nil {
				return err
			}
			storeType, err := cmd.Flags().GetString("storetype")
			if err != nil {
				return err
			}
			storeName, err := cmd.Flags().GetString("storename")
			if err != nil {
				return err
			}
			if storeName == "" {
				return fmt.Errorf("store name is required")
			}
			creationPolicy, err := cmd.Flags().GetString("creation-policy")
			if err != nil {
				return err
			}
			if creationPolicy == "" {
				return fmt.Errorf("creation policy is required")
			}
			resolve, err := cmd.Flags().GetBool("resolve")
			if err != nil {
				return err
			}
			refreshPolicy, err := cmd.Flags().GetString("refresh-policy")
			if err != nil {
				return err
			}
			refreshInterval, err := cmd.Flags().GetString("refresh-interval")
			if err != nil {
				return err
			}

			// Validate refresh policy and interval combination
			if refreshPolicy == "Periodic" && refreshInterval == "" {
				return fmt.Errorf("refresh-interval is required when refresh-policy is Periodic")
			}
			if refreshPolicy != "Periodic" && refreshInterval == "" {
				refreshInterval = "0s" // Set default for non-Periodic policies
			}

			err = converter.ConvertSecret(inputPath, storeType, storeName, esv1.ExternalSecretCreationPolicy(creationPolicy), resolve, esv1.ExternalSecretRefreshPolicy(refreshPolicy), refreshInterval)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringP("input", "i", "", "Input path of corev1 secret file (required)")
	cmd.Flags().StringP("storetype", "s", "SecretStore", "Store type (optional)")
	cmd.Flags().StringP("storename", "n", "", "Store name (required)")
	cmd.Flags().StringP("creation-policy", "c", "Owner", "Create policy, only Owner, Orphan")
	cmd.Flags().BoolP("resolve", "r", false, "Resolve the <% ENV %> from env")
	cmd.Flags().StringP("refresh-policy", "p", "OnChange", "Refresh policy: CreatedOnce, Periodic, or OnChange")
	cmd.Flags().StringP("refresh-interval", "t", "", "Refresh interval (e.g., 30s, 5m, 1h) - only required when refresh-policy is Periodic")

	err := cmd.MarkFlagRequired("input")
	if err != nil {
		return nil
	}

	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of secret2es",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("secret2es version %s\n", version)
			fmt.Printf("Built at %s\n", buildTime)
		},
	}
}
