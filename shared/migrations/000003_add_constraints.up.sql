BEGIN;

-- Only 1 item of a level can exist in a collection
ALTER TABLE collection_items
    ADD CONSTRAINT unique_collection_item UNIQUE (collection_id, level);

-- Only 1 collection with the same name can exist
ALTER TABLE collections
    ADD CONSTRAINT unique_collection UNIQUE (name);

-- Only 1 creation rule for an item, means item can not create 2 items somehow
ALTER TABLE creation_rules
    ADD CONSTRAINT unique_creation_rule UNIQUE (collection_item_id);

COMMIT;