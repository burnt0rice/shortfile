# shortFile

Service to upload Files and download it with an unique URL.
Under the hood the file will be compressed and saved.

## Run frontend and backend

```shell
go run main.go -port=3000
```

## Build frontend

```shell
cd frontend
npm install
npm run build
```

## Requests

### Upload new file

#### Request

```shell
curl -F "file=@/path/to/file.jpg" http://localhost:3000/api/u
```

#### Response

- message `string`
- downloadUrl `string`

### Get file info

#### Request

```shell
curl http://localhost:3000/api/i/:id
```

#### Response

- message `string`
- file `object`
- file.name `string`
- file.timestamp `string`
- file.uuid `number`

### Download file

#### Request

```shell
curl http://localhost:3000/api/d/:id
```

#### Response

- File `file`

## Todo

- [ ] Job to delete files after 30 days

## Others

<a target="_blank" href="https://icons8.com/icon/g8mOI88XbBJX/happy-file">Happy File</a> icon by <a target="_blank" href="https://icons8.com">Icons8</a>

[Embedding Vue.js Apps in Go](https://hackandsla.sh/posts/2021-06-18-embed-vuejs-in-go/)
