import fetch from 'node-fetch';

export interface Column {
  name: string;
  dataType: 'string' | 'int' | 'bool' | 'float' | 'timestamp';
}

export interface Row {
  [key: string]: any;
}

export class DBClient {
  private baseURL: string;

  constructor(baseURL: string = 'http://localhost:8080') {
    this.baseURL = baseURL;
  }

  async createTable(name: string, columns: Column[]): Promise<void> {
    try {
      // Check if table exists
      try {
        await this.select(name);
        console.log(`Table ${name} already exists, skipping creation`);
        return;
      } catch (e) {
        // Table doesn't exist, continue
      }

      const response = await fetch(`${this.baseURL}/tables`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, columns })
      });

      if (!response.ok) {
        throw new Error(await response.text());
      }
    } catch (error) {
      if (error instanceof Error) {
        throw new Error(`Failed to create table: ${error.message}`);
      } else {
        throw new Error('Failed to create table: Unknown error');
      }
    }
  }

  async insert(tableName: string, row: Row): Promise<void> {
    const response = await fetch(`${this.baseURL}/tables/${tableName}/rows`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(row)
    });

    if (!response.ok) {
      throw new Error(await response.text());
    }
  }

  async select(tableName: string, conditions: Row = {}): Promise<Row[]> {
    const params = new URLSearchParams({ where: JSON.stringify(conditions) });
    const response = await fetch(`${this.baseURL}/tables/${tableName}/rows?${params}`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' }
    });

    if (!response.ok) {
      throw new Error(await response.text());
    }

    return response.json();
  }

  async update(tableName: string, where: Row, updates: Row): Promise<number> {
    const response = await fetch(`${this.baseURL}/tables/${tableName}/rows`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ where, updates })
    });

    if (!response.ok) {
      throw new Error(await response.text());
    }

    const result = await response.json();
    return result.updated;
  }
}

export default DBClient;
