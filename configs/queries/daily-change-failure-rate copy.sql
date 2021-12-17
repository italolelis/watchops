SELECT
CASE WHEN change_fail_rate <= .15 then \"0-15%\"
     WHEN change_fail_rate < .46 then \"16-45%\"
     ELSE \"46-60%\" 
     END AS change_fail_rate
FROM 
 (SELECT
  SUM(IF(i.incident_id is NULL, 0, 1)) / COUNT(DISTINCT change_id) as change_fail_rate
  FROM four_keys.deployments d, d.changes
  LEFT JOIN four_keys.changes c ON changes = c.change_id
  LEFT JOIN(SELECT
          incident_id,
          change,
          time_resolved
          FROM four_keys.incidents i,
          i.changes change) i ON i.change = changes
  WHERE d.time_created > TIMESTAMP(DATE_SUB(CURRENT_DATE(), INTERVAL 3 MONTH))
  )
LIMIT 1
