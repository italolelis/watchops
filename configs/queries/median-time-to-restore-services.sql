SELECT
CASE WHEN med_time_to_resolve < 24  then 'One day'
     WHEN med_time_to_resolve < 168  then 'One week'
     WHEN med_time_to_resolve < 730  then 'One month'
     WHEN med_time_to_resolve < 730 * 6 then 'Six months'
     ELSE 'One year'
     END as med_time_to_resolve
FROM (
  SELECT
  PERCENTILE_CONT(0.5) within group (order by date_diff('hour', time_resolved, time_created))
    OVER() as med_time_to_resolve
  FROM watchops.incidents
  WHERE time_created > dateadd(month,-3,GETDATE())
  )
LIMIT 1;
