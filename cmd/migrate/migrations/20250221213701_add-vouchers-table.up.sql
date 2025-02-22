CREATE TABLE IF NOT EXISTS vouchers (
    id INT AUTO_INCREMENT,
    brand_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    cost_in_points INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY unique_voucher_name (name),
    CONSTRAINT fk_voucher_brand
        FOREIGN KEY (brand_id) REFERENCES brands (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);