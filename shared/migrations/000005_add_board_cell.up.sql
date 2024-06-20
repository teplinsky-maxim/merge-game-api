BEGIN;

create table if not exists board_cells
(
    id                 bigserial primary key,
    board_id           bigint,
    cell_w             bigint,
    cell_h             bigint,
    collection_id      bigint,
    collection_item_id bigint,

    time_created       timestamp default now()
);

CREATE UNIQUE INDEX board_cells_board_id_cell_w_cell_h_idx ON board_cells (board_id, cell_w, cell_h);

COMMIT;