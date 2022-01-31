WITH last_three_months AS (
SELECT (dateadd(month,-3,GETDATE())+ n)::date as day
from watchops.dates 
where day <=  GETDATE()
), 
deployments AS 
(
SELECT
TRUNC(time_created) AS day,
deploy_id
FROM
watchops.deployments
)
SELECT
CASE WHEN daily THEN 'Daily' 
     WHEN weekly THEN 'Weekly' 
     WHEN PERCENTILE_CONT(0.5) within group (order by monthly_deploys) OVER () >= 1 THEN  'Monthly' 
     ELSE 'Yearly'
     END as deployment_frequency
FROM (
  SELECT
  PERCENTILE_CONT(0.5) within group (order by days_deployed) OVER() >= 3 AS daily,
  PERCENTILE_CONT(0.5) within group (order by week_deployed) OVER() >= 1 AS weekly,
  SUM(week_deployed) OVER(partition by date_trunc('month', week)) monthly_deploys
  FROM(
      SELECT
      DATE_TRUNC('week', last_three_months.day) as week,
      MAX(CASE WHEN deployments.day is not null THEN 1 ELSE 0 END) as week_deployed,
      COUNT(distinct deployments.day) as days_deployed
	FROM last_three_months, deployments
	WHERE deployments.day = last_three_months.day
      GROUP BY week)
 )
LIMIT 1;  
