//用户管理 表格数据类型
export interface TUserTableData {
  Id: number;
  UserName: string;
  Mobile: string;
  Avatar: string;
  NickName: string;
  Introduction: string;
  Status: number;
  Creator: string;
  RoleIds: number[];
}

//用户管理 请求参数类型
export type TUserQuery = {
  PageNum: number;
  PageSize: number;
  user_name?: string;
  nick_name?: string;
  mobile?: string;
  status?: number;
};

//用户管理 新增 编辑 表单数据类型
export type TUserFormData = {
  Id?: number;
  UserName: string;
  Password: string;
  NickName: string;
  Status: number;
  Mobile: string;
  Avatar: string;
  Introduction: string;
  RoleIds: string;
};

export interface TLogin {
  username: string;
  password: string;
}
