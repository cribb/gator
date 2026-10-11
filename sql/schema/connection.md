# Database info for gator

## connection string for gator

`postgres://postgres:postgres@localhost:5432/gator`

## Migrations

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


