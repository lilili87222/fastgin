package service

import (
	"fastgin/sys/dao"
	"fastgin/sys/model"
	"fastgin/sys/model/request"
)

type LogService struct {
	logDao dao.OperationLogDao
}

func NewLogService() LogService {
	return LogService{logDao: dao.OperationLogDao{}}
}

// 获取操作日志列表
func (s LogService) GetOperationLogs(req *request.OperationLogListRequest) ([]model.OperationLog, int64, error) {
	return s.logDao.GetOperationLogs(req)
}

// 批量删除操作日志
func (s LogService) BatchDeleteOperationLogByIds(ids []uint) error {
	return s.logDao.BatchDeleteOperationLogByIds(ids)
}

// 保存操作日志到数据库
func (s LogService) SaveOperationLogChannel(olc <-chan *model.OperationLog) {
	s.logDao.SaveOperationLogChannel(olc)
}
