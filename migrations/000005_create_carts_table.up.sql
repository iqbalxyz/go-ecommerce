CREATE TABLE IF NOT EXISTS carts (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id    BIGINT NOT NULL,
    created_at DATETIME(3),
    updated_at DATETIME(3),
    PRIMARY KEY (id),
    UNIQUE KEY idx_carts_user_id (user_id),
    CONSTRAINT fk_carts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
