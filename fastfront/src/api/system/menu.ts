import type { TMenuForm } from "@/types/system/menu";
import { requestApi } from "../type";

// 获取菜单树
export function getMenuTree() {
  return requestApi<void>("/api/auth/menu/tree", "get");
}

// 获取菜单列表
export function getMenus() {
  return requestApi<void>("/api/auth/menu/index", "get");
}

// 创建菜单
export function createMenu(data: TMenuForm) {
  return requestApi<TMenuForm>("/api/auth/menu/index", "post", data);
}

// 更新菜单
export function updateMenuById(Id: number, data: TMenuForm) {
  return requestApi<TMenuForm>("/api/auth/menu/index/" + Id, "patch", data);
}

// 批量删除菜单
export function batchDeleteMenuByIds(data: any) {
  return requestApi<any>("/api/auth/menu/index", "delete", data);
}

// 获取用户的可访问菜单列表
export function getUserMenusByUserId(Id: number) {
  return requestApi<void>("/api/auth/menu/user/" + Id, "get");
}

// 获取用户的可访问菜单树
export function getUserMenuTreeByUserId(Id: number) {
  return requestApi<void>("/api/auth/menu/user_tree/" + Id, "get");
}
