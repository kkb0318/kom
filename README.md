# KOM (Kubernetes Operator Manager)

[![GitHub release](https://img.shields.io/github/release/kkb0318/kom.svg?maxAge=60)](https://github.com/kkb0318/kom/releases)
[![CI](https://github.com/kkb0318/kom/actions/workflows/ci.yaml/badge.svg)](https://github.com/kkb0318/kom/actions/workflows/ci.yaml)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/kom)](https://artifacthub.io/packages/search?repo=kom)

## Overview

KOM, which stands for Kubernetes Operator Manager, is an open-source software tool designed to streamline the management of Kubernetes operators. It acts as an operator itself, facilitating the deployment, management, and removal of Kubernetes operators with minimal hassle.

## Features

- **Git Repository Support**: Enables integration with Git repositories to manage operator configurations and their versions effectively.
- **Chart Release Mechanism**: Incorporates a system for deploying Helm charts, allowing for easy installation and management of Kubernetes applications.

## Getting Started

### Prerequisites

- Kubernetes cluster with Flux 2.x or Argo CD installed: KOM requires a Kubernetes cluster that is already equipped with either Flux version 2.x or Argo CD for GitOps-based management.

### Installation

To install KOM on your Kubernetes cluster, follow these steps:

1. Add the KOM Helm repository:

```sh
helm repo add kkb0318 https://kkb0318.github.io/kom
helm repo update
```

2. Install KOM using Helm:

```sh
helm install kom kkb0318/kom -n kom-system --create-namespace
```

This command deploys KOM on the Kubernetes cluster in the default configuration. For more advanced configurations, refer to the Configuration section.

## Deploying with KOM

After installing KOM, you can deploy the operator using `OperatorManager` manifest.

```yaml
apiVersion: kom.kkb0318.github.io/v1alpha1
kind: OperatorManager
metadata:
  name: kom
  namespace: kom-system
spec:
  tool: flux
  cleanup: true
  resource:
    helm:
      - name: jetstack
        url: https://charts.jetstack.io
        charts:
          - name: cert-manager
            version: v1.14.4
            values:
              installCRDs: true
              prometheus:
                enabled: false
      - name: repo1
        url: https://helm.github.io/examples
        charts:
          - name: hello-world
            version: x.x.x
```

You can find more details about the example manifests in the `examples/` directory.

## API Reference

You can find the reference in the [Reference](./docs/api.md) file.

## Future Plans

- **Access to Private Repositories**: We are planning to enhance KOM's capabilities by enabling it to access and manage operators from private Git repositories.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
