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
	"testing"
)

var dockerComposeYaml = []byte(`
version: "3.2"
services:
  zookeeper:
    image: wurstmeister/zookeeper:3.4.6
    ports:
      - "2181:2181"
  the_right_kafka:
    image: wurstmeister/kafka:2.11-1.1.1
    ports:
      - "9092:9092"
    depends_on:
      - zookeeper
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      KAFKA_CREATE_TOPICS: "topic1:1:1"
      KAFKA_LISTENERS: "PLAINTEXT://0.0.0.0:9092"
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: "PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT"
      KAFKA_ADVERTISED_LISTENERS: "PLAINTEXT://${DOCKER_IP}:9092"
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
  the_wrong_kafka:
    image: wurstmeister/kafka:2.11-1.1.1
    ports:
      - "9093:9092"
    depends_on:
      - zookeeper
    environment:
      KAFKA_LISTENERS: "PLAINTEXT://0.0.0.0:9092"
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://${DOCKER_IP}:9093
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
  the_not_right_kafka:
    image: wurstmeister/kafka:2.11-1.1.1
    ports:
      - "9094:9092"
    depends_on:
      - zookeeper
    environment:
      KAFKA_LISTENERS: "PLAINTEXT://0.0.0.0:9092"
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://${DOCKER_IP}:9094
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
`)

func TestParse(t *testing.T) {
	dockerCompose, err := Parse(dockerComposeYaml)
	assertEqual(t, nil, err, "Does not throw an error")
	assertEqual(t, "3.2", dockerCompose.Version, "Parses the version correctly")
}

func TestFindKafkaContainer(t *testing.T) {
	dockerCompose, err := Parse(dockerComposeYaml)
	kafkaContainer, err := dockerCompose.FindKafkaContainer()
	assertEqual(t, nil, err, "Does not throw an error")
	assertEqual(t, "the_right_kafka", kafkaContainer.Name, "finds the correct container")

}

func TestGetTopics(t *testing.T) {
	service := Service{
		Name:  "kafka",
		Image: "wurstmeister/kafka:2.11-1.1.1",
		Environment: map[string]string{
			"KAFKA_CREATE_TOPICS": "topic1:12:3",
		},
	}
	topics := service.GetTopics()
	assertEqual(t, 1, len(topics), "Only one topic is found")
	assertEqual(t, topics[0], "topic1", "Finds the topics in KAFKA_CREATE_TOPICS")
}

func assertEqual(t *testing.T, expected interface{}, got interface{}, message string) {
	if expected != got {
		t.Fatalf("%s %s: Expected %v, got %v", t.Name(), message, expected, got)
	}
}
