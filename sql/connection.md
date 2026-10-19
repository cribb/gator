# Database info for gator

## connection string for gator

`postgres://postgres:postgres@localhost:5432/gator`

## Connecting directly to postgresql cli (psql) for gator db

`psql <connection string>`

```bash
psql "postgres://postgres:postgres@localhost:5432/gator"
```


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

## Tables

### Users

Columns:  id, created_at, updated_at, name 

### Feeds

Columns:  id, created_at, updated_at, name, url, last_fetched_at, user_id (from JOIN) 

### Feed-follows

Columns:  id, created_at, updated_at, user_id (from JOIN with Users), feed_id (from JOIN with Feeds)

### Posts

Columns: id, created_at, updated_at, title, url, description, published_at, feed_id (from JOIN with Feeds)

