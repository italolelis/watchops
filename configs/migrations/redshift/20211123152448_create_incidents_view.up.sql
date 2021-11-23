CREATE VIEW watchops.incidents AS (
    WITH issue AS
    (
    SELECT 
        source,
        CASE WHEN source LIKE 'github%' THEN json_extract_path_text(json_serialize(metadata), 'issue', 'number')
            WHEN source LIKE 'gitlab%' AND event_type = 'note' THEN json_extract_path_text(json_serialize(metadata), 'object_attributes', 'noteable_id')
            WHEN source LIKE 'gitlab%' AND event_type = 'issue' THEN json_extract_path_text(json_serialize(metadata), 'object_attributes', 'id') end as incident_id,
        CASE WHEN source LIKE 'github%' THEN TO_TIMESTAMP(json_extract_path_text(json_serialize(metadata), 'issue', 'created_at'), 'YYYY-MM-DD HH24:MI:SS')
            WHEN source LIKE 'gitlab%' THEN TO_TIMESTAMP(json_extract_path_text(json_serialize(metadata), 'object_attributes', 'created_at'), 'YYYY-MM-DD HH24:MI:SS') end as time_created,
        CASE WHEN source LIKE 'github%' THEN TO_TIMESTAMP(json_extract_path_text(json_serialize(metadata), 'issue', 'closed_at'), 'YYYY-MM-DD HH24:MI:SS')
            WHEN source LIKE 'gitlab%' THEN TO_TIMESTAMP(json_extract_path_text(json_serialize(metadata), 'object_attributes', 'closed_at'), 'YYYY-MM-DD HH24:MI:SS') end as time_resolved,
        CASE WHEN source LIKE 'github%' THEN json_extract_path_text(json_serialize(metadata), 'issue', 'labels') like '%bug%'
            WHEN source LIKE 'gitlab%' THEN json_extract_path_text(json_serialize(metadata), 'object_attributes', 'labels','title') like '%ncident%' end as bug
    FROM watchops.events_raw 
    WHERE event_type LIKE '%issue%' OR (event_type = 'note' and json_extract_path_text(json_serialize(metadata), 'object_attributes', 'noteable_type') = 'Issue')
    )
    SELECT
    source,
    incident_id,
    MIN(time_created) as time_created,
    MAX(time_resolved) as time_resolved
    FROM issue
    WHERE bug = true
    GROUP BY source, incident_id
)
