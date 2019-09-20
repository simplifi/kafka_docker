package dockercompose

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// Service represents a service defined in a docker-compose file
type Service struct {
	Name        string
	Image       string
	Environment interface{}
}

// Casts the Environment (which could be a []string or map[string]string, but has to be parsed as an
// interface{}, which becomes a map[interface{}]interface{} in the YAML library to a map[string]string.
func (s *Service) getEnvironment() map[string]string {
	switch val := s.Environment.(type) {
	case map[string]string:
		return val
	case map[interface{}]interface{}:
		r := make(map[string]string)
		for k, v := range val {
			r[k.(string)] = v.(string)
		}
		return r
	default:
		return make(map[string]string)
	}
}

// DockerCompose represents a docker-compose file
type DockerCompose struct {
	Version  string
	Services map[string]Service
}

// Parse converts a yaml byte slice into a DockerCompose struct
func Parse(yamlFile []byte) (DockerCompose, error) {
	var dc DockerCompose
	err := yaml.Unmarshal(yamlFile, &dc)
	return dc, err

}

// FindKafkaContainer looks for a service that has the KAFKA_CREATE_TOPICS environment variable set, and returns it
// as a Service struct.
func (dc *DockerCompose) FindKafkaContainer() (Service, error) {
	for name, service := range dc.Services {
		if strings.HasPrefix(service.Image, "wurstmeister/kafka") && service.getEnvironment()["KAFKA_CREATE_TOPICS"] != "" {
			service.Name = name
			return service, nil
		}
	}
	return Service{}, errors.New("No Kafka container found with KAFKA_CREATE_TOPICS set")
}

// GetTopics returns a list of the kafka topics that a service creates
func (service *Service) GetTopics() []string {
	rawTopicString := service.getEnvironment()["KAFKA_CREATE_TOPICS"]
	if rawTopicString == "" {
		return make([]string, 0)
	}
	topicDefinitions := strings.Split(rawTopicString, ",")
	topics := make([]string, len(topicDefinitions))
	for i := 0; i < len(topicDefinitions); i++ {
		topicDefinition := strings.Split(topicDefinitions[i], ":")
		topics[i] = topicDefinition[0]
	}
	return topics
}

// DockerIP finds the first private IP found in the network interfaces
func DockerIP() (string, error) {
	var ip net.IP
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Could not determine interfaces")
		panic(err)
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		// If there's an error getting the addresses for an interface, then ignore it
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if isPrivate(ip) {
				return ip.String(), nil
			}
		}
	}

	return "", errors.New("Unable to determine private IP address")
}

func isPrivate(ip net.IP) bool {
	// 10.x.x.x
	r := "^10\\.|"
	// 172.[16-31].x.x
	r += "^172\\.1[6-9]\\.|"
	r += "^172\\.2[0-9]\\.|"
	r += "^172\\.3[0-1]\\.|"
	// 192.168.x.x
	r += "^192\\.168\\."
	var privateIPRegex, _ = regexp.Compile(r)

	return privateIPRegex.MatchString(ip.String())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
