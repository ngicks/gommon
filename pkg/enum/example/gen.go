package enum

//go:generate go run ../cmd/generate_enum/generate_enum.go -o enum.go -pkg-name enum -ty int,string,*os.File -matcher-returns=true -panic-on-no-match=true -disable-goimports=false
