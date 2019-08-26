# kafka_docker

A docker-compose wrapper that brings up Kafka brokers, ensures the advertised listeners are correct, and ensures
that any topics that the Kafka container is creating via KAFKA_CREATE_TOPICS are in place before returning

## Usage

```
kafka_docker up [-f|--file <docker-compose-file>]
```

## Testing

Tests are run via a Makefile. To download all dependencies, build, and run tests:
```
make
```

Other make commands:
```
# Download dependencies:
make get
# Compile project:
make build
```
