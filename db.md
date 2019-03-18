# database

Uses a postgres database, configurable in `cmd/main.go`

```sql
CREATE TABLE files (
    id text PRIMARY KEY,
    serial SERIAL,
    owner text NOT NULL,
    s3 text NOT NULL,
    type text,
    labels text[]
    created_at timestamp with time zone NOT NULL DEFAULT now()
);
```
