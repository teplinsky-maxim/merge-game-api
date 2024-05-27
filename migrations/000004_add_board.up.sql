BEGIN;

create table if not exists boards
(
    id     bigserial primary key,
    width  int,
    height int
);

create table if not exists tasks
(
    id                     bigserial primary key,
    uuid                   uuid,
    type                   varchar,
    status                 int,
    args                   jsonb,
    result                 jsonb,

    time_created           timestamp default now(),
    time_started_executing timestamp,
    time_done_executing    timestamp
);

COMMIT;