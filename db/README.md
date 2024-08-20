# DB Setup

## Postgres

```bash
postgres -D path/to/data

createuser -P -d pi-vitrine
createdb pi-vitrine

psql -U pi-vitrine -f testdata/setup.sql -d pi-vitrine
```
