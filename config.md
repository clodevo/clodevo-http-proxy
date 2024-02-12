# Application Configuration Guide

This document outlines the necessary steps and provides an example configuration for setting up the application using the `config.json` file. The `config.json` file centralizes the configuration for the database, Git synchronization, proxy settings, administrative API key, ACL data path, administrative server address, and logging level.

## Config.json Structure

The `config.json` file consists of several sections, each responsible for different aspects of the application configuration:

### Database Configuration

- **Type:** Database type (`sqlite3`, `postgres`, `mysql`, etc.).
- **FilePath:** File path for SQLite database.
- **Host:** Host address for PostgreSQL and MySQL.
- **Port:** Port number for PostgreSQL and MySQL.
- **User:** Username for database access.
- **Password:** Password for database access.
- **DBName:** Database name.

### Git Sync Configuration

- **repo-url:** URL of the Git repository for ACL synchronization.
- **branch-name:** Branch name to use in the Git repository.
- **username:** Username for Git repository access.
- **password:** Password for Git repository access.
- **repo-path:** Local file system path where the Git repository will be cloned.
- **sync-interval:** Interval between synchronization operations (e.g., `"1m"` for 1 minute).

### Proxy Configuration

- **addr:** Address the proxy server listens on.
- **maxConcurrent:** Maximum number of concurrent connections.
- **dns:** List of DNS servers for the proxy to use.
- **timeout:** Timeout for proxy connections.

### Admin and ACL Configuration

- **admin-api-key:** API key for securing admin endpoints.
- **acl-data-path:** File system path to ACL (Access Control List) data.
- **admin-addr:** Address on which the admin server listens.

### Logging Level

- **log-Level:** Specifies the logging level (`info`, `debug`, `warn`, `error`).

### Example config.json

```json
{
  "database": {
    "Type": "sqlite3",
    "FilePath": "/opt/clodevo/data.db",
    "Host": "",
    "Port": 3306,
    "User": "",
    "Password": "",
    "DBName": "clod-proxy"
  },
  "git-acl": {
    "repo-url": "https://example.com/git-repo.git",
    "branch-name": "main",
    "username": "gituser",
    "password": "gitpassword",
    "repo-path": "/opt/clodevo/acl/tenants",
    "sync-interval": "1m"
  },
  "proxy": {
    "addr": ":8080",
    "maxConcurrent": 512,
    "dns": [],
    "timeout": "20s"
  },
  "admin-api-key": "your_admin_api_key_here",
  "acl-data-path": "/opt/clodevo/acl/tenants",
  "admin-addr": ":9090",
  "log-Level": "info"
}
```

### Configuration Instructions

1. Copy the example `config.json` provided above into your application's root directory.
2. Replace placeholders and adjust values according to your specific environment and requirements.
3. Ensure the application has read access to the `config.json` file at runtime.

### Further Assistance

For further assistance or if you encounter issues with the configuration, please consult the application documentation or contact support.


## Global Config Environment Variables

This table details each configuration option within `AppConfig`, its corresponding environment variable (if applicable), and a brief description along with any default values provided:

| Configuration Option | Environment Variable | Description                                                        | Default Value                |
|----------------------|----------------------|--------------------------------------------------------------------|------------------------------|
| AdminAPIKey          | ADMIN_API_KEY        | The API key for securing admin endpoints.                          | (none)                       |
| ACLDataPath          | ACL_DATA_PATH        | The file system path to ACL (Access Control List) data.            | `/opt/clodevo/acl/tenants`   |
| AdminAddr            | ADMIN_ADDR           | The address on which the admin server listens.                     | `:9090`                      |
| LogLevel             | LOG_LEVEL            | The logging level of the application.                              | `info`                       |
| DatabaseConfig       | (various)            | Embedded struct for database configuration. Uses its own set of environment variables as described earlier. | (see DatabaseConfig table) |
| ProxyConfig          | (various)            | Embedded struct for proxy configuration. Uses its own set of environment variables as described earlier.   | (see ProxyConfig table)   |
| GitSyncConfig        | (various)            | Embedded struct for Git synchronization configuration. Uses its own set of environment variables as described earlier. | (see GitSyncConfig table) |

