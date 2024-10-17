import request from "@/utils/request";
import type { ApiResponse } from "../type";

// 获取菜单树
export function getMenuTree() {
  return request({
    url: "/api/auth/menu/tree",
    method: "get",
  }) as Promise<ApiResponse<any>>;
}

// 获取菜单列表
export function getMenus() {
  return request({
    url: "/api/auth/menu/index",
    method: "get",
  }) as Promise<ApiResponse<any>>;
}

// 创建菜单
export function createMenu(data) {
  return request({
    url: "/api/auth/menu/index",
    method: "post",
    data,
  }) as Promise<ApiResponse<any>>;
}

// 更新菜单
export function updateMenuById(Id, data) {
  return request({
    url: "/api/auth/menu/index/" + Id,
    method: "patch",
    data,
  }) as Promise<ApiResponse<any>>;
}

// 批量删除菜单
export function batchDeleteMenuByIds(data) {
  return request({
    url: "/api/auth/menu/index",
    method: "delete",
    data,
  }) as Promise<ApiResponse<any>>;
}

// 获取用户的可访问菜单列表
export function getUserMenusByUserId(Id) {
  return request({
    url: "/api/auth/menu/user/" + Id,
    method: "get",
  }) as Promise<ApiResponse<any>>;
}

// 获取用户的可访问菜单树
export function getUserMenuTreeByUserId(Id) {
  return request({
    url: "/api/auth/menu/user_tree/" + Id,
    method: "get",
  }) as Promise<ApiResponse<any>>;
}
