# structfields

Extract metadata from `struct` fields.

``` go
type Example struct {
  A int `json:"ja" yaml:"ya"`
}

s, _ := structfields.New(Example{})
s.Fields[0].TagKeys // ["json", "yaml"]
```
