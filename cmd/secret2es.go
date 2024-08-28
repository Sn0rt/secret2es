package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Sn0rt/sercert2extsecret/pkg/converter"
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
		Short: "A tool to convert Kubernetes secrets to External Secrets",
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
			err = converter.ConvertSecret(inputPath, storeType, storeName)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringP("input", "i", "", "Input path of corev1 secret file (required)")
	cmd.Flags().StringP("storetype", "s", "ClusterSecretStore", "Store type (optional)")
	cmd.Flags().StringP("storename", "n", "", "Store name (required)")

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
