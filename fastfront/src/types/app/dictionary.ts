export interface TDictionaryTableData {
  ID: number;
  CreatedAt: Date;
  UpdatedAt: Date;
  DeletedAt: null;
  Key: string;
  Value: string;
  Desc: string;
}
export interface TDictionaryQuery {
  PageNum: number;
  PageSize: number;
  Key?: string;
  Value?: string;
  Desc?: string;
}
export interface TDictionaryFormData {
  ID?: number;
  Key: string;
  Value: string;
  Desc: string;
}
