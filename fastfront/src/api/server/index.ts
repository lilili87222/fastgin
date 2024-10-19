import { requestApi } from "../type";

// 获取操作日志列表
export function getSystemInfo() {
  return requestApi<any>("/api/auth/system/info", "get");
}

export function restartServer() {
  return requestApi<any>("/api/auth/system/restart", "get");
}

export function stopServer() {
  return requestApi<any>("/api/auth/system/stop", "get");
}
