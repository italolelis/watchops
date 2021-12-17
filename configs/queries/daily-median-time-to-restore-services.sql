SELECT
  TIMESTAMP_TRUNC(time_created, DAY) as day,
  -- Median time to resolve
  PERCENTILE_CONT(
        TIMESTAMP_DIFF(time_resolved, time_created, HOUR), 0.5)
    OVER(PARTITION BY TIMESTAMP_TRUNC(time_created, DAY)
    ) as daily_med_time_to_restore,
  FROM four_keys.incidents
ORDER BY day
