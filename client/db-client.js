class DBClient {
    constructor(baseURL = 'http://localhost:8080') {
        this.baseURL = baseURL;
    }

    async createTable(name, columns) {
        try {
            // Check if table exists first
            try {
                await this.select(name);
                console.log(`Table ${name} already exists, skipping creation`);
                return;
            } catch (e) {
                // Table doesn't exist, continue with creation
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
            throw new Error(`Failed to create table: ${error.message}`);
        }
    }

    async insert(tableName, row) {
        // Convert all numbers to float64 (JSON number type)
        const convertedRow = {};
        for (const [key, value] of Object.entries(row)) {
            if (typeof value === 'number') {
                // All numbers in JavaScript are float64
                // Just send them as-is, the server will validate
                convertedRow[key] = value;
            } else {
                convertedRow[key] = value;
            }
        }

        const response = await fetch(`${this.baseURL}/tables/${tableName}/rows`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(convertedRow)
        });

        if (!response.ok) {
            const error = await response.text();
            throw new Error(error);
        }
    }

    async select(tableName, conditions = {}) {
        const params = new URLSearchParams({ where: JSON.stringify(conditions) });
        const response = await fetch(`${this.baseURL}/tables/${tableName}/rows?${params}`);
        
        if (!response.ok) {
            throw new Error(await response.text());
        }
        
        return response.json();
    }
}

module.exports = DBClient;
