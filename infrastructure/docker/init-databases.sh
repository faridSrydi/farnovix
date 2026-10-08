#!/bin/bash
set -e

# Create all service databases
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE DATABASE user_db;
    CREATE DATABASE event_db;
    CREATE DATABASE booking_db;
    CREATE DATABASE payment_db;
    CREATE DATABASE ticket_db;
    CREATE DATABASE notification_db;

    -- Grant privileges
    GRANT ALL PRIVILEGES ON DATABASE user_db TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON DATABASE event_db TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON DATABASE booking_db TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON DATABASE payment_db TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON DATABASE ticket_db TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON DATABASE notification_db TO $POSTGRES_USER;
EOSQL

echo "All service databases created successfully."
