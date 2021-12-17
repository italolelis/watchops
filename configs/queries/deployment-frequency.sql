WITH last_three_months AS
(SELECT
TIMESTAMP(day) AS day
FROM
UNNEST(
    GENERATE_DATE_ARRAY(
        DATE_SUB(CURRENT_DATE(), INTERVAL 3 MONTH), CURRENT_DATE(),
    INTERVAL 1 DAY)) AS day
-- FROM the start of the data
WHERE day > (SELECT date(min(time_created)) FROM four_keys.events_raw)
)
SELECT
CASE WHEN daily THEN \"Daily\" 
     WHEN weekly THEN \"Weekly\" 
      # If at least one per month, then Monthly
     WHEN PERCENTILE_CONT(monthly_deploys, 0.5) OVER () >= 1 THEN  \"Monthly\" 
     ELSE \"Yearly\"
     END as deployment_frequency
FROM (
      SELECT
  -- If the median number of days per week is more than 3, then Daily
  PERCENTILE_CONT(days_deployed, 0.5) OVER() >= 3 AS daily,
  -- If most weeks have a deployment, then Weekly
  PERCENTILE_CONT(week_deployed, 0.5) OVER() >= 1 AS weekly,
  -- Count the number of deployments per month.  
  -- Cannot mix aggregate and analytic functions, so calculate the median in the outer select statement
  SUM(week_deployed) OVER(partition by TIMESTAMP_TRUNC(week, MONTH)) monthly_deploys
  FROM(
          SELECT
      TIMESTAMP_TRUNC(last_three_months.day, WEEK) as week,
      MAX(if(deployments.day is not null, 1, 0)) as week_deployed,
      COUNT(distinct deployments.day) as days_deployed
      FROM last_three_months
      LEFT JOIN(
            SELECT
        TIMESTAMP_TRUNC(time_created, DAY) AS day,
        deploy_id
        FROM four_keys.deployments) deployments ON deployments.day = last_three_months.day
      GROUP BY week)
 )
LIMIT 1
