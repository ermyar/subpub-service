# SubPub-service
Publisher-Subcriber service.
In first part was implemented `SubPub`package. This is simple bus that works on the principle `Publisher-Subscriber`. Testing is also enable.

Then was implemented in-memory subscription service based on this bus. This is gRPC-service, describing `proto`-file located in `api/subpub`.


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
  make generate              # generate grpc files
  make build                 # build server&client in bin directory
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
  go run ./cmd/server/...             # run server
  ./bin/server --config [configPath]  # also running server
  go run ./cmd/client/...             # run client
  ./bin/client --config [configPath]  # also runnign client
```

Server supports Graceful Shutdown.

Client accepts 2 options (according to service description):
- `publish` `"subs"` `"msg"`
- `subscribe` `"subs"`
