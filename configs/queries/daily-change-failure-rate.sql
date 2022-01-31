SELECT
TRUNC(d.time_created) as day,
SUM(CASE WHEN i.incident_id is NULL THEN 0 ELSE 1 END) / COUNT(DISTINCT change_id) as change_fail_rate
FROM watchops.deployments d --, d.changes
LEFT JOIN watchops.changes c 
-- ON changes = c.change_id
ON main_commit = c.change_id
LEFT JOIN(SELECT
        incident_id,
        changes,
        time_resolved
        FROM watchops.incidents) i 
        ON i.changes = change_id
GROUP BY day;
