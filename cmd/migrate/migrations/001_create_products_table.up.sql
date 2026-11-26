-- Migration: 001_create_products_table
-- Creates the products table for MySQL

CREATE TABLE IF NOT EXISTS `products` (
    `id` CHAR(36) NOT NULL PRIMARY KEY,
    `sku` VARCHAR(100) NOT NULL UNIQUE,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT,
    `price_cents` INT NOT NULL CHECK (`price_cents` > 0),
    `category` VARCHAR(100) NOT NULL,
    `created_at` TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Indexes for common query patterns
CREATE INDEX `idx_products_category` ON `products` (`category`);
CREATE INDEX `idx_products_name` ON `products` (`name`);
CREATE INDEX `idx_products_created_at` ON `products` (`created_at`);