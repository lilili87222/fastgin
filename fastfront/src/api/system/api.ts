import type { TApiFormData, TApiQuery } from "@/types/system/api";
import { requestApi } from "../type";

// 获取接口列表
export function getApis(params: TApiQuery) {
  return requestApi<TApiQuery>("/api/auth/api/index", "get", params);
}

// 获取接口树(按接口Category字段分类)
export function getApiTree() {
  return requestApi<any>("/api/auth/api/tree", "get");
}

// 创建接口
export function createApi(data: TApiFormData) {
  return requestApi<TApiFormData>("/api/auth/api/index", "post", data);
}

// 更新接口
export function updateApiById(Id: number, data: TApiFormData) {
  return requestApi<TApiFormData>(`/api/auth/api/index/${Id}`, "patch", data);
}

// 批量删除接口
export function batchDeleteApiByIds(data: any) {
  return requestApi<any>("/api/auth/api/index", "delete", data);
}
