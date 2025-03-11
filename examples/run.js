import cql from 'k6/x/cql';

export let options = {
    scenarios: {
        setup: {
            executor: 'shared-iterations',
            vus: 1,
            iterations: 1,
            exec: 'setup',
            startTime: '0s',
        },
        stress_test: {
            executor: 'ramping-vus',
            stages: [
                { duration: '10s', target: 3 },
                { duration: '5s', target: 1 },
                { duration: '1s', target: 0 }
            ],
            exec: 'stressTest',
            startTime: '10s'
        },
    }
};

export function setup() {
    let client = cql;
    let hosts = __ENV.CQL_HOSTS || "127.0.0.1:9042";
    const keyspace = "test_keyspace";

    try {
        client.session({
            hosts: hosts.split(","),
            keyspace: "system",
        });

        console.log("Creating keyspace...");
        client.exec(`
            CREATE KEYSPACE IF NOT EXISTS test_keyspace 
            WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
        `);
    } catch (error) {
        console.error("Failed to create keyspace", error);
        throw error;
    }

    try {
        client.session({
            hosts: hosts.split(","),
            keyspace: keyspace,
        });

        console.log("Creating table...");
        client.exec(`
            CREATE TABLE IF NOT EXISTS users (
                id UUID PRIMARY KEY, 
                name TEXT
            );
        `);
    } catch (error) {
        console.error("Failed to create table", error);
        throw error;
    }
}

export function stressTest() {
    let client = cql;
    let hosts = __ENV.CQL_HOSTS || "127.0.0.1:9042";
    const keyspace = "test_keyspace";

    client.session({
        hosts: hosts.split(","),
        keyspace: keyspace,
    });

    let id = Math.floor(Math.random() * 1000000);
    let name = `User${id}`;

    client.exec(`INSERT INTO users (id, name) VALUES (uuid(), '${name}');`);

    client.close();
}
