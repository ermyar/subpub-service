# SubPub-service
Publisher-Subcriber service.
In first part was implemented `SubPub`package. This is simple bus that works on the principle `Publisher-Subscriber`. Testing is also enable.

Then was implemented subscription service based on this bus. This is gRPC-service, `proto`-file located in `api/subpub`.


## Installing dependencies 

To generate actual files was used `protoc` utility.
Usage described in Makefile.

To download required packages on Linux (Ubuntu) run follow
```
  sudo apt update
  sudo apt install protoc-gen-go
  sudo apt install protoc-gen-go-grcp
```

## Generate and build
```
  make generate              # generate
  go build ./cmd/server/...  # build server
  go build ./cmd/client/...  # build client
```

### Test/Run

You can test main part(subpub bus) by running
```
  go test -v -race ./cmd/subpub/...
```

Server and client supports flags:

- config - set path to config json path. Default set to `configs/config.json` 


Start service by running
```
  go run cmd/server/main.go  # run server
  go run cmd/client/main.go  # run client
```
