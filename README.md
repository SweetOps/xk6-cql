# xk6-cql

### Configuration options

The following table describes the available configuration options for connecting to a Cassandra cluster:

| Option Name       | Type        | Default Value | Required | Description |
|------------------|------------|--------------|----------|-------------|
| **`hosts`** | `[]string` | `[]` | Yes | A list of Cassandra node addresses to connect to. |
| **`keyspace`** | `string` | `""` | Yes | The keyspace to use for queries. |
| **`username`** | `string` | `""` | No | The username for authentication. |
| **`password`** | `string` | `""` | No | The password for authentication. |
| **`protocolVersion`** | `int` | `""` | No | The Cassandra protocol version to use. |
| **`timeout`** | `duration` | `"10s"` | No | The maximum duration to wait before a connection attempt times out. |
| **`consistency`** | `string` | `"quorum"` | No | The consistency level for queries (Possible values: `all`, `any`, `one`, `two`, `three`, `each_quorum`, `quorum`, `local_one`, `local_quorum`). |
| **`tls`** | `tls` | `{}` | No | The TLS settings for secure connections. |


#### TLS Configuration

| Option Name                  | Type     | Default Value | Required | Description |
|------------------------------|---------|--------------|----------|-------------|
| **`cert_path`** | `string` | `""` | No | Path to the client certificate file. |
| **`key_path`** | `string` | `""` | No | Path to the client private key file. |
| **`ca_path`** | `string` | `""` | No | Path to the CA certificate file for verifying the server. |
| **`enable_host_verification`** | `bool` | `false` | No | Whether to verify the server hostname in the TLS certificate. |


### Available Methods

#### `session()`

Initializes a new session for interacting with the Cassandra database.

```javascript
let client = cql;

client.session({
    hosts: ["10.10.0.1"],
    keyspace: "system",
});
```

#### `exec(query string)`
Executes a CQL query.

```javascript
let client = cql;

client.session({
    hosts: ["10.10.0.1"],
    keyspace: "system",
});

client.exec(`
    CREATE KEYSPACE IF NOT EXISTS test_keyspace 
    WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
`);
```
#### `batch(batchType string, queries []string)`
Executes a batch of CQL queries with the specified batch type.

`batchType` – The type of batch to execute. Possible values: `logged`, `unlogged`, `counter`,

```javascript
let client = cql;

client.session({
    hosts: ["10.10.0.1"],
    keyspace: "system",
});

client.batch("", [
        `SELECT * FROM users WHERE id = 'test';`,
        `SELECT * FROM users WHERE name = 'test';`,
    ])
```

#### `close()`
Closes the active session and releases resources.

```javascript
let client = cql;

client.session({
    hosts: ["10.10.0.1"],
    keyspace: "system",
});

client.close()
```


### Build for development

```sh
git clone git@github.com:SweetOps/xk6-cql.git
cd xk6-cql
xk6 build --with github.com/SweetOps/xk6-cql@latest=.
```
