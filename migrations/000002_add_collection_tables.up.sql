BEGIN;

CREATE TABLE IF NOT EXISTS collection_items
(
    id            bigserial primary key,
    collection_id BIGINT,
    name          TEXT,
    level         BIGINT,
    mergeable     BOOLEAN,
    can_create    BOOLEAN
);

CREATE TABLE IF NOT EXISTS creation_rules
(
    id                          bigserial primary key,
    collection_item_id          BIGINT,
    generate_collection_item_id BIGINT
);

COMMIT;