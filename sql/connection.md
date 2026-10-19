# Database info for gator

## connection string for gator

`postgres://postgres:postgres@localhost:5432/gator`

## Connecting directly to postgresql cli (psql) for gator db

`psql <connection string>`

```bash
psql "postgres://postgres:postgres@localhost:5432/gator"
```


## Goose Migrations

Run from the `sql/schema` directory...

```bash
goose postgres <connection string> [up|down]
```


### up

```bash
goose postgres "postgres://postgres:postgres@localhost:5432/gator" up
```

### down

```bash
goose postgres "postgres://postgres:postgres@localhost:5432/gator" down
```


