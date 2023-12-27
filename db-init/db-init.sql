CREATE TABLE committed_transactions (
    commitment_index BYTEA PRIMARY KEY,
    transaction VARCHAR(255),
    block_number BIGINT,
    builder_address BYTEA
);