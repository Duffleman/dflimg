# database

Uses a postgres database, configurable in `cmd/main.go`

```sql
CREATE TABLE resources (
    id text PRIMARY KEY,
    type text NOT NULL,
    serial BIGSERIAL,
    owner text NOT NULL,
    link text NOT NULL,
    nsfw bool NOT NULL DEFAULT false,
    mime_type text,
    shortcuts text[],
    created_at timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE labels (
    id text PRIMARY KEY,
    name text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE labels_resources (
    label_id text NOT NULL,
    resource_id text NOT NULL,
    PRIMARY KEY("label_id", "resource_id")
);
```
