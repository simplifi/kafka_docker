/*
Copyright © 2019 Simpli.fi Holdings

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
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
