import request from "@/utils/request";
import type { ApiResponse } from "../type";

// 获取操作日志列表
export function getOperationLogs(params) {
  return request({
    url: "/api/auth/log/index",
    method: "get",
    params,
  }) as Promise<ApiResponse<any>>;
}

// 批量删除操作日志
export function batchDeleteOperationLogByIds(data) {
  return request({
    url: "/api/auth/log/index",
    method: "delete",
    data,
  }) as Promise<ApiResponse<any>>;
}
