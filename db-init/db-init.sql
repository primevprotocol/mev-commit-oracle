CREATE TABLE committed_transactions (
    commitment_index BBYTEA PRIMARY KEY,
    transaction VARCHAR(255),
    block_number BIGINT,
    builder_address BYTEA
);