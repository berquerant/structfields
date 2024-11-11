package internal

//go:generate go run golang.org/x/tools/cmd/goyacc -o tag_goyacc_generated.go -v tag_goyacc_generated.output tag.y

type Tags struct {
	List []*Tag
}

type Tag struct {
	Key   string
	Value string
}
