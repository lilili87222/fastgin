export interface TDictionaryTable {
  ID: number;
  created_at: Date;
  updated_at: Date;
  deleted_at: null;
  key: string;
  value: string;
  des: string;
}
export interface TDictionaryQuery {
  page_num: number;
  page_size: number;
  key?: string;
  value?: string;
  des?: string;
}
export interface TDictionaryForm {
  id?: number;
  key: string;
  value: string;
  des: string;
}
