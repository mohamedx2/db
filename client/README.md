# GoDB Client

A Node.js client for GoDB - A simple database system.

## Installation

```bash
npm install godb-client
```

## Usage

```javascript
const { DBClient } = require('godb-client');
// or
import { DBClient } from 'godb-client';

const db = new DBClient('http://localhost:8080');

// Create a table
await db.createTable('users', [
  { name: 'id', dataType: 'int' },
  { name: 'name', dataType: 'string' },
  { name: 'active', dataType: 'bool' }
]);

// Insert data
await db.insert('users', {
  id: 1,
  name: 'John Doe',
  active: true
});

// Query data
const results = await db.select('users', { active: true });
console.log(results);
```

## API

### `new DBClient(baseURL?: string)`
Creates a new client instance. Default URL is 'http://localhost:8080'.

### `createTable(name: string, columns: Column[]): Promise<void>`
Creates a new table with specified columns.

### `insert(tableName: string, row: Row): Promise<void>`
Inserts a row into the specified table.

### `select(tableName: string, conditions?: Row): Promise<Row[]>`
Queries data from the specified table with optional conditions.

## Types

```typescript
interface Column {
  name: string;
  dataType: 'string' | 'int' | 'bool' | 'float' | 'timestamp';
}

interface Row {
  [key: string]: any;
}
```

## License

MIT
