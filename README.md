# Hotaru Server

## Build Releases

```sh
go build -ldflags "-X 'main.Version=$(cat VERSION)'" -o ./build/hotaru ./cmd/hotaru
```