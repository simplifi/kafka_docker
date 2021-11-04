# kafka_docker

A docker-compose wrapper that brings up Kafka brokers, ensures the advertised listeners are correct, and ensures
that any topics that the Kafka container is creating via KAFKA_CREATE_TOPICS are in place before returning

## Installation

### Homebrew (preferred)

- Setup the Simplifi Homebrew Tap by following these directions: https://github.com/simplifi/homebrew-tap#how-do-i-set-this-up
- Run `brew install kafka_docker`

### Manual

- Browse to the [Releases](https://github.com/simplifi/kafka_docker/releases)
- Navigate to the latest release
- Click to download the appropriate release for your platform
- Once downloaded, extract, then move the binary into your path and make it executable
  - ex: `mv ~/Downloads/kafka_docker /usr/local/bin/kafka_docker && chmod +x /usr/local/bin/kafka_docker`

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

You can see the options available to you by running kafka_docker without any arguments:

```bash
Start up docker-compose with kafka.
	Scans the docker-compose.yml file and finds a kafka container, and ensures that the advertised connection
	is set correctly to allow the host to connect, but still allow inter-container communication.

	Usage:
	Starting up docker-compose:
	kafka_docker up [-f /path/to/docker-compose.yml]

	For symmetry also there is kafka_docker down, which just calls docker-compose down.

	Defaults to looking in $PWD for the docker-compose.yml

Usage:
  kafka_docker [command]

Available Commands:
  down        Runs docker-compose down
  help        Help about any command
  ip          Outputs Docker IP
  up          docker-compose up with extra options

Flags:
  -f, --file string   docker-compose file
  -h, --help          help for kafka_docker

Use "kafka_docker [command] --help" for more information about a command.
```

## Example Usages
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

## Development

### Tool Setup

[Install Golang](https://golang.org/doc/install)

[Install golint](https://github.com/golang/lint#installation)

[Install goreleaser](https://goreleaser.com/install/)

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

### Releasing

#### Automated

This project is using [goreleaser](https://goreleaser.com). GitHub release
creation is automated using Travis CI. New releases are automatically created
when new tags are pushed to the repo.

```shell script
$ TAG=v0.1.0 make tag
```

#### Manual

A release can also be manually pushed using goreleaser. You must have the
`GITHUB_TOKEN` environment variable set to a GitHub token with the `repo` scope.
You can create a new github token [here](https://github.com/settings/tokens/new).

```shell script
$ make release
```
