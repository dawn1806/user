# User Service

This is the User service

Generated with

```
micro new user
```

## Usage

Generate the proto code

```
make proto
```

Run the service

```
micro run .
```

打包成二进制文件

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o user *.go
```