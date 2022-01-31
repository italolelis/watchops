SELECT
TRUNC(time_created) AS day,
COUNT(distinct deploy_id) AS deployments
FROM
watchops.deployments
GROUP BY day;
