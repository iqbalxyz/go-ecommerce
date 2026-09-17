CREATE TABLE IF NOT EXISTS cart_items (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    cart_id    BIGINT UNSIGNED NOT NULL,
    product_id BIGINT UNSIGNED NOT NULL,
    quantity   INT NOT NULL DEFAULT 1,
    created_at DATETIME(3),
    updated_at DATETIME(3),
    PRIMARY KEY (id),
    UNIQUE KEY idx_cart_items_cart_product (cart_id, product_id),
    CONSTRAINT fk_cart_items_cart FOREIGN KEY (cart_id) REFERENCES carts(id) ON DELETE CASCADE,
    CONSTRAINT fk_cart_items_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;