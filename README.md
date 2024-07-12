# SCA

## Requirements

- Installed SQLite3

For Ubuntu:
```sh
sudo apt install sqlite3
```

## Build and Run

```shell
cd src/cmd/ && go build -o server && ./server
```

## Run in Docker

```shell
docker build --no-cache -t go-sca-app . && docker run -p 8080:8080 go-sca-app -d
```