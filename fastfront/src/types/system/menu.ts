//菜单管理 表格数据
export interface TMenuTableData {
  Id: number;
  CreatedAt: Date;
  UpdatedAt: Date;
  DeletedAt: null;
  Name: string;
  Title: string;
  Icon: string;
  Path: string;
  Redirect: null | string;
  Component: string;
  Sort: number;
  Status: number;
  Hidden: number;
  NoCache: number;
  AlwaysShow: number;
  Breadcrumb: number;
  ActiveMenu: null;
  ParentId: number;
  Creator: string;
  Children: TMenuTableData[];
  Roles: null;
}

export interface TMenuFormData {
  Id: number;
  Title: string;
  Name: string;
  Sort: number;
  Icon: string;
  Path: string;
  Component: string;
  Redirect: string;
  Status: number;
  Hidden: number;
  NoCache: number;
  ActiveMenu: string;
  ParentId: number;
}

//菜单管理 查询参数
export interface TMenuQuery {
  PageNum: number;
  PageSize: number;
}
