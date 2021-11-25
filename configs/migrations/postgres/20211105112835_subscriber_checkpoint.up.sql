CREATE TABLE IF NOT EXISTS kinesis_consumer (
	namespace text NOT NULL,
	shard_id text NOT NULL,
	sequence_number numeric NOT NULL,
	CONSTRAINT kinesis_consumer_pk PRIMARY KEY (namespace, shard_id)
);
