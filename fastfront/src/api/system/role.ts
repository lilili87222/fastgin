import request from "@/utils/request";
import type { ApiResponse } from "../type";

// 获取角色列表
export function getRoles(params?: any) {
  return request({
    url: "/api/auth/role/index",
    method: "get",
    params,
  }) as Promise<ApiResponse<any>>;
}

// 创建角色
export function createRole(data) {
  return request({
    url: "/api/auth/role/index",
    method: "post",
    data,
  }) as Promise<ApiResponse<any>>;
}

// 更新角色
export function updateRoleById(roleId, data) {
  return request({
    url: "/api/auth/role/index/" + roleId,
    method: "patch",
    data,
  }) as Promise<ApiResponse<any>>;
}

// 获取角色的权限菜单
export function getRoleMenusById(roleId) {
  return request({
    url: "/api/auth/role/menus/" + roleId,
    method: "get",
  }) as Promise<ApiResponse<any>>;
}

// 更新角色的权限菜单
export function updateRoleMenusById(roleId, data) {
  return request({
    url: "/api/auth/role/menus/" + roleId,
    method: "patch",
    data,
  }) as Promise<ApiResponse<any>>;
}

// 获取角色的权限接口
export function getRoleApisById(roleId) {
  return request({
    url: "/api/auth/role/apis/" + roleId,
    method: "get",
  }) as Promise<ApiResponse<any>>;
}

// 更新角色的权限接口
export function updateRoleApisById(roleId, data) {
  return request({
    url: "/api/auth/role/apis/" + roleId,
    method: "patch",
    data,
  }) as Promise<ApiResponse<any>>;
}

// 批量删除角色
export function batchDeleteRoleByIds(data) {
  return request({
    url: "/api/auth/role/index",
    method: "delete",
    data,
  }) as Promise<ApiResponse<any>>;
}
