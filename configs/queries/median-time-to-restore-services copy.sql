SELECT
CASE WHEN med_time_to_resolve < 24  then \"One day\"
     WHEN med_time_to_resolve < 168  then \"One week\"
     WHEN med_time_to_resolve < 730  then \"One month\"
     WHEN med_time_to_resolve < 730 * 6 then \"Six months\"
     ELSE \"One year\"
     END as med_time_to_resolve,
FROM (
    SELECT
  -- Median time to resolve
  PERCENTILE_CONT(
      TIMESTAMP_DIFF(time_resolved, time_created, HOUR), 0.5)
    OVER() as med_time_to_resolve,
  FROM four_keys.incidents
  -- Limit to 3 months
  WHERE time_created > TIMESTAMP(DATE_SUB(CURRENT_DATE(), INTERVAL 3 MONTH)))
LIMIT 1
