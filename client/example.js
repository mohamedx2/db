const { DBClient } = require('godb-hamroun');

async function main() {
    const db = new DBClient();

    try {
        console.log('Creating table...');
        await db.createTable('users', [
            { name: 'id', dataType: 'int' },
            { name: 'name', dataType: 'string' },
            { name: 'active', dataType: 'bool' }
        ]);

        console.log('Inserting data...');
        await db.insert('users', {
            id: 1,
            name: 'John Doe',
            active: true
        });

        // Add update example
        console.log('Updating data...');
        const updated = await db.update('users', 
            { name: 'John Doe' },  // where clause
            { active: false }      // updates
        );
        console.log(`Updated ${updated} rows`);

        console.log('Querying data...');
        const results = await db.select('users', { active: true });
        console.log('Query results:', results);

    } catch (error) {
        console.error('Error:', error.message);
        process.exit(1);
    }
}

main().catch(console.error);
