select date_part('year', events_raw.time_created) as year, extract('week' from events_raw.time_created) as week, date_part('year', events_raw.time_created) || '.' || extract('week' from events_raw.time_created) as week_year, json_extract_path_text(json_serialize(metadata), 'deployment', 'status', 'state') as state, count(*)
from fourkeys.events_raw 
where event_type = 'deployment_status' 
and json_extract_path_text(json_serialize(metadata), 'deployment', 'status', 'state') in ('success', 'failure')
[[and json_extract_path_text(json_serialize(metadata), 'deployment', 'environment') = {{environment}}]]
[[and {{timesframe}}]]
group by year, week, state
order by year, week;
