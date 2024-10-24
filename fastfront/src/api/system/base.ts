import type { TLogin } from "@/types/system/user";
import { requestApi } from "../type";

// 登录
export function login(data: TLogin) {
  return requestApi<TLogin>("/api/public/login", "post", data);
}

// 注册
export function register(data: any) {
  return requestApi<TLogin>("/api/public/register", "post", data);
}

// 刷新令牌
export function refreshToken() {
  return requestApi<any>("/api/public/refreshToken", "post");
}

// 登出
export function logout() {
  return requestApi<any>("/api/auth/user/logout", "post");
}

// 获取验证码
export function getCode(params) {
  return requestApi<any>("/api/public/captcha", "get", params);
}

// 发送验证码
export function sendCode(params) {
  return requestApi<any>("/api/public/verifycode", "get", params);
}
