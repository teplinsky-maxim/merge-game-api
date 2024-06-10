BEGIN;

create table if not exists collections
(
    id   bigserial primary key,
    name text
);

COMMIT;