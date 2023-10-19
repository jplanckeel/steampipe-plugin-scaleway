# Table: scaleway_cockpit_plans

Plan is list of pricing for Cockpit in Scaleway.

## Examples

### Basic info

```sql
select
  name,
  retention_metrics_interval_days,
  retention_logs_interval_days,
  sample_ingestion_price,
  logs_ingestion_price,
  retention_price
from
  scaleway_cockpit_plans
```

### List Plan where name is free

```sql
select
  name,
  retention_metrics_interval_days,
  retention_logs_interval_days,
  sample_ingestion_price,
  logs_ingestion_price,
  retention_price
from
  scaleway_cockpit_plans
where name = 'free';
```
