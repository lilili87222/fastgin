//接口管理 表格数据
export interface TApiTableData {
  Id: number;
  Desc: string;
  Category: string;
  Children: Child[];
}

//接口管理 查询参数
export interface TApiQuery {
  Method?: string;
  Creator?: string;
  Category?: string;
  Path?: string;
  PageNum: number;
  PageSize: number;
}

export interface TApiFormData {
  Id: number;
  Method: string;
  Path: string;
  Category: string;
  Desc: string;
}

export interface Child {
  Id: number;
  CreatedAt: Date;
  UpdatedAt: Date;
  DeletedAt: null;
  Method: Method;
  Path: string;
  Category: string;
  Desc: string;
  Creator: string;
}

export enum Method {
  Delete = "DELETE",
  Get = "GET",
  Patch = "PATCH",
  Post = "POST",
  Put = "PUT",
}
