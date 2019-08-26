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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/simplifi/kafka_docker/internal/dockercompose"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upCmd)
}

func doUpCmd(cmd *cobra.Command, _args []string) {
	dockerFilePath := rootCmd.PersistentFlags().Lookup("file").Value.String()
	fmt.Printf("Using %s as the docker-compose configuration\n", dockerFilePath)

	dockerComposeContent, err := ioutil.ReadFile(dockerFilePath)
	if err != nil {
		fmt.Printf("Error opening file `%s`: %v\n", dockerFilePath, err)
		os.Exit(128)
	}

	dc, err := dockercompose.Parse(dockerComposeContent)
	if err != nil {
		fmt.Printf("Error parsing docker-compose configuration\n")
		panic(err)
	}

	kafkaService, err := dc.FindKafkaContainer()
	topics := kafkaService.GetTopics()
	fmt.Printf("docker-compose service found: %s, with topics %v\n", kafkaService.Name, topics)

	dockerIP, err := dockercompose.DockerIP()
	if err != nil {
		panic(err)
	}
	fmt.Printf("export DOCKER_IP=%s\n", dockerIP)
	os.Setenv("DOCKER_IP", dockerIP)

	fmt.Printf("docker-compose -f %s up -d\n", dockerFilePath)
	dockerComposeUp(dockerFilePath)

	fmt.Println("waiting for topics to be created...")
	waitForTopics(&kafkaService, topics)
	fmt.Println("done.")
}

func dockerComposeUp(dockerFilePath string) {
	cmd := exec.Command("docker-compose", "-f", dockerFilePath, "up", "-d")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failure running docker-compose up: %v", err)
		fmt.Printf("STDOUT:\n%s\n", stdout.String())
		fmt.Printf("STDERR:\n%s\n", stderr.String())
		os.Exit(1)
	}

}

// waitForTopics queries the running kafka instance to see if the topics are available
func waitForTopics(service *dockercompose.Service, expectedTopics []string) {
	// Wait for one minute for topics to be created (should take < 5 seconds)
	for i := 0; i < 60; i++ {
		if eql(createdTopics(service), expectedTopics) {
			return
		}
		time.Sleep(time.Second)
	}
}

// Returns true if the sorted content of the slices are equal. Sorts the passed in slices in place
func eql(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func createdTopics(service *dockercompose.Service) []string {
	// docker-compose exec -T -e topic=$topic kafka1 /bin/bash -c '$KAFKA_HOME/bin/kafka-topics.sh --zookeeper $KAFKA_ZOOKEEPER_CONNECT --list | grep -q $topic'
	cmd := exec.Command("docker-compose", "exec", "-T", service.Name, "/bin/bash", "-c", "$KAFKA_HOME/bin/kafka-topics.sh --zookeeper $KAFKA_ZOOKEEPER_CONNECT --list")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failure running docker-compose exec: %v", err)
		fmt.Printf("STDOUT:\n%s\n", stdout.String())
		fmt.Printf("STDERR:\n%s\n", stderr.String())
		os.Exit(1)
	}
	return strings.Split(stdout.String(), "\n")
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "docker-compose up with extra options",
	Args:  cobra.MaximumNArgs(0),
	Run:   doUpCmd,
}
