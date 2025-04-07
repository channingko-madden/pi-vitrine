# DB Setup

## Postgres

```bash
sudo -u postgres createuser -P -d pi-vitrine
createdb pi_vitrine

psql -U pi-vitrine -f testdata/setup.sql -d pi_vitrine
```

Make sure that the pi-vitrine user is the owner of the pi_vitrine tables,
and has read/write permissions.

The host server will automatially run db migrations on startup.

### Notes
```bash
GRANT CONNECT ON DATABASE pi_vitrine to "pi-vitrine";

### postgres 16
GRANT pg_read_all_data TO "pi-vitrine";
GRANT pg_write_all_data TO "pi-vitrine";

### postgres 13
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "pi-vitrine";
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO "pi-vitrine";
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO "pi-vitrine";
GRANT ALL PRIVILEGES ON DATABASE pi_vitrine TO "pi-vitrine";
```
