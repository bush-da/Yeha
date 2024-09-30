-- Creates a MySQL server with:
--   Database yeha_dev_db.
--   User yeha_dev with password yeha_dev_pwd in localhost.
--   Grants all privileges for yeha_dev on yeha_dev_db.
--   Grants SELECT privilege for yeha_dev on performance_schema.

CREATE DATABASE IF NOT EXISTS yeha_dev_db;
CREATE USER
    IF NOT EXISTS 'yeha_dev'@'localhost'
    IDENTIFIED BY 'yeha_dev_pwd';
GRANT ALL PRIVILEGES
   ON `yeha_dev_db`.*
   TO 'yeha_dev'@'localhost';

GRANT SELECT
   ON `performance_schema`.*
   TO 'yeha_dev'@'localhost';

FLUSH PRIVILEGES;
