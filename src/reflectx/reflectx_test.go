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
	mapper:=GetNewMapper("测试", func(s string) string{return "666"+s})
	fmt.Println(reflect.TypeOf(mapper).Kind())

	v:=mapper.GetFieldByName(fv,"A")
	fmt.Println(v.Kind())
	v=mapper.GetFieldByName(fv,"A")
	fmt.Println(v.Kind())
	v=mapper.GetFieldByName(fv,"A")

	fmt.Println(v.Kind())


}

func TestBasicEmbedded(t *testing.T){
	type  Foo struct {
		A int
	}

	type Bar struct {
		Foo Foo
		B int
		C int
	}

	type Baz struct {
		A int
		Bar  Bar
	}
	m:=GetNewMapper("DB", func(s string) string {
		return s
	})
	
	z:=Baz{}
	z.A=1
	z.Bar.C=3
	z.Bar.B=4
	z.Bar.Foo.A=9
	
	zv:=reflect.ValueOf(z)
	fields:=m.TypeMap(reflect.TypeOf(z))

	if len(fields.Index)!=5 {

	}

	v:=m.GetFieldByName(zv,"A")

	v=m.GetFieldByName(zv,"Bar.B")

	v=m.GetFieldByName(zv,"Bar.C")

	v=m.GetFieldByName(zv,"Bar.A")

	fi:=fields.GetFiledByPath("Bar.C")

	fmt.Println(v,fi)


}
