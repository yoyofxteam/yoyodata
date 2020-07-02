/**
反射工具类
*/
package reflectx

import (
	_ "github.com/yoyofxteam/yoyodata/src/model"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

/**
字段属性结构体
*/
type FieldInfo struct {
	Index    []int
	Path     string
	Filed    reflect.StructField
	Zero     reflect.Value
	Name     string
	Options  map[string]string
	Embedded bool
	//子集字段
	Children []*FieldInfo
	//父级字段
	Parent *FieldInfo
}

type StructMap struct {
	Tree  *FieldInfo
	Index []*FieldInfo
	Paths map[string]*FieldInfo
	Names map[string]*FieldInfo
}

/**
字段映射结构体
*/
type Mapper struct {
	cache      map[reflect.Type]*StructMap
	tagName    string
	tagMapFunc func(string) string
	mapFunc    func(string) string
	mutex      sync.Mutex
}

/**
类型队列
*/
type TypeQueue struct {
	typeInfo   reflect.Type
	FieldInfo  *FieldInfo
	parentPath string
}

/**
获取一个新的mapper对象
*/
func getNewMapper(tagName string, mapFunc func(string) string) *Mapper {
	return &Mapper{
		tagName: tagName,
		mapFunc: mapFunc,
		cache:   make(map[reflect.Type]*StructMap),
	}

}

type kinder interface {
	Kind() reflect.Kind
}

/**
判断是否是结构体类型
*/
func mustBeStruct(v kinder, expected reflect.Kind) {
	if k := v.Kind(); k != expected {
		panic(reflect.ValueError{Method: getStackInfo(), Kind: k})
	}
}

/**
获取堆栈信息
*/
func getStackInfo() string {

	//获取函数名，文件，行号
	pc, file, line, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	stockInfo := strings.Builder{}
	stockInfo.WriteString(f.Name())
	stockInfo.WriteString(file)
	stockInfo.WriteString(string(line))
	return stockInfo.String()
}

func (m *Mapper) TypeMap(t reflect.Type) *Mapper {
	m.mutex.Lock()
	mapping, ok := m.cache[t]
	if !ok {
		mapping =
	}

}

/**
字段映射函数委托
*/
type filedMapFunc func(string) string

/**
把指针类型转换成指针代表的原始类型
*/
func parsePtr(t reflect.Type) reflect.Type {

	//如果t类型是指针类型转换成指针代表的原始类型
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

func getMapping(fieldType reflect.Type, tagName string, mapFunc, tagMapFunc filedMapFunc) *StructMap {

	fildInfoItem := []*FieldInfo{}

	root := &FieldInfo{}
	//初始化根节点
	queue := []TypeQueue{}
	queue = append(queue, TypeQueue{parsePtr(fieldType),root,""})
loop:

}
