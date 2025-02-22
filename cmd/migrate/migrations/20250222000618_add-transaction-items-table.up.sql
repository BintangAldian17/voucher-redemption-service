CREATE TABLE IF NOT EXISTS transaction_items ( 
    id INT AUTO_INCREMENT,
    transaction_id INT NOT NULL,
    voucher_id INT NOT NULL,
    quantity INT NOT NULL,
    cost_in_points INT NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_transaction_item_transaction
        FOREIGN KEY (transaction_id) REFERENCES transactions (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_transaction_item_voucher
        FOREIGN KEY (voucher_id) REFERENCES vouchers (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
