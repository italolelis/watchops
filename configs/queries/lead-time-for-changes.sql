SELECT
day, 
IFNULL(ANY_VALUE(med_time_to_change)/60, 0) AS median_time_to_change,
FROM (
    SELECT 
        d.deploy_id,
        TIMESTAMP_TRUNC(d.time_created, DAY) AS day,
        PERCENTILE_CONT(
            IF(TIMESTAMP_DIFF(d.time_created, c.time_created, MINUTE) > 0, TIMESTAMP_DIFF(d.time_created, c.time_created, MINUTE), NULL), 0.5) 
            OVER (PARTITION BY TIMESTAMP_TRUNC(d.time_created, DAY)) AS med_time_to_change, 
             FROM four_keys.deployments d, d.changes
             LEFT JOIN four_keys.changes c ON changes = c.change_id
             )
GROUP BY day
ORDER BY day
