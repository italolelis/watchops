CREATE OR REPLACE VIEW watchops.deployments AS 
   WITH deploys AS (
      SELECT 
      source,
      metadata->'deployment'->>'id' as deploy_id,
      time_created,
      CASE WHEN source like '%circleci%' then metadata->'workflow'->>'id'
           WHEN source like '%github%' then metadata->'deployment'->>'sha'
           WHEN source like '%gitlab%' then metadata->'commit'->>'id' end as main_commit              
      FROM watchops.events_raw
      WHERE ((source like '%circleci%' and metadata->>'status' = 'SUCCESS')
            OR (source like 'github%' and event_type = 'deployment_status' and metadata->'deployment_status'->>'state' = 'success')
            OR (source like 'gitlab%' and event_type = 'pipeline' and metadata->'object_attributes'->>'status' = 'success'))
    )
    select * from deploys
