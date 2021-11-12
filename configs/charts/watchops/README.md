# Watchops
 Helm Chart

![Version: 1.2.3](https://img.shields.io/badge/Version-1.2.3-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: v1.0.0](https://img.shields.io/badge/AppVersion-v1.0.0-informational?style=flat-square)

This repository implements the four key metrics coined by Google and introduced in Accelerate: Time to Restore, Lead Time to Change, Deployment Frequency, and Change Failure Rate.

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| italolelis | me@italovietro.com |  |

## Source Code

* <https://github.com/italolelis/watchops>

## Get Repo Info

```console
$ helm repo add italolelis https://italolelis.github.io/watchops
$ helm repo update
```

_See [helm repo](https://helm.sh/docs/helm/helm_repo/) for command documentation._

## Installing the Chart

To install the chart with the release name `my-release`:

```console
helm install my-release italolelis/watchops
```

## Uninstalling the Chart

To uninstall/delete the my-release deployment:

```console
helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

# Configuration

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| app.affinity | object | `{}` |  |
| app.annotations | object | `{}` |  |
| app.config.githubSecret | string | `""` | Sets the github secret that will be use to validate incoming webhooks from GitHub. |
| app.config.logLevel | string | `"info"` | Define the log level. Accepted values are: debug, info, warn, error. |
| app.config.messageBroker.singleTopic | bool | `false` | Whether to use a single topic for all incomming webhooks or not. |
| app.config.messageBroker.topicPrefix | string | `"watchops_"` | If you defined multiple topics (one for each incoming webhook type), then you can define the prefix of these topics. |
| app.config.opsgenieSecret | string | `""` | Sets the OpsGenie secret that will be used to validate incoming webhooks from OpsGenie. |
| app.config.port | int | `8080` | Configure the port number of the main API. |
| app.config.rest.idleTimeout | string | `"30s"` | Defines the idle server timeout. |
| app.config.rest.readTimeout | string | `"30s"` | Defines the read server timeout. |
| app.config.rest.writeTimeout | string | `"30s"` | Defines the write server timeout. |
| app.enabled | bool | `true` |  |
| app.image.repository | string | `"ghcr.io/italolelis/watchops"` |  |
| app.image.tag | string | `"latest"` |  |
| app.imagePullSecrets | list | `[]` |  |
| app.labels | object | `{}` |  |
| app.livenessProbe.failureThreshold | int | `3` |  |
| app.livenessProbe.httpGet.path | string | `"/live"` |  |
| app.livenessProbe.httpGet.port | string | `"http-probe"` |  |
| app.livenessProbe.initialDelaySeconds | int | `10` |  |
| app.livenessProbe.periodSeconds | int | `10` |  |
| app.livenessProbe.successThreshold | int | `1` |  |
| app.livenessProbe.timeoutSeconds | int | `1` |  |
| app.nodeSelector | object | `{}` |  |
| app.podAnnotations | object | `{}` |  |
| app.podSecurityContext | object | `{}` |  |
| app.readinessProbe.failureThreshold | int | `3` |  |
| app.readinessProbe.httpGet.path | string | `"/ready"` |  |
| app.readinessProbe.httpGet.port | string | `"http-probe"` |  |
| app.readinessProbe.initialDelaySeconds | int | `15` |  |
| app.readinessProbe.periodSeconds | int | `10` |  |
| app.readinessProbe.successThreshold | int | `1` |  |
| app.readinessProbe.timeoutSeconds | int | `1` |  |
| app.replicaCount | int | `1` |  |
| app.resources.requests.cpu | string | `"20m"` |  |
| app.resources.requests.memory | string | `"30Mi"` |  |
| app.securityContext | object | `{}` |  |
| app.tolerations | list | `[]` |  |
| autoscaling.app.enabled | bool | `false` | Enable autoscaling for the Four Keys main API. |
| autoscaling.app.maxReplicas | int | `3` | Number of maximum replicas to scale up. |
| autoscaling.app.minReplicas | int | `1` | Number of minimum replicas to scale down. |
| autoscaling.app.targetCPUUtilizationPercentage | int | `80` | Target CPU utilization percentage. |
| autoscaling.app.targetMemoryUtilizationPercentage | int | `50` | Target memory utilization percentage. |
| ingress.annotations | object | `{}` | Ingress annotations (values are templated) |
| ingress.enabled | bool | `false` | Enables Ingress |
| ingress.hosts[0].backend.serviceName | string | `"chart-example.local"` |  |
| ingress.hosts[0].backend.servicePort | int | `80` |  |
| ingress.hosts[0].host | string | `"chart-example.local"` |  |
| ingress.hosts[0].paths[0] | string | `"/"` |  |
| ingress.tls | list | `[]` | Ingress TLS configuration |
| migrations.enabled | bool | `false` | Enable migrations to your datasource. |
| networkPolicy.enabled | bool | `false` | Whether to enable network policies. If your cluster supports it, I recommend enabling it. |
| pdb.enabled | bool | `true` |  |
| pdb.minAvailable | int | `1` |  |
| prometheusRule.annotations | object | `{}` | PrometheusRule annotations |
| prometheusRule.enabled | bool | `false` | If enabled, a PrometheusRule resource for Prometheus Operator is created |
| prometheusRule.groups | list | `[]` | Contents of Prometheus rules file |
| prometheusRule.labels | object | `{}` | Additional PrometheusRule labels |
| prometheusRule.namespace | string | `nil` | Alternative namespace for the PrometheusRule resource |
| rbac.create | bool | `true` |  |
| rbac.pspEnabled | bool | `true` |  |
| service.port | int | `80` |  |
| service.probePort | int | `9090` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.annotations | object | `{}` | Annotations for the service account |
| serviceAccount.automountServiceAccountToken | bool | `false` | Set this toggle to false to opt out of automounting API credentials for the service account |
| serviceAccount.create | bool | `true` | Specifies whether a ServiceAccount should be created |
| serviceAccount.imagePullSecrets | list | `[]` | Image pull secrets for the service account |
| serviceAccount.name | string | `nil` | The name of the ServiceAccount to use. If not set and create is true, a name is generated using the fullname template |
| serviceMonitor.enabled | bool | `false` | Whether to enable prometheus service monitor. Enable this if you're using https://github.com/coreos/prometheus-operator. |
| subscribers.github.affinity | object | `{}` |  |
| subscribers.github.annotations | object | `{}` |  |
| subscribers.github.config.database.driver | string | `"postgres"` | The main database driver (the event store). You can choose between postgres or redshift. |
| subscribers.github.config.database.dsn | string | `""` | The DSN for the database. If you are using the postgres chart, you don't need to define this value. |
| subscribers.github.config.logLevel | string | `"info"` | Define the log level. Accepted values are: debug, info, warn, error. |
| subscribers.github.config.messageBroker.driver | string | `"kinesis"` | Defines which message broker to use. You can choose between kinesis or awslambda (which is also based on kinesis). |
| subscribers.github.config.messageBroker.region | string | `"eu-central-1"` | When using AWS Kinesis, you need to set the AWS region. |
| subscribers.github.config.messageBroker.store.appName | string | `"watchops-consumer-github"` | The app name that is used as a namespace in the storage. |
| subscribers.github.config.messageBroker.store.driver | string | `"memory"` | The message broker storage that holds the last read record from the broker. You can choose between memory, postgres, mysql, or redis. I do not recommend using the memory store in production. |
| subscribers.github.config.messageBroker.store.mysql.dsn | string | `""` | The mysql database DSN. If you are using the mysql chart, you don't need to define this value. |
| subscribers.github.config.messageBroker.store.mysql.tableName | string | `"kinesis_consumer"` | The mysql table name. |
| subscribers.github.config.messageBroker.store.postgres.dsn | string | `""` | The postgres database DSN. If you are using the postgres chart, you don't need to define this value. |
| subscribers.github.config.messageBroker.store.postgres.tableName | string | `"kinesis_consumer"` | The postgres table name. |
| subscribers.github.config.messageBroker.store.redis.address | string | `""` | The redis address. If you are using the redis chart, you don't need to define this value. |
| subscribers.github.config.messageBroker.store.redis.db | string | `""` | The redis database. If you are using the redis chart, you don't need to define this value. |
| subscribers.github.config.messageBroker.store.redis.password | string | `""` | The redis password. If you are using the redis chart, you don't need to define this value. |
| subscribers.github.config.messageBroker.store.redis.username | string | `""` | The redis username. If you are using the redis chart, you don't need to define this value. |
| subscribers.github.config.messageBroker.streamName | string | `"watchops_github"` | Sets the name of the stream this subscriber will listen to. If you are using a single stream for all webhook types, then just define the name of that stream. |
| subscribers.github.enabled | bool | `true` |  |
| subscribers.github.image.repository | string | `"ghcr.io/italolelis/watchops"` |  |
| subscribers.github.image.tag | string | `"subscriber-latest"` |  |
| subscribers.github.imagePullSecrets | list | `[]` |  |
| subscribers.github.labels | object | `{}` |  |
| subscribers.github.livenessProbe | object | `{}` |  |
| subscribers.github.nodeSelector | object | `{}` |  |
| subscribers.github.podAnnotations | object | `{}` |  |
| subscribers.github.podSecurityContext | object | `{}` |  |
| subscribers.github.readinessProbe | object | `{}` |  |
| subscribers.github.replicaCount | int | `1` |  |
| subscribers.github.resources.requests.cpu | string | `"20m"` |  |
| subscribers.github.resources.requests.memory | string | `"30Mi"` |  |
| subscribers.github.securityContext | object | `{}` |  |
| subscribers.github.tolerations | list | `[]` |  |
