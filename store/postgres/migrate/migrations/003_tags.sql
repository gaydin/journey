-- name: create-table-tags

CREATE TABLE IF NOT EXISTS
    tags
(
    id               SERIAL PRIMARY KEY,
    uuid             varchar(36)  NOT NULL,
    name             varchar(150) NOT NULL,
    slug             varchar(150) NOT NULL,
    description      varchar(200),
    parent_id        integer,
    meta_title       varchar(150),
    meta_description varchar(200),
    created_at       timestamp    NOT NULL,
    created_by       integer      NOT NULL,
    updated_at       timestamp,
    updated_by       integer
);