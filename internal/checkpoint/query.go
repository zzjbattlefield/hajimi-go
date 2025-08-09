package checkpoint

import (
	"bufio"
	"os"

	"github.com/zzjbattlefield/hajimi-go/internal/logger"
)

type Query string

func ReadQueryFile(filePath string) []Query {
	var querys = make([]Query, 0)
	if filePath != "" {
		// 从文件读取查询
		query, err := os.Open(filePath)
		if err != nil {
			logger.Log.Errorf("读取查询文件失败: %v", err)
			return querys
		}
		scanner := bufio.NewScanner(query)
		for scanner.Scan() {
			queryBytes := scanner.Bytes()
			if len(queryBytes) == 0 || queryBytes[0] == '#' {
				continue
			}
			querys = append(querys, Query(queryBytes))
		}
	} else {
		querys = append(querys, Query("AIzaSy in:file"))
	}
	return querys
}
