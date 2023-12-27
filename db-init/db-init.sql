CREATE TABLE committed_transactions (
    commitment_index BINARY(32) PRIMARY KEY,
    transaction VARCHAR(255),
    block_number BIGINT,
    builder_address BINARY(32)
);