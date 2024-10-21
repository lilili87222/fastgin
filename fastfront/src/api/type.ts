import request from "@/utils/request";

export interface ApiResponse<T> {
  data: T | null; // 数据可以是泛型 T 或 null
  message: string;
  code:number
}

//封装请求 返回指定数据结构
export const requestApi = <T>(
  url: string,
  method: "get" | "post" | "patch" | "delete" | "put",
  data?: T
) => {
  return request({
    url,
    method,
    ...(method === "get" ? { params: data } : { data }),
  }) as Promise<ApiResponse<any>>;
};
