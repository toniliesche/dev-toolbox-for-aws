# Toni's Developer Toolbox for Local AWS Testing

Goal of this project is to provide a toolbox that mimics basic behaviour of AWS services like SQS, SNS, Secrets Manager 
and more. This toolbox is intended to be used for local development and testing purposes only.

There is no guarantee that the behaviour of the tools in this toolbox is 100% identical to the behaviour of the AWS,
especially the parts of authentication and authorization have been left out on purpose for simplicity reasons.

## Local SNS to SQS Publisher

### Description

This application is a simple tool that can receive messages on an SNS topic and forward them to an SQS compatible queue
locally.

### Build

To build the binary locally, you need a running version of `Golang` version `1.24` or above as well as `make`. To start
the build process, simply run the following command inside your local shell:

```sh
make build-sns-publisher
```

### Configure

If you want to run this tool, you have two ways to provide all necessary configuration data. If you want to run the tool
by using a configuration file, you have to set an environment variable to the correct path of the configuration file:

```sh
export SNS_PUBLISHER_CONFIG_FILE=/path/to/config.yaml
```

Otherwise the tool will fallback to the default, environment-based solution.

#### Via Configuration File

Passing the configuration via a file is the most convenient way to provide all necessary data. Here you can find an easy
example config that contains all valid data:

```yaml
sns:
  topics:
    example-topic:
      subscribers:
        - type: sqs
        - identifier: example-queue
sqs:
  region: eu-central-1
  endpoint: http://sqs.fqdn:9324
```

#### Via Environment Variables

Passing the configuration via environment variables is the second way to provide all necessary data. This is especially
useful if you want to run your tools inside docker. Here you can find an easy example that contains all valid data:

```sh
export AWS_REGION=eu-central-1
export SQS_ENDPOINT=http://sqs.fqdn:9324
export SNS_TOPICS=example-topic,example-topic-2
export SNS_SQS_SUBSCRIPTIONS=example-topic:example-queue,example-topic-2:example-queue
```

### Run Locally

To run the tool locally, you can simply use the following command:

```sh
build/sns-publisher
```

### Run in Docker

To run the tool in docker, you can easily use a docker-compose.yml file that will create and run the container:

```yaml
services:
  sns-publisher:
    image: tliesche/sns-publisher:latest
    hostname: sns-publisher
    environment:
      - AWS_REGION=eu-central-1
      - SQS_ENDPOINT=http://sqs.fqdn:9324
      - SNS_TOPICS=example-topic,example-topic-2
      - SNS_SQS_SUBSCRIPTIONS=example-topic:example-queue,example-topic-2:example-queue
```

To run the container with docker compose v2, you can use the following command:

```sh
docker compose -f /path/to/docker-compose.yml -p sns-publisher up -d
```

## Local SQS Lambda Trigger

### Description

This application is a simple tool that can retrieve message events from an SQS compatible queue
(e.g. ElasticMQ) and trigger an AWS Lambda function locally.

### Build

To build the binary locally, you need a running version of `Golang` version `1.24` or above as well as `make`. To start
the build process, simply run the following command inside your local shell:

```sh
make build-sqs-lambda-trigger
```

### Configure

If you want to run this tool, you have two ways to provide all necessary configuration data. If you want to run the tool
by using a configuration file, you have to set an environment variable to the correct path of the configuration file:

```sh
export SQS_LAMBDA_TRIGGER_CONFIG_FILE=/path/to/config.yaml
```

Otherwise the tool will fallback to the default, environment-based solution.

#### Via Configuration File

Passing the configuration via a file is the most convenient way to provide all necessary data. Here you can find an easy
example config that contains all valid data:

```yaml
lambda:
  endpoint: http://lambda.fqdn:8080
  concurrency: 1
  invocation_type: RequestResponse # RequestResponse | Event
sqs:
  region: eu-central-1
  endpoint: http://sqs.fqdn:9324
  queue_name: example-queue
  message_batch_size: 10
  wait_time: 10
payload_format: sqs-event # sqs-event | pure
```

