CREATE OR REPLACE VIEW watchops.deployments AS 
   WITH deploys AS (
      SELECT 
      source,
      cast(metadata."deployment"."id" as varchar) as deploy_id,
      time_created,
      CASE WHEN source like '%circleci%' then cast(metadata."workflow"."id" as varchar)
           WHEN source like '%github%' then cast(metadata."deployment"."sha" as varchar)
           WHEN source like '%gitlab%' then cast(metadata."commit"."id" as varchar)
      end as main_commit              
      FROM watchops.events_raw
      WHERE ((source like '%circleci%' and cast(metadata."status" as varchar) = 'SUCCESS'
            OR (source like 'github%' and event_type = 'deployment_status' and cast(metadata."deployment_status"."state" as varchar) = 'success')
            OR (source like 'gitlab%' and event_type = 'pipeline' and cast(metadata."object_attributes"."status" as varchar) = 'success')))
    )
    select * from deploys
