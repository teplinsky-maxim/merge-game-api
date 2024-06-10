BEGIN;

ALTER TABLE collection_items
    DROP CONSTRAINT unique_collection_item;

ALTER TABLE collections
    DROP CONSTRAINT unique_collection;

ALTER TABLE creation_rules
    DROP CONSTRAINT unique_creation_rule;

COMMIT;