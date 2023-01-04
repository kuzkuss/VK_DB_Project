.PHONY: start_bd, init_db, create_tables, drop_tables

start:
	docker-compose up -d

build-server:
	docker-compose build server

create_tables:
	psql postgresql://db_pg:db_postgres@localhost:5000/postgres -f db/db.sql

drop_tables:
	psql postgresql://db_pg:db_postgres@localhost:5000/postgres -f SQL_drop/drop_all.sql

