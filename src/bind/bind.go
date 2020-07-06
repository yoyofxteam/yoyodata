package bind

import (
	"strconv"
	"strings"
)

/**
参数化占位符
*/
const (
	UNKNOWN = iota
	//?
	QUESTION
	//$
	DOLLAR

	NAMED
	//@
	AT
)

/**
获取当前数据库的参数化占位符
*/
func BindType(driverName string) int {
	switch driverName {
	case "mysql":
		return QUESTION

	}
	return UNKNOWN
}

func Rebind(bindType int, query string) string {
	switch bindType {
	case QUESTION:
		return query
	}
	rqb := make([]byte, 0, len(query)+10)

	var i, j int

	for i = strings.Index(query, "?"); i != -1; i = strings.Index(query, "?") {
		rqb = append(rqb, query[:i]...)
		switch bindType {
		case DOLLAR:
			rqb = append(rqb, '$')
		case NAMED:
			rqb = append(rqb, ':', 'a', 'r', 'g')
		case AT:
			rqb = append(rqb, '@', 'p')
		}
		j++
		rqb=strconv.AppendInt(rqb,int64(j),10)
		query=query[i+1:]
	}
	return string(append(rqb, query...))
}

func rebindBuff(bindType int,query stirng)
