# knative-eventing-demo
This repository contains a simple Knative Eventing demo using event-driven steps orchestrated in workflows. It includes Go-based event processing services and Makefiles to automate building, deploying, testing, and managing the resources.

## Quick Start

```
# Build and deploy
make all
# Invoke workflow
make send
# Check logs
make logs
# Clean up
make clean
```

## Makefile Commands
The Makefile provides automation for building, deploying, testing, and managing Knative eventing demo resources. Below is a concise description of each available command:

- **make all**
  Runs the full workflow: login, build, push, deploy, and show resource status.

- **make login**
  Logs in to the container registry using environment variables `HARBOR_USER` and `HARBOR_PASS` if set.

- **make docker-build**
  Builds the Docker image for the service.

- **make docker-push**
  Pushes the built Docker image to the registry.

- **make deploy**
  Deploys Knative services and sequence resources to the specified namespace (default: `default`).

- **make undeploy**
  Deletes the deployed Knative services and sequence resources from the namespace.

- **make show**
  Displays the status of Knative services and the sequence address in the namespace.

## Testing and Debugging

- **make send**
  Sends a CloudEvent to the sequence endpoint from within the cluster for integration testing.

- **make test-step [STEP=n]**
  Sends a test CloudEvent to a specific Knative service step (default: step 1). Override with `STEP` variable.

- **make logs**
  Tails logs for all step services' user-containers in the namespace.

- **make clean**
  Cleans up all resources, including services, sequence, channels, subscriptions, and temporary pods.

---

For namespace customization, use `NAMESPACE=my-namespace` as a variable with any command.
