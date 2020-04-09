# dflimg

## server

[DockerHub](https://hub.docker.com/r/duffleman/dflimg)

Quick file sharing platform. Accepts images, and files.

This is built from scratch in Go, so you need a public facing server to handle requests and the server will need to run Go with postgres (see [db.md](db.md)).

Requires an AWS account, it uses S3 as a host for the uploaded files.

When you run this, the shorter your domain is, the better.

Inspired by starbs/yeh

### Env variables to set

```
PG_OPTS=postgresql://postgres/dflimg?sslmode=prefer
DFL_USERS={"USERNAME": "PASSWORD"}
DFL_ROOT_URL=https://dfl.mn
DFL_SALT=some-long-string-that-works-as-a-salt-for-the-hasher
ADDR=:8001
AWS_ACCESS_KEY_ID=AWSKEY
AWS_SECRET_ACCESS_KEY=AWSSECRET
AWS_DEFAULT_REGION=AWSREGION
```

### Endpoints

#### `POST /upload_file`

Legacy, really only here to support programs like ShareX on Windows.

Takes a file in the form of multipart/form-data, returns  a short URL that links to the file. You can set the "Accept" header to modify the response. Defaults to JSON for the response.

##### Request

```bash
curl -X POST -H "Authorization: test" -F file=@duffleman.png https://dfl.mn/upload_file
```

##### Response

```json
{
    "resource_id": "file_000000BdAf7MWsYZ6r5wc18cV2sAS",
    "type": "file",
    "hash": "q3A",
    "url": "https://dfl.mn/q3A"
}
```

If the "Accept" header is set to "text/plain":

`https://dfl.mn/q3A`

#### `POST /create_signed_url`

##### Request

```json
{
	"content_type": "image/png"
}
```

##### Response

```json
{
	"resource_id": "file_aaa000",
	"type": "file",
	"hash": "xAx",
	"url": "https://dfl.mn/xAx",
	"s3link": "https://s3.amazon.com/eu-west-1/..."
}
```

You must then post the content of the file to the S3 link returned to you.

#### `POST /delete_resource`

```json
{
	"query": "aVA"
}
```

#### `POST /set_nsfw`

##### Request

```json
{
	"query": "aAb",
	"nsfw": true
}
```

#### `POST /add_shortcut`

##### Request

```json
{
	"query": "aCw",
	"shortcut": "scott"
}
```

#### `POST /remove_shortcut`

##### Request

```json
{
	"query": "aCw",
	"shortcut": "scott"
}
```

#### `POST /shorten_url`

Shorten a URL. It requires `url` which is the URL to shorten.

##### Request

```json
{
	"url": "https://google.com"
}
```

##### Response

```json
{
    "resource_id": "url_000000BdAf7MWsYZ6r5wc18cV2sAS",
    "type": "url",
    "hash": "aaB",
    "url": "https://dfl.mn/aaB"
}
```

#### `GET /{hash}`

Links to the resource. Serves the content directly!

#### `GET /:{shortcut}`

Links to the resource through one of it's shortcuts. Serves the content directly!

## client/cli

A CLI tool that allows you to upload files to the above server! More information on this soon.

### Install

Install it into your PATH

`go install cmd/dflimg/...`

### Env variables to set

```
DFLIMG_ROOT_URL=https://dfl.mn
DFLIMG_AUTH_TOKEN=some-token-thats-on-the-server
```

### Upload a file

Upload a single file:

`dflimg u {file}`

Upload a file with some shorcuts, you can give it a CSV for the shortcuts (`-s`)

`dflimg u -s test,srs {file}`

Upload a file, mark it as NSFW (`-n`)

`dflimg u -n {file}`

It will attempt to automatically put the URL in your clipboard too!

### Shorten URL

Shorten a URL

`dflimg s {url}`

See other params above.

### Copy a URL

When given a long URL leading to an image, it'll attempt to download the file and reupload it to the dflimg server.

`dflimg c {url}`

`-n` for NSFW works here, along with `-s` for shortcuts.
