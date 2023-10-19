# Table: scaleway_cockpit_grafana_users

A Grafana user is any individual who can log in to Grafana. Each user is associated with a role. There are two types of roles a user can have:

a viewer: can only view dashboards
an editor: can build and view dashboards
Managed dashboards in the “Scaleway” folder are always read-only, regardless of your role.

This table requires an Organization ID to be configured in the scaleway.spc file. Because we use table scaleway_project in ParentHydrate. 

## Examples

### Basic info

```sql
select
  login,
  project,
  role
from
  scaleway_cockpit_grafana_users;
```

### List Grafana Users with role editor

```sql
select
  login,
  project,
  role
from
  scaleway_cockpit_grafana_users
where role = 'editor';
```

### List Grafana Users for a project

```sql
select
   login,
   role,
   project as project_id,
   p.name as project_name 
from
   scaleway_cockpit_grafana_users as c 
   inner join
      scaleway_project p 
      on c.project = p.project_id 
where
   p.name = 'default';
```
