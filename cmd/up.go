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

var timeout uint32

func init() {
	rootCmd.AddCommand(upCmd)
	upCmd.PersistentFlags().Uint32VarP(&timeout, "timeout", "t", 30, "seconds to wait on the topics to become ready")
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

	fmt.Println("waiting for topics to be created")
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
	for i := uint32(0); i < timeout; i++ {
		foundTopics := createdTopics(service)
		if eql(foundTopics, expectedTopics) {
			return
		}
		fmt.Print(".")
		time.Sleep(time.Second)
	}
	if eql(createdTopics(service), expectedTopics) {
		return
	}
	msg := fmt.Sprintf("Kafka Topics did not become available after %d seconds", timeout)
	panic(msg)
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
	rawTopics := strings.Split(strings.TrimSpace(stdout), "\n")
	// Internal topic
	topics := removeElem(rawTopics, "__consumer_offsets")
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

func removeElem(s []string, e string) []string {
	k := -1
	for i := 0; i < len(s); i++ {
		if s[i] == e {
			k = i
			break
		}
	}
	if k >= 0 {
		s[k] = s[len(s)-1]
		// We do not need to put s[i] at the end, as it will be discarded anyway
		return s[:len(s)-1]
	}
	return s
}
