#!/usr/bin/env bash

goose -dir migrations sqlite3 db/local.db up

/app/migrator
/app/bot