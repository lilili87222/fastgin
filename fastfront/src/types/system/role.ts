//角色管理 表格数据类型
export interface TRoleTableData {
  id: number;
  created_at: Date;
  updated_at: Date;
  deleted_at: null;
  name: string;
  keyword: string;
  des: string;
  status: number;
  sort: number;
  creator: string;
  users: null;
  menus: null;
}

//角色管理 添加/编辑表单数据类型
export interface TRoleFormData {
  id?: number;
  name: string;
  keyword: string;
  des: string;
  status: number;
  sort: number;
}

//角色管理 搜索表单数据类型
export interface TRoleQuery {
  page_num: number;
  page_size: number;
  name?: string;
  keyword?: string;
  status?: number;
}
