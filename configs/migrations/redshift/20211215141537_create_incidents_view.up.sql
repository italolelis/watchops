CREATE OR REPLACE VIEW watchops.incidents AS (
    WITH issue AS
    (
    SELECT 
        source,
        CASE WHEN source LIKE 'github%' THEN cast(metadata."issue"."number" as varchar) 
            WHEN source LIKE 'gitlab%' AND event_type = 'note' THEN cast(metadata."object_attributes"."noteable_id" as varchar)
            WHEN source LIKE 'gitlab%' AND event_type = 'issue' THEN cast(metadata."object_attributes"."id" as varchar)
            WHEN source LIKE 'opsgenie%' THEN json_extract_path_text(json_serialize(metadata), 'alert', 'tinyId')
        end as incident_id,
        CASE WHEN source LIKE 'github%' THEN TO_TIMESTAMP(cast(metadata."issue"."created_at" as varchar), 'YYYY-MM-DD HH24:MI:SS')
            WHEN source LIKE 'gitlab%' THEN TO_TIMESTAMP(cast(metadata."object_attributes"."created_at" as varchar), 'YYYY-MM-DD HH24:MI:SS') 
            WHEN source LIKE 'opsgenie' THEN dateadd(ms, CAST(json_extract_path_text(json_extract_path_text(json_serialize(metadata), 'alert'), 'createdAt') AS bigint), '1970-01-01')
        end as time_created,
        CASE WHEN source LIKE 'github%' THEN TO_TIMESTAMP(cast(metadata."issue"."closed_at" as varchar), 'YYYY-MM-DD HH24:MI:SS')
            WHEN source LIKE 'gitlab%' THEN TO_TIMESTAMP(cast(metadata."object_attributes"."closed_at" as varchar), 'YYYY-MM-DD HH24:MI:SS') 
            WHEN source LIKE 'opsgenie' THEN dateadd(ms, CAST(json_extract_path_text(json_extract_path_text(json_serialize(metadata), 'alert'), 'updatedAt') AS bigint)/1000000, '1970-01-01')
        end as time_resolved,
        regexp_substr(json_serialize(metadata), '(?<=root cause: )(.*?)(?=\")', 1, 1,'p') as root_cause,
        CASE WHEN source LIKE 'github%' THEN cast(metadata."issue"."labels" as varchar) like '%bug%'
            WHEN source LIKE 'gitlab%' THEN cast(metadata."object_attributes"."labels"."title" as varchar) like '%ncident%' 
            WHEN source LIKE 'opsgenie' THEN true
        end as bug
    FROM watchops.events_raw 
    WHERE 
    event_type LIKE '%issue%' 
    OR (event_type = 'note' and cast(metadata."object_attributes"."noteable_type" as varchar) = 'Issue') 
    OR source = 'opsgenie'
    )
    SELECT
    issue.source,
    issue.incident_id,
    date_trunc('second',MIN(issue.time_created)::timestamp) as time_created,
    date_trunc('second',MAX(issue.time_resolved)::timestamp) as time_resolved,
  	LISTAGG(issue.root_cause IGNORE NULLS) changes
    FROM issue 
    LEFT JOIN watchops.deployments d on d.main_commit = root_cause
    WHERE bug = true 
    and time_resolved >= TO_DATE('1900-01-01', 'YYYY-MM-DD') --filter null time_resolved 
    GROUP BY issue.source, issue.incident_id
)
