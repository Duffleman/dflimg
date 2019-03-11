# dflimg

Quick file sharing platform. Accepts images, and files.

This is built from scratch in Go, so you need a public facing server to handle requests and the server will need to run Go with postgres.

Requires an AWS account, it uses S3 as a host for the uploaded files.

When you run this, the shorter your domain is, the better.

Inspired by starbs/yeh

## Endpoints

### `upload_file`

Takes a file in the form of multipart/form-data, stores it in S3, keeps a local cached copy for quick retrieval for the next 10 minutes (configurable). Returns  a short URL that links to the file.

### `/{hash}`

Links to the image or file.
