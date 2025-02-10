export interface Column {
    name: string;
    dataType: 'string' | 'int' | 'bool' | 'float' | 'timestamp';
}
export interface Row {
    [key: string]: any;
}
export declare class DBClient {
    private baseURL;
    constructor(baseURL?: string);
    createTable(name: string, columns: Column[]): Promise<void>;
    insert(tableName: string, row: Row): Promise<void>;
    select(tableName: string, conditions?: Row): Promise<Row[]>;
}
export default DBClient;
