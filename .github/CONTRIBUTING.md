
```sh
pip install awscli-local
```

```sh
awslocal kinesis create-stream --stream-name fourkeys_github --shard-count 1
awslocal kinesis create-stream --stream-name fourkeys_opsgenie --shard-count 1
```
