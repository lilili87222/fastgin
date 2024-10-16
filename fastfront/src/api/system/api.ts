import request from "@/utils/request";
import type { ApiResponse } from "../type";

// 获取接口列表
export function getApis(params) {
  return request({
    url: "/api/auth/api/index",
    method: "get",
    params,
  }) as Promise<ApiResponse<any>>;
}

// 获取接口树(按接口Category字段分类)
export function getApiTree() {
  return request({
    url: "/api/auth/api/tree",
    method: "get",
  }) as Promise<ApiResponse<any>>;
}

// 创建接口
export function createApi(data) {
  return request({
    url: "/api/auth/api/index",
    method: "post",
    data,
  }) as Promise<ApiResponse<any>>;
}

// 更新接口
export function updateApiById(Id, data) {
  return request({
    url: "/api/auth/api/index/" + Id,
    method: "patch",
    data,
  }) as Promise<ApiResponse<any>>;
}

// 批量删除接口
export function batchDeleteApiByIds(data) {
  return request({
    url: "/api/auth/api/index",
    method: "delete",
    data,
  }) as Promise<ApiResponse<any>>;
}
