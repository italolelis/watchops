SELECT
CASE WHEN change_fail_rate <= .15 then '0-15%'
     WHEN change_fail_rate < .46 then '16-45%'
     ELSE '46-60%' end as change_fail_rate
FROM 
 (SELECT
  SUM( CASE WHEN i.incident_id is NULL THEN 0 ELSE 1 END) / COUNT(DISTINCT change_id) as change_fail_rate
  FROM watchops.deployments d
  LEFT JOIN watchops.changes c ON d.main_commit = c.change_id
  LEFT JOIN(SELECT
          incident_id,
          changes,
          time_resolved
          FROM watchops.incidents) i ON i.changes = d.main_commit
  WHERE d.time_created > dateadd(month,-3,GETDATE())
  )
LIMIT 1;
