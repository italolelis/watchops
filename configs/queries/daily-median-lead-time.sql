SELECT
  day,
  NVL(ANY_VALUE(med_time_to_change)/60, 0) AS median_time_to_change
FROM (
  SELECT
    d.deploy_id,
    TRUNC(d.time_created) AS day,
    PERCENTILE_CONT(0.5) within group (order by 
   CASE WHEN date_DIFF('minute', c.time_created, d.time_created) > 0 THEN date_DIFF('minute', c.time_created, d.time_created) ELSE NULL END) -- I inverted c and d 
      OVER (PARTITION BY TRUNC(d.time_created)) AS med_time_to_change
  FROM watchops.deployments d --,d.changes 
  LEFT JOIN watchops.changes c 
  -- ON changes = c.change_id
  ON d.main_commit = c.change_id
)
GROUP BY day ORDER BY day;
