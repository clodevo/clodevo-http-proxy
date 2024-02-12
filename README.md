# Clodevo Forward Proxy Documentation

The Clodevo Forward Proxy is designed to efficiently manage web traffic with a focus on security, scalability, and flexibility.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
  - [Prerequisites](#prerequisites)
  - [Building from Source](#building-from-source)
  - [Using Docker](#using-docker)
- [Configuration](#configuration)
  - [Configuration File Structure](#configuration-file-structure)
  - [Database Configuration](#database-configuration)
  - [Git Sync Configuration](#git-sync-configuration)
  - [Proxy Configuration](#proxy-configuration)
  - [Admin and ACL Configuration](#admin-and-acl-configuration)
  - [Logging Level](#logging-level)
- [Running the Proxy](#running-the-proxy)
- [API Reference](#api-reference)
  - [Authentication](#authentication)
  - [Tenant Management](#tenant-management)
  - [API Key Management](#api-key-management)
- [Database Management](#database-management)
  - [Initializing the Database](#initializing-the-database)
  - [Managing Connections](#managing-connections)
- [Git Synchronization](#git-synchronization)
  - [Setting Up Git Synchronization](#setting-up-git-synchronization)
- [Logging](#logging)
  - [Logging Levels](#logging-levels)
  - [Configuring Log Output](#configuring-log-output)
- [Troubleshooting](#troubleshooting)
  - [Database Connection Errors](#database-connection-errors)
  - [Git Sync Problems](#git-sync-problems)
  - [Proxy Routing Issues](#proxy-routing-issues)
- [License](#license)
- [Support](#support)
- [Configuring Proxy Environment Variables](#configuring-proxy-environment-variables)
- [ACL Manager Usage Guide](#acl-manager-usage-guide)

## Introduction

The Clodevo Forward Proxy is engineered for organizations seeking to manage web traffic efficiently. It emphasizes security, scalability, and flexibility, making it an ideal choice for high-performance environments.

## Features

- **Flexible Routing:** Direct requests based on rules and policies.
- **Access Control Lists (ACL):** Manage access with fine-grained control.
- **Git Sync for ACLs:** Dynamically update ACLs using Git repositories.
- **Multi-Database Support:** Compatible with SQLite, PostgreSQL, and MySQL.
- **Comprehensive Logging:** Detailed logs for effective monitoring and debugging.
- **API for Management:** RESTful API for managing tenants and API keys.
- **Docker Support:** Simplifies deployment.

## Installation

### Prerequisites

- Go 1.20+
- Docker (optional)
- Git repository access (optional for Git sync)

### Building from Source

```sh
git clone https://github.com/clodevo/raven-proxy.git
cd raven-proxy
go build -o clodevo-proxy .
```

### Using Docker

```sh
docker build -t clodevo-proxy .
```

## Configuration

Configuration is managed through the `config.json` file, encompassing database connections, Git sync settings, proxy configurations, and more.

### Configuration File Structure

```json
{
  "DatabaseConfig": {...},
  "GitSyncConfig": {...},
  "ProxyConfig": {...},
  "AdminAPIKey": "your_admin_api_key_here",
  "ACLDataPath": "./path/to/acl",
  "AdminAddr": ":9090",
  "LogLevel": "info"
}
```

#### Database Configuration

Configure the database type, connection details, and initialization parameters.

#### Git Sync Configuration

Set up Git repository details for ACL synchronization.

#### Proxy Configuration

Define proxy server settings like address, max concurrent connections, and timeouts.

#### Admin and ACL Configuration

Secure admin endpoints with an API key and specify the ACL data path.

#### Logging Level

Adjust the verbosity of logs according to your needs.

## Running the Proxy

```sh
./clodevo-proxy
```

Or using Docker:

```sh
docker run -p 8080:8080 -p 9090:9090 clodevo-proxy
```

## API Reference

### Authentication

Use the `X-Admin-API-Key` header for authenticating API requests.

### Tenant Management

Create, retrieve, update, and delete tenants through the Admin API.

### API Key Management

Manage API keys for tenants, allowing for creation, retrieval, rotation, and deletion.

## Database Management

Supports SQLite, PostgreSQL, and MySQL, with automatic table creation and efficient connection pooling.

## Git Synchronization

Facilitates dynamic ACL updates through Git repository synchronization, supporting private repositories and dedicated branches.

## Logging

Offers comprehensive logging with adjustable levels for detailed monitoring and troubleshooting.

## Troubleshooting

Addresses common issues such as database connection errors, Git sync problems, and proxy routing issues, providing solutions for each.

## License

Clodevo Forward Proxy is licensed under the Apache 2.0 License.

## Support

For support, contact support@clodevo.com or visit the support page.

## Configuring Proxy Environment Variables

Guide on setting `http_proxy` and `https_proxy` environment variables for routing traffic through the proxy.

## ACL Manager Usage Guide

Explains managing access control with ACLs, including JSON file structure and request evaluation logic.

This comprehensive guide aims to provide all necessary information to get started with the Clodevo Forward Proxy. For further details or specific use cases, please refer to the API documentation or contact our support team.