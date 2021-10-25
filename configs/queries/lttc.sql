select events_raw.time_created, extract(epoch from (select pe.time_created as base_ref from fourkeys.events_raw pe where pe.event_type = 'push' and events_raw.metadata -> 'pull_request' ->> 'head_ref' = SUBSTRING(pe.metadata -> 'push' ->> 'ref', 12) order by pe.time_created desc limit 1) - events_raw.time_created)/3600 as time_to_merge
    from fourkeys.events_raw
    where events_raw.event_type = 'pull_request' 
        and events_raw.metadata ->> 'action' = 'closed' 
        and events_raw.metadata -> 'pull_request' ->> 'merged' = 'true'
        and (select pe.time_created as base_ref from fourkeys.events_raw pe where pe.event_type = 'push' and events_raw.metadata -> 'pull_request' ->> 'head_ref' = SUBSTRING(pe.metadata -> 'push' ->> 'ref', 12) order by pe.time_created desc limit 1) is not null
        [[and {{created_at}}]]
;
