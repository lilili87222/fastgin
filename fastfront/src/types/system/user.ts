//用户管理 表格数据类型
export interface TUserTableData {
  id: number;
  user_name: string;
  mobile: string;
  avatar: string;
  nick_name: string;
  des: string;
  status: number;
  creator: string;
  role_ids: number[];
}

//用户管理 请求参数类型
export type TUserQuery = {
  page_num: number;
  page_size: number;
  user_name?: string;
  nick_name?: string;
  mobile?: string;
  status?: number;
};

//用户管理 新增 编辑 表单数据类型
export type TUserFormData = {
  id?: number;
  user_name: string;
  password: string;
  nick_name: string;
  status: number;
  mobile: string;
  avatar: string;
  des: string;
  role_ids: string;
};

export interface TLogin {
  username: string;
  password: string;
}
