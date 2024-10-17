import request from "@/utils/request";
import type { ApiResponse } from "../type";

// 获取操作日志列表
export function getSystemInfo() {
  return request({
    url: "/api/auth/system/info",
    method: "get",
  }) as Promise<ApiResponse<any>>;
}

export function restarServer() {
  return request({
    url: "/api/auth/system/restar",
    method: "get",
  }) as Promise<ApiResponse<any>>;
}
export function stopServer() {
  return request({
    url: "/api/auth/system/stop",
    method: "get",
  }) as Promise<ApiResponse<any>>;
}