The `AppConfig` structure aggregates configurations for different aspects of the application, including database settings, proxy server settings, Git synchronization settings, and administrative controls. The `LoadAppConfig` function initializes these configurations by loading them from a JSON configuration file and environment variables, with a fallback to default values for certain parameters if they are not explicitly set. This setup facilitates a flexible and dynamic configuration approach, allowing easy adjustments without needing to recompile the application.

## ProxyConfig
Below is the `ProxyConfig` structure represented as a table, detailing each configuration option, its environment variable, and a brief description, along with default values:

| Configuration Option | Environment Variable    | Description                                                      | Default Value         |
|----------------------|-------------------------|------------------------------------------------------------------|-----------------------|
| Addr                 | PROXY_ADDR              | The address the proxy server listens on.                         | `:8080`               |
| MaxConcurrent        | PROXY_MAXCONCURRENT     | The maximum number of concurrent connections the proxy supports. | `512`                 |
| DNS                  | PROXY_DNS               | A list of DNS servers for the proxy to use.                      | `""` (empty string)   |
| Timeout              | PROXY_TIMEOUT           | The timeout for proxy connections.                               | `20s` (20 seconds)    |

This table reflects the configuration options available for the proxy server functionality within the application. The environment variables correspond to the specific settings that can be adjusted to customize the behavior of the proxy. Default values are provided and will be used if the respective environment variables are not set, ensuring the proxy has sensible defaults to fall back on.


## DatabaseConfig

Below is the `DatabaseConfig` structure represented as a table, detailing each configuration option, its environment variable, and a brief description:

| Configuration Option | Environment Variable    | Description                                                                 | Default Value          |
|----------------------|-------------------------|-----------------------------------------------------------------------------|------------------------|
| Type                 | DATABASE_TYPE           | The type of database (e.g., `sqlite3`, `postgres`, `mysql`).                | `sqlite3`              |
| FilePath             | DATABASE_FILEPATH       | The file path for the SQLite database.                                      | `/opt/clodevo/data.db` |
| Host                 | DATABASE_HOST           | The host address for PostgreSQL and MySQL databases.                        | (none)                 |
| Port                 | DATABASE_PORT           | The port number for PostgreSQL and MySQL databases.                         | `3306`                 |
| User                 | DATABASE_USER           | The username for accessing PostgreSQL and MySQL databases.                  | (none)                 |
| Password             | DATABASE_PASSWORD       | The password for accessing PostgreSQL and MySQL databases.                  | (none)                 |
| DBName               | DATABASE_DBNAME         | The name of the database to connect to in PostgreSQL and MySQL databases.   | `clod-proxy`           |

This table reflects the configuration options available in the `DatabaseConfig` structure and corresponds to the environment variables used to configure the database connection parameters. The default values are provided for some of the options and will be used if the respective environment variables are not set.

## GitSyncConfig

Below is the `GitSyncConfig` structure represented as a table, detailing each configuration option, its environment variable, and a brief description:

| Configuration Option | Environment Variable      | Description                                                      | Default Value |
|----------------------|---------------------------|------------------------------------------------------------------|---------------|
| RepoURL              | GIT_ACL_REPO_URL          | The URL of the Git repository to be synchronized.                | (none)        |
| BranchName           | GIT_ACL_BRANCH_NAME       | The name of the branch in the Git repository to use.             | (none)        |
| Username             | GIT_ACL_USERNAME          | The username for authentication with the Git repository.         | (none)        |
| Password             | GIT_ACL_PASSWORD          | The password for authentication with the Git repository.         | (none)        |
| RepoPath             | GIT_ACL_REPO_PATH         | The local filesystem path where the Git repository will be cloned. | `/opt/clodevo/acl`  |
| SyncInterval         | GIT_ACL_SYNC_INTERVAL     | The interval between synchronization operations.                 | `1m` (1 minute) |

This table maps directly to the `GitSyncConfig` structure and the corresponding environment variables used to configure the Git synchronization process. Remember, the default values are used only if the environment variables are not set.