import { requestApi } from "../type";
import type { TLogsQuery } from "@/types/log/operation-log";

export function getOperationLogs(params: TLogsQuery) {
  return requestApi<TLogsQuery>("/api/auth/log/index", "get", params);
}

// 批量删除操作日志
export function batchDeleteOperationLogByIds(data: any) {
  return requestApi<any>("/api/auth/log/index", "delete", data);
}
