//角色管理 表格数据类型
export interface TRoleTableData {
  Id: number;
  CreatedAt: Date;
  UpdatedAt: Date;
  DeletedAt: null;
  Name: string;
  Keyword: string;
  Desc: string;
  Status: number;
  Sort: number;
  Creator: string;
  Users: null;
  Menus: null;
}

//角色管理 添加/编辑表单数据类型
export interface TRoleFormData {
  Id?: number;
  Name: string;
  Keyword: string;
  Desc: string;
  Status: number;
  Sort: number;
}

//角色管理 搜索表单数据类型
export interface TRoleQuery {
  PageNum: number;
  PageSize: number;
  Name?: string;
  Keyword?: string;
  Status?: number;
}
