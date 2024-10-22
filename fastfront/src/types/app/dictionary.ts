export interface TDictionaryTable {
    id: number;
    key: string;
    value: string;
    des: string;
}
export interface TDictionaryQuery {
    page_num: number;
    page_size: number;
    id?: number;
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
