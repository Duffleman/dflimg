# database

Uses a postgres database, configurable in `cmd/main.go`

```sql
CREATE TABLE files (
    id text PRIMARY KEY,
    serial BIGSERIAL,
    owner text NOT NULL,
    s3 text NOT NULL,
    type text,
    shortcuts text[],
    created_at timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE links (
    id text PRIMARY KEY,
    serial BIGSERIAL,
    owner text NOT NULL,
    url text NOT NULL,
    nsfw bool NOT NULL DEFAULT false,
    shortcuts text[],
    comment text,
    created_at timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE labels (
    id text PRIMARY KEY,
    name text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE file_labels (
    file_id text NOT NULL,
    label_id text NOT NULL,
    PRIMARY KEY("file_id", "label_id")
);

CREATE TABLE labels_links (
    link_id text NOT NULL,
    label_id text NOT NULL,
    PRIMARY KEY("link_id", "label_id")
);
```
