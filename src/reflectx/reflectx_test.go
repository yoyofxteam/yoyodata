package reflectx

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBasic(t *testing.T){

	type Foo struct {
		A int
		B string
		C float64
	}

	type Zoo struct {
		D Foo
		F string
	}

	f:=Foo{1,"小崔",3.14}
	z:=Zoo{f,"曹"}
	fv:=reflect.ValueOf(z)
	mapper:=getNewMapper("测试", func(s string) string{return "666"+s})
	fmt.Println(reflect.TypeOf(mapper).Kind())

	v:=mapper.FieldByName(fv,"A")
	fmt.Println(v.Kind())
	v=mapper.FieldByName(fv,"A")
	fmt.Println(v.Kind())
	v=mapper.FieldByName(fv,"A")

	fmt.Println(v.Kind())


}
