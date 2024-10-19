import type { TRoleFormData, TRoleQuery } from "@/types/system/role";
import { requestApi } from "../type";

// 获取角色列表
export function getRoles(params?: TRoleQuery) {
  return requestApi<TRoleQuery>("/api/auth/role/index", "get", params);
}

// 创建角色
export function createRole(data: TRoleFormData) {
  return requestApi<TRoleFormData>("/api/auth/role/index", "post", data);
}

// 更新角色
export function updateRoleById(roleId: number, data: TRoleFormData) {
  return requestApi<TRoleFormData>(
    "/api/auth/role/index/" + roleId,
    "patch",
    data
  );
}

// 获取角色的权限菜单
export function getRoleMenusById(roleId: number) {
  return requestApi<void>("/api/auth/role/menus/" + roleId, "get");
}

// 更新角色的权限菜单
export function updateRoleMenusById(roleId: number, data: any) {
  return requestApi<any>("/api/auth/role/menus/" + roleId, "patch", data);
}

// 获取角色的权限接口
export function getRoleApisById(roleId: number) {
  return requestApi<void>("/api/auth/role/apis/" + roleId, "get");
}

// 更新角色的权限接口
export function updateRoleApisById(roleId: number, data: any) {
  return requestApi<any>("/api/auth/role/apis/" + roleId, "patch", data);
}

// 批量删除角色
export function batchDeleteRoleByIds(data: any) {
  return requestApi<any>("/api/auth/role/index", "delete", data);
}
