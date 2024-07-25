# Snappy

Snappy is a Command Line Interface (CLI) tool designed to handle snapshots for PostgreSQL databases efficiently.

## Features

- Create database snapshots
- List existing snapshots
- Remove snapshots
- Rename snapshots
- Restore snapshots

## Installation

[Add installation instructions here]

## Usage

The basic syntax for using Snappy is:

```sh
snappy <command> [options]
```

### Available Commands

- `completion`: Generate the autocompletion script for the specified shell
- `list`: List all snapshots
- `remove`: Remove a specific snapshot
- `rename`: Rename a snapshot
- `restore`: Restore a database from a snapshot
- `snapshot`: Create a new snapshot of a database

### Global Flags

- `--help`: Display help information
- `-h, --host string`: Specify PostgreSQL host
- `-p, --port string`: Specify PostgreSQL port
- `-U, --username string`: Specify PostgreSQL user

### Environment Variables

Snappy uses the following environment variables if set:

- `PGUSER`: PostgreSQL username
- `PGPASSWORD`: PostgreSQL password

These environment variables can be used instead of or in addition to the command-line flags.
