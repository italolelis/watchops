CREATE VIEW watchops.deployments AS 
   WITH deploys AS (
      SELECT 
      source,
      json_extract_path_text(json_serialize(metadata), 'deployment', 'id') as deploy_id,
      time_created,
      CASE WHEN source like '%circleci%' then json_extract_path_text(json_serialize(metadata), 'workflow', 'id')
           WHEN source like '%github%' then json_extract_path_text(json_serialize(metadata), 'deployment', 'sha')
           WHEN source like '%gitlab%' then json_extract_path_text(json_serialize(metadata), 'commit', 'id') end as main_commit              
      FROM watchops.events_raw
      WHERE ((source like '%circleci%' and json_extract_path_text(json_serialize(metadata), 'status') = 'SUCCESS')
            OR (source like 'github%' and event_type = 'deployment_status' and json_extract_path_text(json_serialize(metadata), 'deployment_status', 'state') = 'success')
            OR (source like 'gitlab%' and event_type = 'pipeline' and json_extract_path_text(json_serialize(metadata), 'object_attributes', 'status') = 'success'))
    )
    select * from deploys
