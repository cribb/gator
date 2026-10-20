# Gator, the CLI Blog Aggregator

## Description

RSS blog aggregator using go (my implementation of "Build a Blog Aggregator" boot.dev project).

## Gator Commands




  |   **Command**   |          **Usage**         | **Description** |
  |:---------------:|:--------------------------:|:----------------|
  |	login           |                            | logs user in |
  | register        | `register <username>`      | creates user account |
  |	reset           |                            | resets gator database **(WARNING! NO undo)** |
  |	agg             |                            | aggregates posts from feeds _(starts proto-server)_|
  |	users           | `addfeed <feedname> <url>` | lists available users |
  |	feeds           |                            | lists available feeds |
  |	addfeed         |                            | adds feed to feed list |
  |	follow          | `follow <name>`            | follows feed |
  |	following       |                            | lists feeds user follows |
  |	unfollow        | `unfollow <feedname>`      | unfollows feed for user |
  |	browse          | `browse <N>`               | prints N posts to stdout |
  

	


## Dependencies 

  - postgresql and its CLI client, psql
  - goose for database migrations
  - sqlc for generating go code from sql queries (sql-to-go "bridge")

## Build Notes

  - To migrate gator database:
    
    `goose postgres <connection string> [up|down]`
    
    - where `<connection string>` is `"postgres://postgres:postgres@localhost:5432/gator"`
    - For more details see [sql connection info](/sql/connection.md)

  - To generate go code from sql queries:

    `sqlc generate` (while in the project root dir)
