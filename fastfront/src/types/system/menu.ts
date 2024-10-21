//菜单管理 表格数据
export interface TMenuTable {
  id: number;
  created_at: Date;
  updated_at: Date;
  deleted_at: null;
  name: string;
  title: string;
  icon: string;
  path: string;
  redirect: null | string;
  component: string;
  sort: number;
  status: number;
  hidden: number;
  no_cache: number;
  always_show: number;
  breadcrumb: number;
  active_menu: null;
  parent_id: number;
  creator: string;
  children: TMenuTable[];
  roles: null;
}

export interface TMenuForm {
  id: number;
  title: string;
  name: string;
  sort: number;
  icon: string;
  path: string;
  component: string;
  redirect: string;
  status: number;
  hidden: number;
  no_cache: number;
  active_menu: string;
  parent_id: number;
}

//菜单管理 查询参数
export interface TMenuQuery {
  page_num: number;
  page_size: number;
}
