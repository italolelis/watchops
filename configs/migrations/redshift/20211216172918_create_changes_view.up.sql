CREATE OR REPLACE VIEW watchops.changes AS (
    WITH all_commits AS
    (
        SELECT source, event_type, metadata 
        FROM watchops.events_raw
        WHERE event_type = 'push'                                    
    ),
    commit AS
    (
        SELECT index, i.source, i.event_type, element.timestamp AS time_created, element.id AS change_id
        FROM all_commits AS i, i.metadata.commits AS element AT index
    )
    SELECT source, event_type, TRIM('""' FROM CAST(change_id AS VARCHAR)) as change_id, date_trunc('second', TO_TIMESTAMP(CAST(time_created AS VARCHAR), 'YYYY-MM-DDTHH:MI:SS')::timestamp) AS time_created
    FROM commit
)

