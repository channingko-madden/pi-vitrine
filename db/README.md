# DB Setup

## Postgres

```bash
postgres -D path/to/data

createdb pi_vitrine
createuser -P -d pi-vitrine

psql -U {superuser} -f testdata/setup.sql -d pi_vitrine

GRANT CONNECT ON DATABASE pi_vitrine to "pi-vitrine";
GRANT pg_read_all_data TO "pi-vitrine";
GRANT pg_write_all_data TO "pi-vitrine";
```
