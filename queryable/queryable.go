package queryable

import (
	"database/sql"
	"github.com/yoyofxteam/yoyodata/cache"
	"github.com/yoyofxteam/yoyodata/reflectx"
	"reflect"
	"sort"
	"strings"
)

type Queryable struct {
	DB    DbInfo
	Model interface{}
}

/**
执行不带参数化的SQL查询
*/
func (q *Queryable) Query(sql string, res interface{}) {
	db, err := q.DB.CreateNewDbConn()
	if err != nil {
		panic(err)
	}
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	//获取返回值的原始数据类型
	resElem := reflect.ValueOf(res).Elem()
	if resElem.Kind() != reflect.Slice {
		panic("value must be slice")
	}
	//获取对象完全限定名称和元数据
	modelName := reflectx.GetTypeName(q.Model)
	typeInfo := getTypeInfo(modelName, q.Model)
	//获取数据库字段和类型字段的对应关系键值对
	columnFieldSlice := contrastColumnField(rows, typeInfo)
	//创建用于接受数据库返回值的字段变量对象
	scanFieldArray := createScanFieldArray(columnFieldSlice)
	resEleArray := make([]reflect.Value, 0)
	//数据装配
	for rows.Next() {
		//创建对象
		dataModel := reflect.New(reflect.ValueOf(q.Model).Type()).Interface()
		//接受数据库返回值
		rows.Scan(scanFieldArray...)
		//为对象赋值
		setValue(dataModel, scanFieldArray, columnFieldSlice)
		resEleArray = append(resEleArray, reflect.ValueOf(dataModel).Elem())
	}
	//利用反射动态拼接切片
	val := reflect.Append(resElem, resEleArray...)
	resElem.Set(val)
}

/**
数据库字段和类型字段键值对
*/
type ColumnFieldKeyValue struct {
	Index      int
	ColumnName string
	FieldInfo  cache.FieldInfo
}

/**
把数据库返回的值赋值到实体字段上
*/
func setValue(model interface{}, data []interface{}, columnFieldSlice []ColumnFieldKeyValue) {
	modelVal := reflect.ValueOf(model).Elem()
	for i, cf := range columnFieldSlice {
		modelVal.Field(cf.FieldInfo.Index).Set(reflect.ValueOf(data[i]).Elem())
	}
}

/**
创建接受数据库返回值的字段
*/
func createScanFieldArray(columnFieldSlice []ColumnFieldKeyValue) []interface{} {
	var res []interface{}
	for _, data := range columnFieldSlice {
		res = append(res, reflect.New(data.FieldInfo.FieldValue.Type()).Interface())
	}
	return res
}

/**
获取SQL返回的字段名和实际数据类型字段的对比
*/
func contrastColumnField(rows *sql.Rows, typeInfo cache.TypeInfo) []ColumnFieldKeyValue {
	var columnFieldSlice []ColumnFieldKeyValue
	columns, _ := rows.Columns()
	for _, field := range typeInfo.FieldInfo {
		for i, column := range columns {
			if strings.ToUpper(column) == strings.ToUpper(field.FieldName) {
				columnFieldSlice = append(columnFieldSlice, ColumnFieldKeyValue{ColumnName: column, Index: i, FieldInfo: field})
			}
		}
	}

	sort.SliceStable(columnFieldSlice, func(i, j int) bool {
		return columnFieldSlice[i].Index < columnFieldSlice[j].Index
	})
	return columnFieldSlice
}

/**
校验数据类型是否一致
*/
func agreementType(arr *[]interface{}, ele interface{}) {
	//对比元数据是否一致
	isAgreement := reflectx.CompareArrayType(arr, ele)
	if !isAgreement {
		panic("传入的数组类型和查询的类型必须一致")
	}
}

/**
获取元数据
*/
func getTypeInfo(key string, model interface{}) cache.TypeInfo {
	typeInfo, ok := cache.TypeCache.GetTypeInfoCache(key)
	if !ok {
		typeInfo = reflectx.GetTypeInfo(model)
	}
	return typeInfo
}
