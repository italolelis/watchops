CREATE OR REPLACE VIEW watchops.changes AS (
    WITH all_commits AS
    (
    SELECT source, event_type, json_parse(json_serialize(metadata)) AS metadata
    FROM watchops.events_raw
    WHERE event_type = 'push'
    and json_serialize(metadata) LIKE '%"commits":[%'                                         
    ),
    commit AS
    (
    SELECT index, i.source, i.event_type, element.timestamp AS time_created, element.id AS change_id
    FROM all_commits AS i, i.metadata.commits AS element AT index
    )
    SELECT source, event_type, change_id, date_trunc('second', TO_TIMESTAMP(CAST(time_created AS VARCHAR), 'YYYY-MM-DDTHH:MI:SS')) AS time_created 
    FROM commit
)

