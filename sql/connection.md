# Database info for gator

## Connection string

 Gator uses a "connection string" that specifies the connection to its postgresql database. Here, that string is:

```
postgres://postgres:postgres@localhost:5432/gator
```

## Connecting directly to postgresql cli (psql) for gator db

To connect _directly_ to the gator database throu the postgresql cli, run:

```
psql <connection string>
```

Here, that full command would look like:

```bash
psql "postgres://postgres:postgres@localhost:5432/gator"
```

To get help with `psql` type "help". Commands in `psql` are preceded by a backslash, and "\h" will tell you all you need to know.

## Goose Migrations

Note: run goose from the `sql/schema` directory...

```bash
goose postgres <connection string> [up|down]
```

### Migrate up

```bash
goose postgres "postgres://postgres:postgres@localhost:5432/gator" up
```

### Migrate down

```bash
goose postgres "postgres://postgres:postgres@localhost:5432/gator" down
```

## Tables/Schema

### Users

|   Column   |    Type    | Description             |
|:----------:|:----------:|:------------------------|
| id         |  UUID      | primary key (unique id) |
| created_at |  TIMESTAMP | record's creation time  |
| updated_at |  TIMESTAMP | record's update time    |
| name       |  TEXT      | user name               |


### Feeds

|     Column      |    Type    | Description                 |
|:---------------:|:----------:|:----------------------------|
| id              |  UUID      | primary key (unique id)     |
| created_at      |  TIMESTAMP | record's creation time      |
| updated_at      |  TIMESTAMP | record's update time        |
| name            |  TEXT      | rss feed name               |
| url             |  TEXT      | rss feed url                |
| last_fetched_at |  TIMESTAMP | feed's last fetch time      |
| user_id         |  UUID      | users.id (JOIN from `users`)|


### Feed-follows

|     Column      |    Type    | Description                 |
|:---------------:|:----------:|:----------------------------|
| id              |  UUID      | primary key (unique id)     |
| created_at      |  TIMESTAMP | record's creation time      |
| updated_at      |  TIMESTAMP | record's update time        |
| user_id         |  UUID      | users.id (JOIN from `users`)|
| feed_id         |  UUID      | feeds.id (JOIN from `feeds`)|


### Posts

|     Column      |    Type    | Description                  |
|:---------------:|:----------:|:-----------------------------|
| id              |  UUID      | primary key (unique id)      |
| created_at      |  TIMESTAMP | record's creation time       |
| updated_at      |  TIMESTAMP | record's update time         |
| title           |  TEXT      | blog post title              |
| url             |  TEXT      | blog post url                |
| description     |  TEXT      | blog post description        |
| published_at    |  TIMESTAMP | blog post publication time   |
| feed_id         |  UUID      | feeds.id (JOIN from `feeds`) |

