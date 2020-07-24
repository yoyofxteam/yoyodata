package reflectx

import (
	"github.com/yoyofxteam/yoyodata/cache"
	"reflect"
)

/**
根据传入的结构体获取结构体的元数据缓存起来
*/
func ReflectTypeInfo(model interface{}) cache.TypeInfo {
	modelValue := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)
	pkg := modelType.PkgPath()
	typeName := pkg + modelType.Name()
	if modelValue.Kind() != reflect.Struct {
		panic("model must be struct !")
	}
	var fieldInfoArray []cache.FieldInfo
	for i := 0; i < modelValue.NumField(); i++ {
		fieldValue := modelValue.Field(i)
		if fieldValue.Kind() == reflect.Struct {
			continue
		}
		fieldType := modelType.Field(i)
		fieldName := fieldType.Name
		fieldInfoElement := cache.FieldInfo{
			Index:      i,
			FieldName:  fieldName,
			FieldType:  fieldType,
			FieldValue: fieldValue,
		}
		fieldInfoArray = append(fieldInfoArray, fieldInfoElement)
	}
	typeInfo := cache.TypeInfo{
		TypeName:  typeName,
		FieldInfo: fieldInfoArray,
	}
	return typeInfo
}

/**
获取类型元数据信息
*/
func GetTypeInfo(model interface{}) cache.TypeInfo {
	modelType := reflect.TypeOf(model)
	typeName := modelType.PkgPath() + modelType.Name()
	typeInfo, ok := cache.TypeCache.GetTypeInfoCache(typeName)
	if ok {
		return typeInfo
	}
	typeInfo = ReflectTypeInfo(model)
	cache.TypeCache.SetTypeInfoCache(typeName, typeInfo)
	return typeInfo
}

func GetArrayEleType(model *[]interface{}) string {

	return reflect.TypeOf(model).Elem().Name()
}

func CompareArrayType(arr *[]interface{}, ele interface{}) bool {
	arrType := reflect.TypeOf(arr).Elem()
	eleType := reflect.TypeOf(ele)
	arrTRypeName := arrType.PkgPath() + arrType.Name()
	eleTypeName := eleType.PkgPath() + eleType.Name()
	return arrTRypeName == eleTypeName
}

func GetTypeName(model interface{}) string {
	modelType := reflect.TypeOf(model)
	modelTypeName := modelType.PkgPath() + modelType.Name()
	return modelTypeName
}
