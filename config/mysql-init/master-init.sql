-- Create the database if it doesn't exist
CREATE DATABASE IF NOT EXISTS `ai-app`;

-- Create the application user with mysql_native_password for compatibility
CREATE USER IF NOT EXISTS 'gorm'@'%' IDENTIFIED WITH mysql_native_password BY 'secret';

-- Grant all privileges on the ai-app database to 'gorm' user
GRANT ALL PRIVILEGES ON `ai-app`.* TO 'gorm'@'%';

-- Create the replication user (with replication privileges)
CREATE USER IF NOT EXISTS 'replicator'@'%' IDENTIFIED WITH mysql_native_password BY 'replica_pass';
GRANT REPLICATION SLAVE ON *.* TO 'replicator'@'%';

-- Apply changes
FLUSH PRIVILEGES;
