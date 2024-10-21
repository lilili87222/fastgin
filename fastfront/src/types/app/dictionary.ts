export interface TDictionaryTable {
  ID: number;
  CreatedAt: Date;
  UpdatedAt: Date;
  DeletedAt: null;
  Key: string;
  Value: string;
  Desc: string;
}
export interface TDictionaryQuery {
  page_num: number;
  page_size: number;
  Key?: string;
  Value?: string;
  Desc?: string;
}
export interface TDictionaryForm {
  ID?: number;
  Key: string;
  Value: string;
  Desc: string;
}
