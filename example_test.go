package structfields_test

import (
	"fmt"

	"github.com/berquerant/structfields"
)

func ExampleStruct() {
	type Example struct {
		A int `json:"ja" yaml:"ya"`
	}
	s, err := structfields.New(Example{})
	if err != nil {
		panic(err)
	}
	fmt.Println(s.Name)
	fmt.Println(len(s.Fields))
	f := s.Fields[0]
	fmt.Println(f.Name)
	for _, k := range f.TagKeys {
		v, ok := f.Tag.Lookup(k)
		if !ok {
			panic(k)
		}
		fmt.Printf("%s=%s\n", k, v)
	}
	// Output:
	// Example
	// 1
	// A
	// json=ja
	// yaml=ya
}
