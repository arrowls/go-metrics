#!/bin/bash

sleep 2
migrate -path /app/migrations -database "${DATABASE_DSN}" up