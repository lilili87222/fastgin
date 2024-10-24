//日志列表 数据类型
export type TLogs = {
  user_name: string;
  ip: string;
  path: string;
  status: number;
  start_time: string;
  time_cost: number;
  des: string;
  id: number;
};
//日志列表请求参数
export type TLogsQuery = {
  page_size: number;
  page_num: number;
  user_name?: string;
  ip?: string;
  path?: string;
  status?: number;
};
