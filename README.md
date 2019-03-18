# dflimg

[DockerHub](https://hub.docker.com/r/duffleman/dflimg)

Quick file sharing platform. Accepts images, and files.

This is built from scratch in Go, so you need a public facing server to handle requests and the server will need to run Go with postgres.

Requires an AWS account, it uses S3 as a host for the uploaded files.

When you run this, the shorter your domain is, the better.

Inspired by starbs/yeh

## Env variables to set

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

## Endpoints

### `POST /upload`

Takes a file in the form of multipart/form-data, returns  a short URL that links to the file. You can set the "Accept" header to modify the response.

#### Request

```bash
curl -X POST -H "Authorization: test" -F file=@duffleman.png https://dfl.mn/upload
```

#### Response

```json
{
    "file_id": "file_000000BdAf7MWsYZ6r5wc18cV2sAS",
    "hash": "q3A",
    "url": "https://dfl.mn/q3A"
}
```

If the "Accept" header is set to "text/plain":

`https://dfl.mn/q3A`

### `GET /{hash}`

Links to the image or file. Serves the content directly!
