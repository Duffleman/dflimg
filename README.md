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

#### `POST /upload`

Takes a file in the form of multipart/form-data, returns  a short URL that links to the file. You can set the "Accept" header to modify the response.

##### Request

```bash
curl -X POST -H "Authorization: test" -F file=@duffleman.png https://dfl.mn/upload
```

##### Response

```json
{
    "file_id": "file_000000BdAf7MWsYZ6r5wc18cV2sAS",
    "hash": "q3A",
    "url": "https://dfl.mn/q3A"
}
```

If the "Accept" header is set to "text/plain":

`https://dfl.mn/q3A`

#### `GET /{hash}`

Links to the image or file. Serves the content directly!

#### `GET /:{label}`

Links to the image or file through one of it's labels. Serves the content directly!

## client/cli

A CLI tool that allows you to upload files to the above server! Ru

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

Upload a file with some labels, you can give it a CSV for labels

`dflimg u -l test,srs {file}`

It will attempt to automatically put the URL in your clipboard too!
