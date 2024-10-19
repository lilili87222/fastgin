import type { TLogin } from "@/types/system/user";
import { requestApi } from "../type";

// 登录
export function login(data: TLogin) {
  return requestApi<TLogin>("/api/public/login", "post", data);
}

// 刷新令牌
export function refreshToken() {
  return requestApi<any>("/api/public/refreshToken", "post");
}

// 登出
export function logout() {
  return requestApi<any>("/api/auth/user/logout", "post");
}
