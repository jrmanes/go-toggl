

# Load .env in order to use vars
ifneq (,$(wildcard ./.env))
    include .env
    export
endif
#### #### #### #### #### #### ####

setup_db:
		sqlite3 ${DATABASE_NAME} < ./internal/data/init.sql
