# DB Setup

## Postgres

```bash
createdb pi_vitrine
sudo -u postgres createuser -P -d pi-vitrine

psql -U pi-vitrine -f testdata/setup.sql -d pi_vitrine

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
