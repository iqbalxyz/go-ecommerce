-- Rollback: simplify_user_name
-- Restore first_name and last_name from name

ALTER TABLE users ADD COLUMN first_name VARCHAR(100) NOT NULL DEFAULT '' AFTER email;
ALTER TABLE users ADD COLUMN last_name VARCHAR(100) DEFAULT NULL AFTER first_name;

ALTER TABLE users DROP COLUMN name;