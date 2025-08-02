#!/bin/bash

set -e

echo "Waiting for master to be ready..."
until mysqladmin ping -hmaster -uroot -proot --silent; do
  sleep 2
done

echo "Fetching master log position..."
MASTER_STATUS=$(mysql -hmaster -uroot -proot -e "SHOW MASTER STATUS\G")
LOG_FILE=$(echo "$MASTER_STATUS" | grep File | awk '{print $2}')
LOG_POS=$(echo "$MASTER_STATUS" | grep Position | awk '{print $2}')

echo "Configuring slave..."
mysql -uroot -proot <<-EOSQL
  STOP SLAVE;
  CHANGE MASTER TO
    MASTER_HOST='master',
    MASTER_USER='replicator',
    MASTER_PASSWORD='replica_pass',
    MASTER_LOG_FILE='$LOG_FILE',
    MASTER_LOG_POS=$LOG_POS,
    MASTER_PORT=3306;
  START SLAVE;
EOSQL

echo "Replica setup complete."
