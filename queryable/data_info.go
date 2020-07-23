package queryable

import (
	"database/sql"
)

/**
数据库链接信息
*/
type DbInfo struct {
	DriverName string
	DataSource string
}

/**
获取链接字符串
*/
func (db *DbInfo) CreateNewDbConn() (*sql.DB, error) {
	return sql.Open(db.DriverName, db.DataSource)
}
