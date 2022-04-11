# TCP Service -  Initial solution

The initial solution of the TCP service is a pretty simple GO application that will listen to a given port and process the supported commands.
After receiving a connection request, it will read the input and process the command.
If the input is not a valid command, it will return a proper message.

Processing the command consists in returning the appropriate response to the client.
Each connection request will be processed in a separate goroutine (similar to coroutines/threads), allowing multiple clients to connect simultaneously.

For simplicity of implementation, the service will read the input until it finds a new line (`\n`) and it will also send a response back with a new line as the signal for the end of the response.

## Requirements

To run this service, you will need Golang 17.x+. It may work on previous versions, but it was not tested.

Installation instructions are [here](https://go.dev/doc/install).

## Supported commands

- `WHO`: Outputs the total number of clients connected.
- `WHERE`: Outputs the id of the server (a unique identifier).
- `WHY`: Output the string `42`.

## Server (TCP Service)

Build binaries for the server.

```shell
# build server only
make build-server

# build both server and client
make build
```

Start a server in the default port (12345)

```shell
make start-server
```

Start a server in a specific port.

```shell
make start-server PORT=9999
```

## Client

The client was not part of the initial requirements but it proved itself useful for testing the solution.

Build binaries for the client.

```shell
# build client only
make build-client

# build both server and client
make build
```

Start a client that will run all commands in sequence.

```shell
make run-client

# without building the binaries
go run cmd/client/client.go

# building the binaries
./client
```

Options:

`-host=n`: Connect to host `n`. The default is `127.0.0.1`.
`-port=n`: Connect to port `n` on the given host. The default is `12345`.
`-clients=n`: Start `n` clients. The default is `1`.
`-iterations=n`: Repeat `n` times the sequece of commands. The default is `1`.
`-commands`: Comma separated list of commands to run in sequence. the default is `WHERE,WHO,WHY`.

```shell
make run-client OPTS="-host=127.0.0.1 -port=9999 -clients=4 -iterations=3 -commands=WHY,WHO,WHY,WHERE"

# without building the binaries
go run cmd/client/client.go -host=127.0.0.1 -port=9999 -clients=4 -iterations=3 -commands=WHY,WHO,WHY,WHERE

# building the binaries
./client -host=127.0.0.1 -port=9999 -clients=4 -iterations=3 -commands=WHY,WHO,WHY,WHERE
```

## Tests

To test the service, you can use the client provided in this repo according to the documentation on this page.

Unit tests were not written on purpose, due to time constraints. For a production service, unit tests and integration tests would be mandatory.

## Next steps

### CI/CD

#### Containers

Currently the service can run in any type of environment/machine, like bare metal, virtual machine or container, as it just listens to a port and process whatever input it receives.
For most (maybe any) types of applications/services, having a stable environment to run is pretty important, as installing new packages, changing operational system configuration and things like this can affect the behavior of the service.
Using containers will allow us to have a better control of the environment and use a lightweight solution for this.
The container with the service will allow to have much more controller way of deploying it to different environments (development, test, staging, production, etc), reducing the variables of each environment.
Also, running the application in a container allow the developers to run it locally in a much easier way and with an environment much closer to the ones where the application will actually run.

#### CI

It's important to have a CI process, so we can have a frequent and automated feedback about the changes being made.
Having a good Pull Request tool and automated builds for branches, pull requests and main branches is pretty important to receive this feedback with almost no effort.
The CI can run unit tests, integration tests, code coverage analysis, static analysis, security scanning (code and container), along with other steps that will provide us as much feedback as we can/want/can afford.

#### CD

Once the CI is well defined, it's also important/desired to have a good/well defined CD process. This will allow us to move faster to production with less risk. The CD will initially deploy the application to lower environments (developmnet, test, etc) so we can run batteries of tests and validate the changes. If no breaking issues found, the CD process can deploy it to upper environments, like staging and production.
It's important to have a rollback strategy, as any two given environments are not exactly the same, and we should be able to deploy the latest stable version of the application if the new one is not in a healthy state.
Also, depending on the needs, we can have no downtime deployments: simple rolling deployments, blue/green or canary deployments.
Each one has their pros and cons, so it will depend on the requirements.

### Maintainability

Maintainability (code hygiene, code robustness, documentation)

#### Developer Experience (DX)

The DX is also pretty important, as they will work on the service in a daily basis.

#### Documentation

The documentation is a pretty important part of an application, specially if it's a service/API that will be consumed by another team or third party. The current service is simple enough to have a simple readme page, but as we grow, we can add documentation generation tools based on code comments and annotations. This will allow us to have a good and standard way to have documentation without much effort.

#### Code

Codewise is important to organize the files in folders/packages/modules so it's easy to find what we are looking for and avoid having files with hundreds lines of code. Keeping the methods/classes/types short, simple and always having the single responsability concept in mind, will help us have a great source code.

#### Testing

Testing is another way of making the code neat and help with documentation.
Unit tests help us to work towards simple, short methods and provide eamples of their usage.
Also, it allow us to validate changes as part of the CI/CD and make sure we are not introducing broken code into production.

### Observability

#### Logs

The current solution is using the builtin `log` package for simplicity.
To improve the logs, we can use a package like `zap` ([here](https://github.com/uber-go/zap)), which will provide us a much more powerful way of logging.
Logging in JSON format is an important improvement, as we ingest them into an ELK stack based tool and then this tool can process the logs in a much more efficient way.

Also, we can write logs for specific scenarions that we want to monitor and then use these log entries to define SLIs and alerts.

ELK stack based tools allow us to set alerts based on number of log entries over time and also query the logs to ingest data for the SLIs.

#### SLO / SLI / SLA

We can define a few SLIs, like availability, latency and error rate for the service.
It will help us understand the health of the service and also set SLOs based on these metrics.
The SLOs will them allow us to agree on SLAs with the consumers of the service.

To define an SLO, we should monitor the SLI for a period of time (like 2-4 weeks) and then based on the historical data, define an appropriate SLO.
We can start setting it to a given/desired number, like 90%, and adjust accordingly.

The current solution is simple enough to say that should be able to have a SLO of 90-95%, as there aren't many parts that can fail other then the connection between the client and the server.

#### Dashboards and alerts

Using dashboard and alerts tool like Grafana we can monitor the application by setting thresholds for some metrics and alert if the thresholds are broken.
Alerting on lower thresholds than those used for SLO, will allow us to react before the issues become a bigger problem.
