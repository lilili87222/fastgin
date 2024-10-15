package service

import (
	"fastgin/sys/dao"
	"fastgin/sys/dto"
	"fastgin/sys/model"
)

type OperationLogService struct {
	logDao *dao.OperationLogDao
}

func NewLogService() *OperationLogService {
	return &OperationLogService{logDao: &dao.OperationLogDao{}}
}

// 获取操作日志列表
func (s *OperationLogService) GetOperationLogs(req *dto.OperationLogListRequest) ([]model.OperationLog, int64, error) {
	return s.logDao.GetOperationLogs(req)
}

// 批量删除操作日志
func (s *OperationLogService) BatchDeleteOperationLogByIds(ids []uint) error {
	return s.logDao.BatchDeleteOperationLogByIds(ids)
}

// 保存操作日志到数据库
//
//	func (s *OperationLogService) SaveOperationLogChannel(olc <-chan *model.OperationLog) {
//		s.logDao.SaveOperationLogChannel(olc)
//	}
//
// 保存操作日志到数据库
func (s *OperationLogService) SaveOperationLogChannel(olc <-chan *model.OperationLog) {
	var logs []model.OperationLog

	for log := range olc {
		logs = append(logs, *log)
		if len(logs) > 5 {
			s.logDao.SaveOperationLogs(logs)
			logs = make([]model.OperationLog, 0)
		}
	}
	if len(logs) > 0 {
		s.logDao.SaveOperationLogs(logs)
	}
}
