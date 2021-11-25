CREATE TABLE IF NOT EXISTS kinesis_consumer (
	namespace VARCHAR NOT NULL,
	shard_id VARCHAR NOT NULL,
	sequence_number numeric NOT NULL,
	CONSTRAINT kinesis_consumer_pk PRIMARY KEY (namespace, shard_id)
);
