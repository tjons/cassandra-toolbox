CREATE TABLE test.table1 (
    entity_name text,
    entity_attribute_1 text,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY (entity_name)
);

CREATE INDEX ON test.table1 (entity_attribute_1) USING 'SAI';