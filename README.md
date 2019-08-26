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

Takes a file in the form of multipart/form-data, returns  a short URL that links to the file. You can set the "Accept" header to modify the response.

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

#### `POST /tag_resource`

Tags a resource with a label. It requires `tags` which is a CSV set of tags, and `url` which is either a full URL or just the hash of a resource.

See [`POST /upload_file`](https://github.com/Duffleman/dflimg-go#post-upload_file) for the expected response.

#### `POST /shorten_url`

Shorten a URL. It requires `url` which is the URL to shorten. You can apply `nsfw` and `shortcuts` here too.

See [`POST /upload_file`](https://github.com/Duffleman/dflimg-go#post-upload_file) for the expected response.

#### `GET /{hash}`

Links to the resource. Serves the content directly!

#### `GET /:{shortcut}`

Links to the resource through one of it's shortcuts. Serves the content directly!

#### `GET /list_labels`

Returns a list of usable labels

##### Response

```json
[
    {
        "id": "label_000000Bjb0S6DSIaTiW8hSaAo6OOy",
        "name": "education"
    }
]
```

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

### Shorten URL

Shorten a URL

`dflimg s {url}`

See other params above.

### Tag a resource

`dflimg t {url} {labels}`

Where labels is a CSV of labels to apply. The labels must exist on the server.

### Copy a URL

When given a long URL leading to an image, it'll attempt to download the file and reupload it to the dflimg server.

`dflimg c {url}`

`-n` for NSFW works here, along with `-s` for shortcuts.
