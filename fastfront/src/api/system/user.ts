import request from "@/utils/request";
import type { ApiResponse } from "../type";

// 获取当前登录用户信息
export function getInfo(token: string) {
  return request({
    url: "/api/auth/user/info",
    method: "get",
    params: { token },
  }) as Promise<ApiResponse<any>>;
}

// 获取用户列表
export function getUsers(params) {
  return request({
    url: "/api/auth/user/index",
    method: "get",
    params,
  }) as Promise<ApiResponse<any>>;
}

// 更新用户登录密码
export function changePwd(data) {
  return request({
    url: "/api/auth/user/changePwd",
    method: "put",
    data,
  }) as Promise<ApiResponse<any>>;
}

// 创建用户
export function createUser(data) {
  return request({
    url: "/api/auth/user/index",
    method: "post",
    data,
  }) as Promise<ApiResponse<any>>;
}

// 更新用户
export function updateUserById(id, data) {
  return request({
    url: "/api/auth/user/index/" + id,
    method: "patch",
    data,
  }) as Promise<ApiResponse<any>>;
}

// 批量删除用户
export function batchDeleteUserByIds(data) {
  return request({
    url: "/api/auth/user/index",
    method: "delete",
    data,
  }) as Promise<ApiResponse<any>>;
}
