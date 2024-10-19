import { requestApi } from "../type";
import type {
  TDictionaryFormData,
  TDictionaryQuery,
} from "@/types/app/dictionary";

// 获取字典列表
export function getDictionary(params: TDictionaryQuery) {
  return requestApi<TDictionaryQuery>(
    "/api/auth/dictionary/index",
    "get",
    params
  );
}

// 获取字典列表详情
export function getDictionaryDetail(Id: number) {
  return requestApi<any>(`/api/auth/dictionary/index/${Id}`, "get");
}

// 新增字典列表
export function createDictionary(data: TDictionaryFormData) {
  return requestApi<TDictionaryFormData>(
    "/api/auth/dictionary/index",
    "post",
    data
  );
}

// 更新字典
export function updateDictionaryById(Id: number, data: TDictionaryFormData) {
  return requestApi<TDictionaryFormData>(
    `/api/auth/dictionary/index/${Id}`,
    "patch",
    data
  );
}

// 批量删除字典
export function batchDeleteDictionaryByIds(data: any) {
  return requestApi<any>("/api/auth/dictionary/index", "delete", data);
}

// 单个删除字典
export function batchDeleteDictionaryById(Id: number) {
  return requestApi<any>(`/api/auth/dictionary/index/${Id}`, "delete");
}
