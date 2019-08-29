package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(downCmd)
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Runs docker-compose down",
	Args:  cobra.MaximumNArgs(0),
	Run:   doDownCmd,
}

func doDownCmd(cmd *cobra.Command, _args []string) {
	dockerFilePath := rootCmd.PersistentFlags().Lookup("file").Value.String()
	fmt.Printf("Using %s as the docker-compose configuration\n", dockerFilePath)

	fmt.Printf("docker-compose -f %s down\n", dockerFilePath)
	dockerComposeDown(dockerFilePath)
}

func dockerComposeDown(dockerFilePath string) {
	stdout, stderr, err := bash("docker-compose", "-f", dockerFilePath, "down")
	if err != 0 {
		fmt.Println("Failure running docker-compose down")
		fmt.Printf("STDOUT:\n%s\n", stdout)
		fmt.Printf("STDERR:\n%s\n", stderr)
		os.Exit(1)
	}
}
