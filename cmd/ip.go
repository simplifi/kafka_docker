package cmd

import (
	"fmt"

	"github.com/simplifi/kafka_docker/internal/dockercompose"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ipCmd)
}

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Outputs Docker IP",
	Args:  cobra.MaximumNArgs(0),
	Run:   doipCmd,
}

func doipCmd(cmd *cobra.Command, _args []string) {
	dockerComposeIP()
}

func dockerComposeIP() {
	dockerIP, err := dockercompose.DockerIP()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", dockerIP)
}
