# Table: scaleway_cockpit_contact_points

Contact points define who is notified when an alert fires. Contact points include emails, Slack, on-call systems and texts. When an alert fires, all contact points are notified.

This table requires an Organization ID to be configured in the scaleway.spc file. Because we use table scaleway_project in ParentHydrate. 

## Examples

### Basic info

```sql
select
  email,
  project
from
  scaleway_cockpit_contact_points;
```

### List contact for a project

```sql
select
   email,
   project as project_id,
   p.name as project_name 
from
   scaleway_cockpit_contact_points as c 
   inner join
      scaleway_project p 
      on c.project = p.project_id 
where
   p.name = 'default';
```
