import type { TUserForm, TUserQuery } from "@/types/system/user";
import { requestApi } from "../type";

// 获取当前登录用户信息
export function getInfo(token: string) {
  return requestApi<{ token: string }>("/api/auth/user/info", "get", { token });
}

// 获取用户列表
export function getUsers(params: TUserQuery) {
  return requestApi<TUserQuery>("/api/auth/user/index", "get", params);
}

// 更新用户登录密码
export function changePwd(data: { oldPassword: string; newPassword: string }) {
  return requestApi<{ oldPassword: string; newPassword: string }>(
    "/api/auth/user/changePwd",
    "put",
    data
  );
}

// 创建用户
export function createUser(data: TUserForm) {
  return requestApi<TUserForm>("/api/auth/user/index", "post", data);
}

// 更新用户
export function updateUserById(id: number, data: TUserForm) {
  return requestApi<TUserForm>("/api/auth/user/index/" + id, "patch", data);
}

// 批量删除用户
export function batchDeleteUserByIds(data: any) {
  return requestApi<any>("/api/auth/user/index", "delete", data);
}
