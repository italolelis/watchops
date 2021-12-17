SELECT 
 CASE
  WHEN median_time_to_change < 24 * 60 then \"One day\"
  WHEN median_time_to_change < 168 * 60 then \"One week\"
  WHEN median_time_to_change < 730 * 60 then \"One month\"
  WHEN median_time_to_change < 730 * 6 * 60 then \"Six months\"
  ELSE \"One year\"
 END as lead_time_to_change
FROM (
         SElECT
  IFNULL(ANY_VALUE(med_time_to_change), 0) AS median_time_to_change
 FROM (
           SELECT
   PERCENTILE_CONT(
            IF(TIMESTAMP_DIFF(d.time_created, c.time_created, MINUTE) > 0, TIMESTAMP_DIFF(d.time_created, c.time_created, MINUTE), NULL),
    0.5)
    OVER () AS med_time_to_change, # Minutes
  FROM four_keys.deployments d, d.changes
  LEFT JOIN four_keys.changes c ON changes = c.change_id
  WHERE d.time_created > TIMESTAMP(DATE_SUB(CURRENT_DATE(), INTERVAL 3 MONTH))
 )
)
