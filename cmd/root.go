package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kafka_docker",
	Short: "Start up docker-compose with kafka",
	Long: `Start up docker-compose with kafka.
	Scans the docker-compose.yml file and finds a kafka container, and ensures that the advertised connection
	is set correctly to allow the host to connect, but still allow inter-container communication.

	Usage:
	Starting up docker-compose:
	kafka_docker up [-f /path/to/docker-compose.yml]

	For symmetry also there is kafka_docker down, which just calls docker-compose down.

	Defaults to looking in $PWD for the docker-compose.yml
	`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "file", "f", defaultDockerCompose(), "docker-compose (default is $PWD/docker-compose.yml)")
}

func defaultDockerCompose() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Couldn't find current directory")
		os.Exit(1)
	}
	return dir + "/docker-compose.yml"
}

// Runs a bash command, returning stdout, stderr, and exit code if any.
func bash(command string, args ...string) (string, string, int) {
	cmd := exec.Command(command, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return stdout.String(), stderr.String(), exitError.ExitCode()
		}
		// Unknown error code, return 255
		return stdout.String(), stderr.String(), 255
	}
	return stdout.String(), stderr.String(), 0
}
