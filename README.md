# Gator, a Blog Aggregator CLI 

## Description

RSS blog aggregator using go (This is my implementation of the "Build a Blog Aggregator" boot.dev project).

## Dependencies 

  - Go toolchain (for building)
  - postgresql
    - For Debian style linux and WSL2: 
      ```
      sudo apt update
      sudo apt install postgresql postgresql-contrib
      ```
  - psql, the default postrgresql CLI client (included in postgresql)
  - [goose](https://github.com/pressly/goose) for database migrations
    ```
    go install github.com/pressly/goose/v3/cmd/goose@latest
    ```    
  - [sqlc](https://sqlc.dev/) for generating go code from sql queries (sql-to-go "bridge")
    ```
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    ```


## Installation

To install gator:
```
go install github.com/cribb/gator@latest
```

## Gator Config file

Gator uses a json config file (default name: `~/.gatorconfig.json`) located in the user's home directory that contains the database connection string/url and the current username. Its format is:

```json
{ 
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable", 
  "current_user_name":"fred"
}
```

## Setting up Gator's database

Before running gator for the first time, set up your database schema by following the migration steps located in [`sql/connection.md`](sql/connection.md).

## Running Gator and its Commands

To run `gator` after installation, merely run `gator` as the base command.
```
gator
```

A number of commands in `gator` are used on the command line in the following manner:

```
gator <command> <command arguments>
```

The table below describes each command and the appropriate usage/arguments for each, when necessary.

  ------------------------------------------------------------------
  |   **Command**   |          **Usage**         | **Description** |
  |:---------------:|:--------------------------:|:----------------|
  | help            |                            | prints help summary to stdout |
  |	login           | `login <username>`         | logs user in |
  | register        | `register <username>`      | creates user account |
  |	reset           |                            | resets gator database **(WARNING! NO undo)** |
  |	agg             |                            | aggregates posts from feeds _(starts proto-server)_|
  |	users           |                            | lists available users |
  |	feeds           |                            | lists available feeds |
  |	addfeed         | `addfeed <feedname> <url>` | adds feed to feed list |
  |	follow          | `follow <name>`            | follows feed |
  |	following       |                            | lists feeds user follows |
  |	unfollow        | `unfollow <feedname>`      | unfollows feed for user |
  |	browse          | `browse <N>`               | prints N posts to stdout (default: 2) |
  -----------------------------------------------------------------

## Build Notes

  - To migrate gator database:
    
    `goose postgres <connection string> [up|down]`
    
    - where `<connection string>` is `"postgres://postgres:postgres@localhost:5432/gator"`
    - Run goose command from `sql/schema` subdirectory
    - For more details see [sql connection info](/sql/connection.md)

  - To generate go code from sql queries:

    `sqlc generate` (while in the project root dir)
  
  - To install necessary go packages
    ``` 
      go install <packages>
    ```

  - To build gator

    ` go build . `
