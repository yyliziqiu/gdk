package xdb

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
)

func IsNoRowsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func IsMysqlTableNotExistError(err error) bool {
	return strings.Contains(err.Error(), "Error 1146 (42S02)")
}

func IsMysqlDuplicateKeyError(err error) bool {
	err2, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}
	return err2.Number == 1062
}

func IsPgsqlDuplicateKeyError(err error) bool {
	err2, ok := err.(*pq.Error)
	if !ok {
		return false
	}
	return err2.Code == "23505"
}

func JoinIntValue(values []int) string {
	length := len(values)
	if length == 0 {
		return ""
	}
	sb := strings.Builder{}
	for i := 0; i < length-1; i++ {
		sb.WriteString(strconv.Itoa(values[i]))
		sb.WriteString(", ")
	}
	sb.WriteString(strconv.Itoa(values[length-1]))
	return sb.String()
}

func JoinStringValue(values []string) string {
	length := len(values)
	if length == 0 {
		return ""
	}
	return "'" + strings.Join(values, "', '") + "'"
}

func SafeJoinStingValue(values []string) string {
	length := len(values)
	for i := 0; i < length; i++ {
		values[i] = EscapeString(values[i])
	}
	return JoinStringValue(values)
}

func EscapeString(str string) string {
	chars := []rune(str)
	temp := make([]rune, 0, len(chars))
	for _, c := range chars {
		if c == '\\' || c == '"' || c == '\'' {
			temp = append(temp, '\\')
			temp = append(temp, c)
		} else {
			temp = append(temp, c)
		}
	}
	return string(temp)
}
