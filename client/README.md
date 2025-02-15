# GoDB Client

A Node.js client for GoDB - A simple database management system with persistent storage and transaction history.

## Installation

```bash
npm install godb-hamroun
```

## Features
- Create tables with typed columns
- Insert, update, delete, and query data
- Data type validation
- Transaction history
- Persistent storage
- Error handling
- TypeScript support

## Usage

```javascript
const { DBClient } = require('godb-hamroun');
// or with TypeScript
import { DBClient } from 'godb-hamroun';

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

// Update data
const updatedCount = await db.update('users', 
  { name: 'John Doe' },  // where clause
  { active: false }      // updates to apply
);
console.log(`Updated ${updatedCount} rows`);

// Query data
const results = await db.select('users', { active: false });
console.log(results);
```

## API Reference

### `new DBClient(baseURL?: string)`
Creates a new client instance. Default URL is 'http://localhost:8080'.

### `createTable(name: string, columns: Column[]): Promise<void>`
Creates a new table with specified columns.

### `insert(tableName: string, row: Row): Promise<void>`
Inserts a row into the specified table.

### `update(tableName: string, where: Row, updates: Row): Promise<number>`
Updates rows matching the where clause with the specified updates.
Returns the number of rows updated.

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

## Error Handling

```javascript
try {
  await db.insert('users', { /* ... */ });
} catch (error) {
  console.error('Error:', error.message);
}
```

## Example with All Operations

```javascript
const db = new DBClient();

// Create table
await db.createTable('users', [
  { name: 'id', dataType: 'int' },
  { name: 'name', dataType: 'string' },
  { name: 'active', dataType: 'bool' },
  { name: 'score', dataType: 'float' }
]);

// Insert
await db.insert('users', {
  id: 1,
  name: 'John',
  active: true,
  score: 95.5
});

// Update
const updated = await db.update('users',
  { id: 1 },              // where
  { score: 98.5 }         // updates
);

// Query
const activeUsers = await db.select('users', { active: true });
```

## License

MIT

## Author

Mohamed Ali Hamroun
