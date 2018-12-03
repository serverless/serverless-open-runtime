# sidecartester - A tool for testing language side car implementations

## Build
```
go build
```

## Test the Node JS runtime implementation
```
$ ./sidecartester node ../nodejs/language-runtime handlers.hello
START open-runtime-nodejs
EVENT open-runtime-nodejs: "{\"asdf\": 53423}"
"Hello World!"
2018/12/03 11:01:48 result from sidecar: "Hello World!"
```
