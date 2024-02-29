# KOM (Kubernetes Operator Manager)

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
```

2. Install KOM using Helm:

```sh
helm install kom kkb0318/kom
```

This command deploys KOM on the Kubernetes cluster in the default configuration. For more advanced configurations, refer to the Configuration section.

## Deploying Your First Operator with KOM

After installing KOM, you can deploy your first operator by following the example manifest provided in the examples/ directory.
This will help you understand how to specify the Helm chart URL and version for the operator you wish to deploy.

## Future Plans

- **Access to Private Repositories**: We are planning to enhance KOM's capabilities by enabling it to access and manage operators from private Git repositories.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
