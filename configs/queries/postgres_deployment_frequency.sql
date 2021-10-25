select extract('week' from time_created) as week, metadata -> 'deployment' -> 'status' ->> 'state' as state, count(*) 
from fourkeys.events_raw 
where event_type = 'deployment_status' 
and metadata -> 'deployment' -> 'status' ->> 'state' in ('success', 'failure')
[[and metadata -> 'deployment' ->> 'environment' = {{environment}}]]
[[and {{timesframe}}]]
group by week, state;
