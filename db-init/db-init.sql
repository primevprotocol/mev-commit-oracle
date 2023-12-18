CREATE TABLE committed_transactions (
    transaction VARCHAR(255),
    block_number BIGINT,
    PRIMARY KEY (transaction, block_number)
);