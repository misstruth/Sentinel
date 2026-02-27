package tools

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlCrudInput struct {
	DSN         string `json:"dsn" jsonschema:"description=The Data Source Name for connecting to the MySQL database, including username, password, host, port, and database name"`
	SQL         string `json:"sql" jsonschema:"description=The SQL query to execute against the MySQL database"`
	OperateType string `json:"operate_type" jsonschema:"description=The type of SQL operation to perform: query, insert, update, or delete"`
}

func NewMysqlCrudTool() tool.InvokableTool {
	t, err := utils.InferOptionableTool(
		"mysql_crud",
		"Execute SQL queries against the MySQL database and return results in JSON format. Use this tool when you need to query, insert, update or delete data from the database. The results will be formatted as JSON for easy parsing.",
		func(ctx context.Context, input *MysqlCrudInput, opts ...tool.Option) (output string, err error) {
			db, err := gorm.Open(mysql.Open(input.DSN), &gorm.Config{})
			if err != nil {
				return fmt.Sprintf("db connect error: %v", err), nil
			}

			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("\n请确定是否执行本sql(y/n): ", input.SQL)
			scanner.Scan()
			fmt.Println()
			if scanner.Text() != "y" {
				return "用户取消执行sql", nil
			}

			if err = db.Exec(input.SQL).Error; err != nil {
				return fmt.Sprintf("exec error: %v", err), nil
			}
			if input.OperateType == "query" {
				var results []interface{}
				if err = db.Raw(input.SQL).Scan(&results).Error; err != nil {
					return fmt.Sprintf("query error: %v", err), nil
				}
				resBytes, _ := json.Marshal(results)
				return string(resBytes), nil
			}
			return "success", nil
		})
	if err != nil {
		return nil
	}
	return t
}
