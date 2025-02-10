"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.DBClient = void 0;
const node_fetch_1 = __importDefault(require("node-fetch"));
class DBClient {
    constructor(baseURL = 'http://localhost:8080') {
        this.baseURL = baseURL;
    }
    async createTable(name, columns) {
        try {
            // Check if table exists
            try {
                await this.select(name);
                console.log(`Table ${name} already exists, skipping creation`);
                return;
            }
            catch (e) {
                // Table doesn't exist, continue
            }
            const response = await (0, node_fetch_1.default)(`${this.baseURL}/tables`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name, columns })
            });
            if (!response.ok) {
                throw new Error(await response.text());
            }
        }
        catch (error) {
            if (error instanceof Error) {
                throw new Error(`Failed to create table: ${error.message}`);
            }
            else {
                throw new Error('Failed to create table: Unknown error');
            }
        }
    }
    async insert(tableName, row) {
        const response = await (0, node_fetch_1.default)(`${this.baseURL}/tables/${tableName}/rows`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(row)
        });
        if (!response.ok) {
            throw new Error(await response.text());
        }
    }
    async select(tableName, conditions = {}) {
        const params = new URLSearchParams({ where: JSON.stringify(conditions) });
        const response = await (0, node_fetch_1.default)(`${this.baseURL}/tables/${tableName}/rows?${params}`, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' }
        });
        if (!response.ok) {
            throw new Error(await response.text());
        }
        return response.json();
    }
}
exports.DBClient = DBClient;
exports.default = DBClient;
