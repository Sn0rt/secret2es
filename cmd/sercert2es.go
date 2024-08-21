package main

import (
	"fmt"
	"os"

	"github.com/Sn0rt/sercert2extsecret/pkg/converter"
	"github.com/spf13/cobra"
)

var (
	version   string
	buildTime string
)

func main() {
	if err := setupAndExecute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func setupAndExecute() error {
	rootCmd := &cobra.Command{
		Use:   "sercert2es",
		Short: "A tool to convert Kubernetes secrets to External Secrets",
	}

	rootCmd.AddCommand(extSecretGenCmd())
	rootCmd.AddCommand(versionCmd())

	return rootCmd.Execute()
}

func extSecretGenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extsecret-gen",
		Short: "Generate external secrets from corev1 secrets",
		RunE: func(cmd *cobra.Command, args []string) error {
			input, _ := cmd.Flags().GetString("input")
			_, _ = cmd.Flags().GetString("output")
			store, _ := cmd.Flags().GetString("store")
			storeName, _ := cmd.Flags().GetString("storename")
			namespace, _ := cmd.Flags().GetString("namespace")
			secretName, _ := cmd.Flags().GetString("secret-name")
			verbose, _ := cmd.Flags().GetBool("verbose")

			err := converter.ConvertSecret(input, store, storeName, namespace, secretName, verbose)
			if err != nil {
				return err
			}

			if verbose {
				fmt.Println("Conversion completed successfully")
			}
			return nil
		},
	}

	cmd.Flags().StringP("input", "i", "", "Input path of corev1 secret file (required)")
	cmd.Flags().StringP("output", "o", "", "Output path external secret file (required)")
	cmd.Flags().StringP("store", "s", "ClusterSecretStore", "Store type (optional)")
	cmd.Flags().StringP("storename", "n", "", "Store name (required)")
	cmd.Flags().String("namespace", "default", "External namespace (optional)")
	cmd.Flags().String("secret-name", "", "Secret name (optional)")
	cmd.Flags().Bool("verbose", false, "Enable verbose output (optional)")

	cmd.MarkFlagRequired("input")
	cmd.MarkFlagRequired("output")
	cmd.MarkFlagRequired("storename")

	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of sercert2extsecret",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("sercert2extsecret version %s\n", version)
			fmt.Printf("Built at %s\n", buildTime)
		},
	}
}
