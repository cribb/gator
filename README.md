# gator

This is my implementation of an RSS blog aggregator using go that follows the accompanying boot.dev project.

## Description

## Build Notes

Uses:
  - postgresql and its CLI client, psql
  - goose for database migrations
  - sqlc, a sql-to-golang code generator

To generate go code from sql queries:
``` sqlc generate``` while in the project root dir