#### Via Environment Variables

Passing the configuration via environment variables is the second way to provide all necessary data. This is especially
useful if you want to run your tools inside docker. Here you can find an easy example that contains all valid data:

```sh
export AWS_REGION=eu-central-1
export LAMBDA_ENDPOINT=http://lambda.fqdn:8080
export LAMBDA_CONCURRENCY=1
export LAMBDA_INVOCATION_TYPE=RequestResponse
export SQS_ENDPOINT=http://sqs.fqdn:9324
export SQS_QUEUE_NAME=example-queue
export SQS_MESSAGE_BATCH_SIZE=10
export SQS_WAIT_TIME=10
export PAYLOAD_FORMAT=sqs-event # sqs-event | pure
```

### Run Locally

To run the tool locally, you can simply use the following command:

```sh
build/sqs-lambda-trigger
```

### Run in Docker

To run the tool in docker, you can easily use a docker-compose.yml file that will create and run the container:

```yaml
services:
  sqs-lambda-trigger:
    image: tliesche/sqs-lambda-trigger:latest
    hostname: sqs-lambda-trigger
    environment:
      - AWS_REGION=eu-central-1
      - LAMBDA_ENDPOINT=http://lamba.fqdn:8080
      - PAYLOAD_FORMAT=sqs-event
      - SQS_ENDPOINT=http://sqs.fqdn:9324
      - SQS_QUEUE_NAME=example-queue
      - SQS_MESSAGE_BATCH_SIZE=10
      - SQS_WAIT_TIME=20
```

To run the container with docker compose v2, you can use the following command:

```sh
docker compose -f /path/to/docker-compose.yml -p sqs-lambda-trigger up -d
```

## Local Secrets Manager (get-only)

### Description

This application is a simple tool that mimics the AWS Secrets Manager API. It can be used to retrieve secrets from a 
local in-memory storage.

### Build

To build the binary locally, you need a running version of `Golang` version `1.24` or above as well as `make`. To start
the build process, simply run the following command inside your local shell:

```sh
make build-secrets-manager
```

### Configure

If you want to run this tool, you have two ways to provide all necessary configuration data. If you want to run the tool
by using a configuration file, you have to set an environment variable to the correct path of the configuration file:

```sh
export SECRETS_MANAGER_CONFIG_FILE=/path/to/config.yaml
```

Otherwise the tool will fallback to the default, environment-based solution.

#### Via Configuration File

Passing the configuration via a file is the most convenient way to provide all necessary data. Here you can find an easy
example config that contains all valid data:

```yaml
account: "000000000000"
region: "eu-central-1"
secrets:
  - name: example-secret
    versions:
      - version1_value
      - version2_value
```

Note: The secrets values are stored in a list from NEWEST to OLDEST.

#### Via Environment Variables

Passing the configuration via environment variables is the second way to provide all necessary data. This is especially
useful if you want to run your tools inside docker. Here you can find an easy example that contains all valid data:

```sh
export AWS_REGION=eu-central-1
export AWS_ACCOUNT_ID=000000000000
export "SECRETS=example-secret=value1,value2;example-secret-2=value1,value2"
```

### Run Locally

To run the tool locally, you can simply use the following command:

```sh
build/secrets-manager
```

### Run in Docker

To run the tool in docker, you can easily use a docker-compose.yml file that will create and run the container:

```yaml
services:
  secrets-manager:
    image: tliesche/secrets-manager:latest
    hostname: secrets-manager
    environment:
      - AWS_REGION=eu-central-1
      - AWS_ACCOUNT_ID=000000000000
      - SECRETS_MANAGER_SECRETS=example-secret=value1,value2;example-secret-2=value1,value2
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This project is not affiliated with, endorsed by, or sponsored by Amazon Web Services (AWS). It is an independent 
project that provides mock implementations for AWS-like services for local development and testing. The names of AWS
services (SQS, SNS, Secrets Manager) are used solely for reference purposes.
