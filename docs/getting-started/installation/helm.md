# ⚓ Helm

### Requirements

* Kubernetes 1.16+ cluster
* Helm 3.0+

### Features

* Supported database backend: `PostgresSQL`, `MySQL, Redshift, Apache Kafka`
* Supported message brokers backend: `AWS Kinesis, Cloud Pub/Sub, Apache Kafka`&#x20;
* Autoscaling for `Publisher` provided
*   Monitoring:

    > * Prometheus metrics for WatchOps
* Automatic database migration after a new deployment

### Installing the Chart

[Helm](https://helm.sh) must be installed to use the charts. Please refer to Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm is set up properly, add the repo as follows:

```
$ helm repo add watchops https://italolelis.github.io/watchops
$ helm repo update
$ helm install my-release watchops/watchops --namespace watchops
```

### Upgrading the Chart

To upgrade the chart with the release name `watchops`:

```
helm upgrade watchops watchops/watchops --namespace watchops
```

\> To upgrade to a new version of the chart, run `helm repo update` first.



### Uninstalling the Chart

To uninstall/delete the `watchops` deployment:

```
helm delete watchops --namespace watchops
```
