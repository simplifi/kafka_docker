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
	Environment map[string]string
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
		if strings.HasPrefix(service.Image, "wurstmeister/kafka") && service.Environment != nil && service.Environment["KAFKA_CREATE_TOPICS"] != "" {
			service.Name = name
			return service, nil
		}
	}
	return Service{}, errors.New("No Kafka container found with KAFKA_CREATE_TOPICS set")
}

// GetTopics returns a list of the kafka topics that a service creates
func (service *Service) GetTopics() []string {
	rawTopicString := service.Environment["KAFKA_CREATE_TOPICS"]
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
	var privateIPRegex, _ = regexp.Compile(
		// 10.x.x.x
		"^10\\.|" +

			// 172.[16-31].x.x
			"^172\\.1[6-9]\\.|" +
			"^172\\.2[0-9]\\.|" +
			"^172\\.3[0-1]\\.|" +

			// 192.168.x.x
			"^192\\.168\\.")

	return privateIPRegex.MatchString(ip.String())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
