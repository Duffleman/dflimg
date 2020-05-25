# dflimg

## server

[DockerHub](https://hub.docker.com/r/duffleman/dflimg)

Quick file sharing and URL shortening platform. Accepts images, files, pretty much anything.

This is built from scratch in Go, so you'll need to handle the dependancies yourself for now. It requires

- redis
- web-ingress (TLS termination)
- postgres (see [db.md](db.md))

When you run this, the shorter your domain is, the better.

Inspired by [starbs/yeh](https://github.com/starbs/yeh)

### Env variables to set

```bash
PG_OPTS=postgresql://postgres/dflimg?sslmode=prefer
DFL_USERS={"USERNAME": "PASSWORD"}
DFL_ROOT_URL=https://dfl.mn
DFL_SALT=some-long-string-that-works-as-a-salt-for-the-hasher
ADDR=:8001
```

### Endpoints

Any case where the response is ommited, the response should be a 204 (No content). Any case where the request is ommited, you are not expected to provide a body.

#### `POST /upload_file`

This endpoint should be used as little as you can, it's better to use `POST /create_signed_url` where you can. This is for cases when the storage provider cannot provide a signed URL.

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

Creates a signed URL to upload a file to directly.

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

### storage providers

#### aws s3

I personally use this one, so you can argue it's well tested. You need to set the following environment variables and the rest works itself out. It also supports URL signing so you get the best speed and results with this one!

```bash
STORAGE_PROVIDER=aws
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_BUCKET_NAME=dflimg
AWS_ROOT=files
AWS_REGION=eu-west-1
```

#### local filesystem

Save to the local filesystem. Please make sure the folder you give it already exists and is writable!

```bash
STORAGE_PROVIDER=lfs
LFS_FOLDER=/Users/duffleman/Downloads/dflimg
LFS_PERMISSIONS=0777
```

## client/cli

A CLI tool that allows you to upload files to the above server!

### Install

Install it into your PATH

`go install ./cmd/dflimg/...`

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

`dflimg c https://avatars1.githubusercontent.com/u/1222340?s=460&v=4`

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
