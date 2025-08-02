-- Create the database if it doesn't exist
CREATE DATABASE IF NOT EXISTS `ai-app`;

-- Create the application user with mysql_native_password for compatibility
CREATE USER IF NOT EXISTS 'gorm'@'%' IDENTIFIED WITH mysql_native_password BY 'secret';

-- Grant SELECT (read-only) permissions on the ai-app database to 'gorm'
GRANT SELECT ON `ai-app`.* TO 'gorm'@'%';

-- Apply changes
FLUSH PRIVILEGES;

-- Sleep to give some time for the master to be fully available
SELECT SLEEP(10);

-- Set up replication, connecting to the master using the replication user
CHANGE MASTER TO
    MASTER_HOST='master',  -- Hostname of the master (can be the Docker service name or IP)
    MASTER_PORT=3306,
    MASTER_USER='replicator',  -- Use the replicator user for replication
    MASTER_PASSWORD='replica_pass',  -- Password of the replication user
    MASTER_LOG_FILE='mysql-bin.000001',  -- This should be dynamically fetched in real scenarios
    MASTER_LOG_POS=0;  -- Position should be dynamically set based on the master's binlog

-- Start replication on the slave
START SLAVE;
