CREATE TABLE IF NOT EXISTS kinesis_consumer (
	namespace varchar(255) NOT NULL,
	shard_id varchar(255) NOT NULL,
	sequence_number numeric(65,0) NOT NULL,
	CONSTRAINT kinesis_consumer_pk PRIMARY KEY (namespace, shard_id)
);
