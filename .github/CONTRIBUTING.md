
```sh
pip install awscli-local
```

```sh
awslocal kinesis create-stream --stream-name watchops_github --shard-count 1
awslocal kinesis create-stream --stream-name watchops_opsgenie --shard-count 1
```

Gets the shard iterator

```
awslocal kinesis get-shard-iterator --stream-name watchops_github --shard-iterator-type LATEST --shard-id shardId-000000000000
```


Get records from the stream

```
awslocal kinesis get-records --shard-iterator <whatever-iterator-from-previous-command>
```
