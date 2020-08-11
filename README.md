## Simple app for remote install CIA files

| argument | desctiption | example |
| --- | --- | --- |
| path | path to yours cia files | ~/Downloads  |
| port | webserver port | 8090 | 

## installation 

```
go get github.com/ice2heart/fbi_dir
```



## How to run

```
fbi_dir -path ~/Downloads -port 8090
```

## How to build

1. preapare embedded data

```
packr2
```

2. Build 

```
go build
```

## Run gorealeaser for test

```
docker run --rm --privileged \
  -v $PWD:/go/src/github.com/ice2heart/fbi_dir \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -w /go/src/github.com/ice2heart/fbi_dir \
  goreleaser/goreleaser --snapshot --skip-publish --rm-dist
```