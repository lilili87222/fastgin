// types/apiResponse.ts
export interface ApiResponse<T> {
  Data: T | null; // 数据可以是泛型 T 或 null
  Message: string;
}
