# kafka_docker

A docker-compose wrapper that brings up Kafka brokers, ensures the advertised listeners are correct, and ensures
that any topics that the Kafka container is creating via KAFKA_CREATE_TOPICS are in place before returning

## Requirements

`kafka_docker` expects that you are using the [wurstmeister/kafka](https://hub.docker.com/r/wurstmeister/kafka/) images in a docker-compose.yml file, and that you are
setting the `KAFKA_CREATE_TOPICS` environment variable in the docker-compose.yml file.

It sets a `DOCKER_IP` environment variable which can be used in KAFKA_ADVERTISED_LISTENERS. For example, the following
configuration brings up a 3-broker cluster, creates topics `topic1` and `topic2` with 12 partitions and 3 replicas,
and allows access from the host system and also other containers:

```
version: "3.2"
services:
  zookeeper:
    image: wurstmeister/zookeeper:3.4.6
    ports:
      - "2181:2181"
  kafka1:
    image: wurstmeister/kafka:2.11-1.1.1
    ports:
      - "9092:9092"
    depends_on:
      - zookeeper
    environment:
      KAFKA_CREATE_TOPICS: "topic1:12:3,topic2:12:3"
      KAFKA_LISTENERS: "PLAINTEXT://0.0.0.0:9092"
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://${DOCKER_IP}:9092
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
  kafka2:
    image: wurstmeister/kafka:2.11-1.1.1
    ports:
      - "9093:9092"
    depends_on:
      - zookeeper
    environment:
      KAFKA_LISTENERS: "PLAINTEXT://0.0.0.0:9092"
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://${DOCKER_IP}:9093
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
  kafka3:
    image: wurstmeister/kafka:2.11-1.1.1
    ports:
      - "9094:9092"
    depends_on:
      - zookeeper
    environment:
      KAFKA_LISTENERS: "PLAINTEXT://0.0.0.0:9092"
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://${DOCKER_IP}:9094
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
```

## Usage
Bring kafka containers online
```
kafka_docker up [-f|--file <docker-compose-file>]
```
Halt kafka docker containers
```
kafka_docker down [-f|--file <docker-compose-file>]
```
Display ip associated to docker containers
```
kafka_docker ip 
```

## Testing

Tests are run via a Makefile. To download all dependencies, build, and run tests:
```
make
```

Additionally you can run linting with
```
make lint
```

It is strongly recommended that linting passes, but Travis-CI does not run the linters because golint specifies that
it may have false positives and shouldn't be relied on automatically.

Other make commands:
```
# Download dependencies:
make get
# Compile project:
make build
```
