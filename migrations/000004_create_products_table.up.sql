CREATE TABLE IF NOT EXISTS products (
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name        VARCHAR(200) NOT NULL,
    description TEXT,
    price       BIGINT NOT NULL DEFAULT 0,
    stock       INT NOT NULL DEFAULT 0,
    sku         VARCHAR(50) NOT NULL,
    is_active   TINYINT(1) NOT NULL DEFAULT 1,
    created_at  DATETIME(3),
    updated_at  DATETIME(3),
    deleted_at  DATETIME(3),
    PRIMARY KEY (id),
    UNIQUE KEY idx_products_sku (sku),
    KEY idx_products_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;