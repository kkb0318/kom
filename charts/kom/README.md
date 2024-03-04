# kom

[kom](https://github.com/kkb0318/kom) Kubernetes Operator Manager is an open-source software tool designed to streamline the management of Kubernetes operators.

![Version: CHART_VERSION](https://img.shields.io/badge/Version-CHART_VERSION-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: APP_VERSION](https://img.shields.io/badge/AppVersion-APP_VERSION-informational?style=flat-square)

## Get Repo Info

```console
helm repo add kkb0318 https://kkb0318.github.io/kom
helm repo update
```

## Install Chart

```console
helm install kom kkb0318/kom -n kom-system --create-namespace
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controllerManager.manager.args[0] | string | `"--leader-elect"` |  |
| controllerManager.manager.containerSecurityContext.allowPrivilegeEscalation | bool | `false` |  |
| controllerManager.manager.containerSecurityContext.capabilities.drop[0] | string | `"ALL"` |  |
| controllerManager.manager.image.repository | string | `"ghcr.io/kkb0318/kom"` |  |
| controllerManager.manager.image.tag | string | `"APP_VERSION"` |  |
| controllerManager.manager.resources.limits.cpu | string | `"500m"` |  |
| controllerManager.manager.resources.limits.memory | string | `"128Mi"` |  |
| controllerManager.manager.resources.requests.cpu | string | `"10m"` |  |
| controllerManager.manager.resources.requests.memory | string | `"64Mi"` |  |
| controllerManager.replicas | int | `1` |  |
| controllerManager.serviceAccount.annotations | object | `{}` |  |
| kubernetesClusterDomain | string | `"cluster.local"` |  |
| metricsService.ports[0].name | string | `"https"` |  |
| metricsService.ports[0].port | int | `8443` |  |
| metricsService.ports[0].protocol | string | `"TCP"` |  |
| metricsService.ports[0].targetPort | string | `"https"` |  |
| metricsService.type | string | `"ClusterIP"` |  |
