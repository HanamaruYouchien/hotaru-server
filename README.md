# Hotaru Server

## Build Releases

```sh
go build -ldflags "-X 'main.Version=$(cat VERSION)'" -o ./build/hotaru ./cmd/hotaru
```

## Run

### first time

```sh
hotaru -config ./config/hotaru-server.yml -init
```

### after

```sh
hotaru -config ./config/hotaru-server.yml
```