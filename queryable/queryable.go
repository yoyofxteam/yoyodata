package queryable

import (
	"database/sql"
	"github.com/yoyofxteam/yoyo-reflect"
	"github.com/yoyofxteam/yoyodata/cache"
	"github.com/yoyofxteam/yoyodata/reflectx"
	"reflect"
)

type Queryable struct {
	DB    DbInfo
	Model interface{}
}

func (q *Queryable) Query(sql string, res interface{}) {

	db, err := q.DB.CreateNewDbConn()
	if err != nil {
		panic(err)
	}
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	//校验数据类型是否一致
	//agreementType(res, q.model)
	//获取对象元数据
	modelName := reflectx.GetTypeName(q.Model)
	typeInfo := getTypeInfo(modelName, q.Model)
	//获取数据库字段和类型字段键值对
	columnFieldSlice := contrastColumnField(rows, typeInfo)
	//获取要扫码的字段数组
	scanFieldArray := createScanFieldArray(columnFieldSlice)
	//数据装配
	resPtr := reflect.ValueOf(res).Elem()
	resEleArray := make([]reflect.Value, 0)
	for rows.Next() {
		dataModel := Reflect.CreateInstancePtr(reflect.ValueOf(q.Model).Type())
		rows.Scan(scanFieldArray...)
		resEle := setValue(&dataModel, scanFieldArray, columnFieldSlice)
		resEleArray = append(resEleArray, reflect.ValueOf(resEle))
	}
	val:= reflect.Append(resPtr,resEleArray...)
	resPtr.Set(val)
}

/**
数据库字段和类型字段键值对
*/
type ColumnFieldKeyValue struct {
	ColumnName string
	FieldInfo  cache.FieldInfo
}

func setValue(model *interface{}, data []interface{}, columnFieldSlice []ColumnFieldKeyValue) *interface{} {
	modelVal := reflect.ValueOf(model)
	for i, cf := range columnFieldSlice {
		modelVal.Field(cf.FieldInfo.Index).Set(reflect.ValueOf(data[i]).Elem())
	}
	return model
}

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
		for _, column := range columns {
			if column == field.FieldName {
				columnFieldSlice = append(columnFieldSlice, ColumnFieldKeyValue{ColumnName: column, FieldInfo: field})
			}
		}
	}
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
