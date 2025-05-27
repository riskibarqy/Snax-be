#!/bin/bash

# Create directories if they don't exist
mkdir -p ../sqlc/schema

# Clear existing schema files in sqlc/schema
rm -f ../sqlc/schema/*.sql

# Get all up migrations and sort them
migrations=$(find ../url-shortener/migrations -name "*.up.sql" | sort)

# Concatenate all migrations into a single schema file
echo "-- This file is auto-generated from migrations. DO NOT EDIT DIRECTLY." > ../sqlc/schema/schema.sql
echo "-- Last updated: $(date)" >> ../sqlc/schema/schema.sql
echo "" >> ../sqlc/schema/schema.sql

# Append each migration file
for migration in $migrations; do
    echo "-- Including migration: $(basename $migration)" >> ../sqlc/schema/schema.sql
    echo "" >> ../sqlc/schema/schema.sql
    cat "$migration" >> ../sqlc/schema/schema.sql
    echo "" >> ../sqlc/schema/schema.sql
    echo "" >> ../sqlc/schema/schema.sql
done

echo "Schema sync completed successfully!" 