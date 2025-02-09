const DBClient = require('./db-client');

async function main() {
    const db = new DBClient();

    try {
        // Create a table (will skip if already exists)
        await db.createTable('users', [
            { name: 'id', dataType: 'int' },
            { name: 'name', dataType: 'string' },
            { name: 'active', dataType: 'bool' }
        ]);

        // Wait a moment to ensure table is created
        await new Promise(resolve => setTimeout(resolve, 100));

        console.log('Inserting data...');
        await db.insert('users', {
            id: 1,
            name: 'John Doe',
            active: true
        });
        console.log('Data inserted successfully');

        console.log('Querying data...');
        const results = await db.select('users', { active: true });
        console.log('Query results:', results);

    } catch (error) {
        console.error('Error:', error.message);
        process.exit(1);
    }
}

main().catch(console.error);
