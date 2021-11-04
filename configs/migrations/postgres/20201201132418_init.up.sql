CREATE SCHEMA fourkeys;

CREATE TABLE fourkeys.events_raw (
    id VARCHAR NOT NULL,
    event_type VARCHAR NOT NULL,
    metadata JSONB NOT NULL,
    time_created TIMESTAMP NOT NULL,
    signature VARCHAR NOT NULL,
    msg_id  VARCHAR NOT NULL,
    source  VARCHAR NOT NULL,
    CONSTRAINT "pk_events_raw_id" PRIMARY KEY ("signature")
);
