SELECT
  TRUNC(time_created) as day,
  PERCENTILE_CONT(0.5) within group (order by date_diff('hour', time_created, time_resolved))
    OVER(PARTITION BY TRUNC(time_created)
    ) as daily_med_time_to_restore
  FROM watchops.incidents;
