CREATE OR REPLACE VIEW watchops.incidents AS (
    WITH issue AS
    (
    SELECT 
        source,
        CASE WHEN source LIKE 'github%' THEN json_extract_path_text(json_serialize(metadata), 'issue', 'number')
            WHEN source LIKE 'gitlab%' AND event_type = 'note' THEN json_extract_path_text(json_serialize(metadata), 'object_attributes', 'noteable_id')
            WHEN source LIKE 'gitlab%' AND event_type = 'issue' THEN json_extract_path_text(json_serialize(metadata), 'object_attributes', 'id')
            WHEN source LIKE 'opsgenie%' THEN json_extract_path_text(json_extract_path_text(json_serialize(metadata), 'alert'), 'tinyId')
        end as incident_id,
        CASE WHEN source LIKE 'github%' THEN TO_TIMESTAMP(json_extract_path_text(json_serialize(metadata), 'issue', 'created_at'), 'YYYY-MM-DD HH24:MI:SS')
            WHEN source LIKE 'gitlab%' THEN TO_TIMESTAMP(json_extract_path_text(json_serialize(metadata), 'object_attributes', 'created_at'), 'YYYY-MM-DD HH24:MI:SS') 
            WHEN source LIKE 'opsgenie' THEN dateadd(ms, CAST(json_extract_path_text(json_extract_path_text(json_serialize(metadata), 'alert'), 'createdAt') AS bigint), '1970-01-01')
        end as time_created,
        CASE WHEN source LIKE 'github%' THEN TO_TIMESTAMP(json_extract_path_text(json_serialize(metadata), 'issue', 'closed_at'), 'YYYY-MM-DD HH24:MI:SS')
            WHEN source LIKE 'gitlab%' THEN TO_TIMESTAMP(json_extract_path_text(json_serialize(metadata), 'object_attributes', 'closed_at'), 'YYYY-MM-DD HH24:MI:SS') 
            WHEN source LIKE 'opsgenie' THEN dateadd(ms, CAST(json_extract_path_text(json_extract_path_text(json_serialize(metadata), 'alert'), 'updatedAt') AS bigint), '1970-01-01')
        end as time_resolved,
        CASE WHEN source LIKE 'github%' THEN json_extract_path_text(json_serialize(metadata), 'issue', 'labels') like '%bug%'
            WHEN source LIKE 'gitlab%' THEN json_extract_path_text(json_serialize(metadata), 'object_attributes', 'labels','title') like '%ncident%' 
            WHEN source LIKE 'opsgenie' THEN true
        end as bug
    FROM watchops.events_raw 
    WHERE event_type LIKE '%issue%' OR (event_type = 'note' and json_extract_path_text(json_serialize(metadata), 'object_attributes', 'noteable_type') = 'Issue') OR source = 'opsgenie'
    )
    SELECT
    source,
    incident_id,
    date_trunc('second',MIN(time_created)::timestamp) as time_created,
    date_trunc('second',MAX(time_resolved)::timestamp) as time_resolved
    FROM issue
    WHERE bug = true 
    and time_resolved >= TO_DATE('1900-01-01', 'YYYY-MM-DD') --filter null time_resolved 
    GROUP BY source, incident_id
)
