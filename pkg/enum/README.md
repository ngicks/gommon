# enum

Enum-like struct generator.

```
# go run ./cmd/generate_enum/generate_enum.go --help
Usage of /tmp/go-build581005351/b001/exe/generate_enum:
  -disable-goimports
        disable applying goimports after code generation. forced to be false if outFilename is empty.
  -matcher-returns
        whether matcher returns type [T any]. (default true)
  -o string
        out filename. stdout if empty.
  -panic-on-no-match
        whether panic on non-exhaustive match (default true)
  -pkg-name string
        [required] package name of output
  -ty string
        [required] types to enumerate. comma seperated.
  -type-name string
        output typename (default "Enum")
```

It behaves like Tagged union. It actually is not.

see [./exapmle](./example/). It has go:generate.

```
go run ../cmd/generate_enum/generate_enum.go -o enum.go -pkg-name enum -ty int,string,*os.File -matcher-returns true
```

result is ./example/enum.go

Idea based on https://qiita.com/sxarp/items/cd528a546d1537105b9d.