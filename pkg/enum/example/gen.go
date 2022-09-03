package enum

//go:generate go run ../cmd/generate_enum/generate_enum.go -o enum_a.go -pkg-name enum -type-name EnumA -ty int,string,*os.File -matcher-returns=true -panic-on-no-match=false -disable-goimports=false
//go:generate go run ../cmd/generate_enum/generate_enum.go -o enum_b.go -pkg-name enum -type-name EnumB -ty int,string,*os.File -matcher-returns=false -panic-on-no-match=true -disable-goimports=false
