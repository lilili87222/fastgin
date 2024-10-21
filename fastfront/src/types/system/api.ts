//接口管理 表格数据
export interface TApiTableData {
  id: number;
  des: string;
  category: string;
  children: Child[];
}

//接口管理 查询参数
export interface TApiQuery {
  method?: string;
  creator?: string;
  category?: string;
  path?: string;
  page_num: number;
  page_size: number;
}

export interface TApiFormData {
  id?: number;
  method: string;
  path: string;
  category: string;
  des: string;
}

export interface Child {
  id: number;
  created_at: Date;
  updated_at: Date;
  deleted_at: null;
  method: Method;
  path: string;
  category: string;
  des: string;
  creator: string;
}

export enum Method {
  Delete = "DELETE",
  Get = "GET",
  Patch = "PATCH",
  Post = "POST",
  Put = "PUT",
}
