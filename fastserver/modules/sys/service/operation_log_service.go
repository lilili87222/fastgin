package service

import (
	"fastgin/common/httpz"
	"fastgin/database"
	"fastgin/modules/sys/model"
)

type OperationLogService struct {
}

func NewLogService() *OperationLogService {
	return &OperationLogService{}
}
func (s *OperationLogService) BatchDelete(ids []uint64) error {
	return database.DeleteByIds[model.OperationLog](ids)
}
func (s *OperationLogService) Search(req *httpz.SearchRequest) ([]model.OperationLog, int64, error) {
	return database.SearchTable[model.OperationLog](req)
}

// 保存操作日志到数据库
func (s *OperationLogService) SaveOperationLogChannel(olc <-chan *model.OperationLog) {
	var logs []model.OperationLog
	for log := range olc {
		logs = append(logs, *log)
		if len(logs) > 5 {
			database.DB.Create(logs)
			logs = make([]model.OperationLog, 0)
		}
	}
	if len(logs) > 0 {
		database.DB.Create(logs)
	}
}
