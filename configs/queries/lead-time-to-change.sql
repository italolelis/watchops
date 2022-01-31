SELECT 
  CASE
    WHEN median_time_to_change < 24 * 60 then 'One day'
    WHEN median_time_to_change < 168 * 60 then 'One week'
    WHEN median_time_to_change < 730 * 60 then 'One month'
    WHEN median_time_to_change < 730 * 6 * 60 then 'Six months'
    ELSE 'One year'
    END as lead_time_to_change
FROM (
  SELECT
    NVL(ANY_VALUE(med_time_to_change), 0) AS median_time_to_change
  FROM (
    SELECT
      PERCENTILE_CONT(0.5) within group (order by 
      CASE WHEN date_diff('minute', c.time_created, d.time_created) > 0 
           THEN date_diff('minute', c.time_created, d.time_created)
		   ELSE NULL -- I inverted c with d 
	  END
     )
        OVER () AS med_time_to_change
    FROM watchops.deployments d --, d.changes
    LEFT JOIN watchops.changes c 
    -- ON changes = c.change_id
    ON d.main_commit = c.change_id
    WHERE d.time_created > dateadd(month,-3,GETDATE()) 
  )
)

