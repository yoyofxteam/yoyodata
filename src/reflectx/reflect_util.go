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

/**
使用函数进行字段名转换
*/
func parseFiledName(filed reflect.StructField, tagName string, mapFunc, tagMapFunc filedMapFunc) (tag, fieldName string) {

	//获取原始字段名
	fieldName = filed.Name
	//如果转换函数不为空进行名称转换
	if mapFunc != nil {
		fieldName = mapFunc(fieldName)
	}
	if tagName == "" {
		return "", fieldName
	}

	if !strings.Contains(string(filed.Tag), tagName+":") {
		return "", fieldName
	}
	tag = filed.Tag.Get(tagName)
	if tagMapFunc != nil {
		tag = tagMapFunc(tag)
	}
	parts := strings.Split(tag, ",")
	fieldName = parts[0]
	return tag, fieldName
}

func parseOptions(tag string) map[string]string {
	parts := strings.Split(tag, ",")
	options := make(map[string]string, len(parts))
	if len(parts) > 1 {
		for _, option := range parts[1:] {
			if strings.Contains(option, "=") {
				kv := strings.Split(option, "=")
				options[kv[0]] = kv[1]
				continue
			}
			options[option] = ""
		}
	}
	return options
}

func apnd(is []int, i int) []int {
	x := make([]int, len(is)+1)
	copy(x, is)
	x[len(x)-1] = i
	return x
}

func getMapping(fieldType reflect.Type, tagName string, mapFunc, tagMapFunc filedMapFunc) *StructMap {

	fildInfoItem := []*FieldInfo{}

	root := &FieldInfo{}
	//初始化根节点
	queue := []TypeQueue{}
	queue = append(queue, TypeQueue{parsePtr(fieldType), root, ""})
	//这种写法类似于goto
QueueLoop:
	for len(queue) != 0 {

		typeQueue := queue[0]
		queue = queue[1:]

		for parent := typeQueue.FieldInfo.Parent; parent != nil; parent = parent.Parent {
			if typeQueue.FieldInfo.Filed.Type == parent.Filed.Type {
				continue QueueLoop
			}
		}

		childrenFiledCount := 0
		//如果当前类型是结构体，获取结构体的字段数量
		if typeQueue.typeInfo.Kind() == reflect.Struct {
			childrenFiledCount = typeQueue.typeInfo.NumField()
		}
		typeQueue.FieldInfo.Children = make([]*FieldInfo, childrenFiledCount)

		for filedIndex := 0; filedIndex < childrenFiledCount; filedIndex++ {
			//根据索引获取所代表结构体上的字段
			field := typeQueue.typeInfo.Field(filedIndex)
			//进行字段名转换
			tag, name := parseFiledName(field, tagName, mapFunc, tagMapFunc)
			if name == "-" {
				continue
			}
			fieldInfo := FieldInfo{
				Filed:   field,
				Name:    name,
				Zero:    reflect.New(field.Type),
				Options: parseOptions(tag),
			}

			if typeQueue.parentPath == "" {
				fieldInfo.Path = field.Name
			} else {
				fieldInfo.Path = typeQueue.parentPath + "." + fieldInfo.Name
			}

			//排除匿名类
			if len(field.PkgPath) != 0 && !field.Anonymous {
				continue
			}

			if field.Anonymous {
				parentPath := typeQueue.parentPath
				if tag != "" {
					parentPath = fieldInfo.Path
				}

				fieldInfo.Embedded = true
				fieldInfo.Index = append(typeQueue.FieldInfo.Index, filedIndex)
				childrenFiledCount := 0
				ft := parsePtr(field.Type)
				if ft.Kind() == reflect.Struct {
					childrenFiledCount = ft.NumField()
				}
				fieldInfo.Children = make([]*FieldInfo, childrenFiledCount)
				queue = append(queue, TypeQueue{parsePtr(field.Type), &fieldInfo, parentPath})
			} else if fieldInfo.Zero.Kind() == reflect.Struct ||
				(fieldInfo.Zero.Kind() == reflect.Ptr && fieldInfo.Zero.Type().Elem().Kind() == reflect.Struct) {
				queue = append(queue, TypeQueue{parsePtr(field.Type), &fieldInfo, fieldInfo.Path})
			}
			fieldInfo.Index = apnd(typeQueue.FieldInfo.Index, filedIndex)
			fieldInfo.Parent = typeQueue.FieldInfo
			typeQueue.FieldInfo.Children[filedIndex] = &fieldInfo
			fildInfoItem = append(fildInfoItem, &fieldInfo)

		}

	}
	fields := &StructMap{Index: fildInfoItem, Tree: root, Paths: map[string]*FieldInfo{}, Names: map[string]*FieldInfo{}}
	for _, field := range fields.Index {
		fields.Paths[field.Path] = field
		if field.Name != "" && !field.Embedded {
			fields.Names[field.Path] = field
		}
	}
	return fields
}
