# watchops
 Helm Chart

![Version: 1.2.4](https://img.shields.io/badge/Version-1.2.4-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: v1.0.0](https://img.shields.io/badge/AppVersion-v1.0.0-informational?style=flat-square)

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
| autoscaling.publisher.enabled | bool | `false` | Enable autoscaling for the Four Keys main API. |
| autoscaling.publisher.maxReplicas | int | `3` | Number of maximum replicas to scale up. |
| autoscaling.publisher.minReplicas | int | `1` | Number of minimum replicas to scale down. |
| autoscaling.publisher.targetCPUUtilizationPercentage | int | `80` | Target CPU utilization percentage. |
| autoscaling.publisher.targetMemoryUtilizationPercentage | int | `50` | Target memory utilization percentage. |
| ingress.annotations | object | `{}` | Ingress annotations (values are templated) |
| ingress.enabled | bool | `false` | Enables Ingress |
| ingress.hosts[0].backend.serviceName | string | `"chart-example.local"` |  |
| ingress.hosts[0].backend.servicePort | int | `80` |  |
| ingress.hosts[0].host | string | `"chart-example.local"` |  |
| ingress.hosts[0].paths[0] | string | `"/"` |  |
| ingress.tls | list | `[]` | Ingress TLS configuration |
| migrations.enabled | bool | `false` | Enable migrations to your datasource. |
| migrations.labels | object | `{}` |  |
| networkPolicy.enabled | bool | `false` | Whether to enable network policies. If your cluster supports it, I recommend enabling it. |
| pdb.enabled | bool | `true` |  |
| pdb.minAvailable | int | `1` |  |
| prometheusRule.annotations | object | `{}` | PrometheusRule annotations |
| prometheusRule.enabled | bool | `false` | If enabled, a PrometheusRule resource for Prometheus Operator is created |
| prometheusRule.groups | list | `[]` | Contents of Prometheus rules file |
| prometheusRule.labels | object | `{}` | Additional PrometheusRule labels |
| prometheusRule.namespace | string | `nil` | Alternative namespace for the PrometheusRule resource |
| publisher.affinity | object | `{}` |  |
| publisher.annotations | object | `{}` |  |
| publisher.config.githubSecret | string | `""` | Sets the github secret that will be use to validate incoming webhooks from GitHub. |
| publisher.config.logLevel | string | `"info"` | Define the log level. Accepted values are: debug, info, warn, error. |
| publisher.config.messageBroker.driver | string | `"kinesis"` | Defines which message broker to use. You can choose between kinesis or awslambda (which is also based on kinesis). |
| publisher.config.messageBroker.singleTopic | bool | `false` | Whether to use a single topic for all incomming webhooks or not. |
| publisher.config.messageBroker.topicPrefix | string | `"watchops_"` | If you defined multiple topics (one for each incoming webhook type), then you can define the prefix of these topics. |
| publisher.config.opsgenieSecret | string | `""` | Sets the OpsGenie secret that will be used to validate incoming webhooks from OpsGenie. |
| publisher.config.port | int | `8080` | Configure the port number of the main API. |
| publisher.config.rest.idleTimeout | string | `"30s"` | Defines the idle server timeout. |
| publisher.config.rest.readTimeout | string | `"30s"` | Defines the read server timeout. |
| publisher.config.rest.writeTimeout | string | `"30s"` | Defines the write server timeout. |
| publisher.enabled | bool | `true` |  |
| publisher.image.repository | string | `"ghcr.io/italolelis/watchops-publisher"` |  |
| publisher.image.tag | string | `"latest"` |  |
| publisher.imagePullSecrets | list | `[]` |  |
| publisher.labels | object | `{}` |  |
| publisher.livenessProbe.failureThreshold | int | `3` |  |
| publisher.livenessProbe.httpGet.path | string | `"/live"` |  |
| publisher.livenessProbe.httpGet.port | string | `"http-probe"` |  |
| publisher.livenessProbe.initialDelaySeconds | int | `10` |  |
| publisher.livenessProbe.periodSeconds | int | `10` |  |
| publisher.livenessProbe.successThreshold | int | `1` |  |
| publisher.livenessProbe.timeoutSeconds | int | `1` |  |
| publisher.nodeSelector | object | `{}` |  |
| publisher.podAnnotations | object | `{}` |  |
| publisher.podSecurityContext | object | `{}` |  |
| publisher.readinessProbe.failureThreshold | int | `3` |  |
| publisher.readinessProbe.httpGet.path | string | `"/ready"` |  |
| publisher.readinessProbe.httpGet.port | string | `"http-probe"` |  |
| publisher.readinessProbe.initialDelaySeconds | int | `15` |  |
| publisher.readinessProbe.periodSeconds | int | `10` |  |
| publisher.readinessProbe.successThreshold | int | `1` |  |
| publisher.readinessProbe.timeoutSeconds | int | `1` |  |
| publisher.replicaCount | int | `1` |  |
| publisher.resources.requests.cpu | string | `"20m"` |  |
| publisher.resources.requests.memory | string | `"30Mi"` |  |
| publisher.securityContext | object | `{}` |  |
| publisher.tolerations | list | `[]` |  |
| rbac.create | bool | `true` |  |
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
| subscribers.github.config.database.driver | string | `"postgres"` | The main database driver (the event store). You can choose between biquery, postgres or redshift. |
| subscribers.github.config.database.dsn | string | `""` | The DSN for the database. |
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
| subscribers.github.image.repository | string | `"ghcr.io/italolelis/watchops-subscriber"` |  |
| subscribers.github.image.tag | string | `"latest"` |  |
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
