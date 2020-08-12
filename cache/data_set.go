package cache

type DataSet struct {
	//表名
	TableName string
	//结构体名
	ModelName string
}

var DataSetCache map[string]DataSet


func  AddDataSetCache(model interface{},tableName string){


}