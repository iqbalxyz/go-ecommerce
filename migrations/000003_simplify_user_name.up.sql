-- Migration: simplify_user_name
-- Description: Replace first_name and last_name with single name column

ALTER TABLE users ADD COLUMN name VARCHAR(150) NOT NULL AFTER email;

ALTER TABLE users DROP COLUMN first_name;
ALTER TABLE users DROP COLUMN last_name;