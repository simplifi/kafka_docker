package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/simplifi/kafka_docker/internal/dockercompose"
	"github.com/spf13/cobra"
)

// Seconds to wait for kafka topic creation to finish
const maximumWaitTime = 30

func init() {
	rootCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "docker-compose up with extra options",
	Args:  cobra.MaximumNArgs(0),
	Run:   doUpCmd,
}

func doUpCmd(cmd *cobra.Command, _args []string) {
	dockerFilePath := rootCmd.PersistentFlags().Lookup("file").Value.String()
	fmt.Printf("Using %s as the docker-compose configuration\n", dockerFilePath)

	dockerComposeContent, err := ioutil.ReadFile(dockerFilePath)
	if err != nil {
		fmt.Printf("Error opening file `%s`: %v\n", dockerFilePath, err)
		os.Exit(1)
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

	fmt.Print("waiting for topics to be created")
	waitForTopics(&kafkaService, topics)
	fmt.Println("\ndone.")
}

func dockerComposeUp(dockerFilePath string) {
	stdout, stderr, err := bash("docker-compose", "-f", dockerFilePath, "up", "-d")
	if err != 0 {
		fmt.Println("Failure running docker-compose up")
		fmt.Printf("STDOUT:\n%s\n", stdout)
		fmt.Printf("STDERR:\n%s\n", stderr)
		os.Exit(err)
	}

}

// waitForTopics queries the running kafka instance to see if the topics are available
func waitForTopics(service *dockercompose.Service, expectedTopics []string) {
	// Wait for a bit for topics to be created (should take < 5 seconds)
	for i := 0; i < maximumWaitTime; i++ {
		if eql(createdTopics(service), expectedTopics) {
			return
		}
		fmt.Print(".")
		time.Sleep(time.Second)
	}
}

func createdTopics(service *dockercompose.Service) []string {
	// docker-compose exec -T -e topic=$topic kafka1 /bin/bash -c '$KAFKA_HOME/bin/kafka-topics.sh --zookeeper $KAFKA_ZOOKEEPER_CONNECT --list | grep -q $topic'
	stdout, stderr, err := bash("docker-compose", "exec", "-T", service.Name, "/bin/bash", "-c", "$KAFKA_HOME/bin/kafka-topics.sh --zookeeper $KAFKA_ZOOKEEPER_CONNECT --list")

	if err != 0 {
		fmt.Println("Failure running docker-compose exec")
		fmt.Printf("STDOUT:\n%s\n", stdout)
		fmt.Printf("STDERR:\n%s\n", stderr)
		os.Exit(err)
	}
	topics := strings.Split(strings.TrimSpace(stdout), "\n")
	return topics
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
