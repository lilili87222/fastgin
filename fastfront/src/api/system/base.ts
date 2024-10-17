import request from "@/utils/request";
import type { ApiResponse } from "../type";

export function login(data) {
  return request({
    url: "/api/public/login",
    method: "post",
    data,
  }) as Promise<ApiResponse<any>>;
}

export function refreshToken() {
  return request({
    url: "/api/public/refreshToken",
    method: "post",
  }) as Promise<ApiResponse<any>>;
}

export function logout() {
  return request({
    url: "/api/auth/user/logout",
    method: "post",
  }) as Promise<ApiResponse<any>>;
}
