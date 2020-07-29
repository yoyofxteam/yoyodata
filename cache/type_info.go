package cache

import (
	"reflect"
	"sync"
)

var TypeCache TypeInfoCache

type TypeInfoCache struct {
	sync.RWMutex
	Items map[string]TypeInfo
}

//缓存初始化
func NewTypeInfoCache() {

	TypeCache = TypeInfoCache{
		Items: make(map[string]TypeInfo),
	}
}

//获取缓存
func (c *TypeInfoCache) GetTypeInfoCache(key string) (TypeInfo, bool) {
	c.RLock()
	defer c.RUnlock()
	value, ok := c.Items[key]
	if ok {
		return value, ok
	}
	return  value, false
}

//添加缓存
func (c *TypeInfoCache) SetTypeInfoCache(key string, typeInfo TypeInfo) {
	c.RLock()
	defer c.RUnlock()
	c.Items[key] = typeInfo
}

//类型缓存
type TypeInfo struct {
	//类型名称
	TypeName string
	//类型下的字段
	FieldInfo []FieldInfo
}

//字段缓存
type FieldInfo struct {
	//字段索引值
	Index      int
	//字段名称
	FieldName  string
	FieldValue reflect.Value
	FieldType  reflect.StructField
}
