//日志列表 数据类型
export type TLogs = {
  UserName: string;
  Ip: string;
  Path: string;
  Status: number;
  StartTime: string;
  TimeCost: number;
  Desc: string;
  Id: number;
};
//日志列表请求参数
export type TLogsQuery = {
  PageNum: number;
  PageSize: number;
  user_name?: string;
  ip?: string;
  path?: string;
  status?: number;
};
