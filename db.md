# database

Uses a postgres database, configurable in `cmd/main.go`

```sql
CREATE TABLE resources (
    id text PRIMARY KEY,
    serial SERIAL UNIQUE,
    hash text,
    type text NOT NULL,
    owner text NOT NULL,
    link text NOT NULL,
    mime_type text,
    shortcuts text[],
    nsfw boolean NOT NULL DEFAULT false,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone
);

CREATE TABLE labels (
    id text PRIMARY KEY,
    name text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE labels_resources (
    label_id text,
    resource_id text,
    CONSTRAINT labels_resources_pkey PRIMARY KEY (label_id, resource_id)
);
```
