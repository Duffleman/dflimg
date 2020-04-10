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

Respects the `Accept` request header.

#### `POST /create_signed_url`

##### Request

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

#### `POST /delete_resource`

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

#### `POST /add_shortcut`

##### Request

```json
{
	"query": "axA",
	"shortcut": "hello"
}
```

#### `POST /remove_shortcut`

##### Request

```json
{
	"query": "axA",
	"shortcut": "hello"
}
```

#### `POST /view_details`

##### Request

```json
{
	"query":  "dZM"
}
```

##### Response

```json
{
	"id": "file_000000BslGI66pAIjV27Uvh4ofWKG",
	"type": "file",
	"hash": "dZM",
	"owner": "Duffleman",
	"nsfw": true,
	"mime_type": "image/png",
	"shortcuts": [
		"hello"
	],
	"created_at": "2020-04-10T00:35:44.793661+01:00",
	"deleted_at": null
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

`dflimg signed-upload {file}`

`dflimg u my-file.png`

It will attempt to automatically put the URL in your clipboard too!

### Shorten URL

Shorten a URL

`dflimg shorten {url}`

`dflimg s https://google.com/?query=something-long`

See other params above.

### Copy a URL

When given a long URL leading to an image, it'll attempt to download the file and reupload it to the dflimg server.

`dflimg copy {url}`

`dflimg c mLd`

### Set it as NSFW

Set the file as NSFW so a NSFW primer appears before the content. The user must agree before they continue.

`dflimg nsfw {url or hash}`

`dflimg n ddA`

### Add a shortcut

Add a shortcut to the resource, so there is an easy way to access the resource

`dflimg add-shortcut {url or hash} {shortcut}`

`dflimg asc https://dfl.mn/aaA yolo`

### Remove a shortcut

Remove a shortcut from the resource

`dflimg remove-shortcut {url or hash} {shortcut}`

`dflimg rsc aaA yolo`

### Screenshot

macOS only so far, this one handles the whole screenshot process for you. Bind this to a shortcut on your mac so you can quickly take a snippet of a program and the link appears in your clipboard

`dflimg screenshot`
