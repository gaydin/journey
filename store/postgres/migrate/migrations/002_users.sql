-- name: create-table-users

CREATE TABLE IF NOT EXISTS
    users
(
    id               SERIAL PRIMARY KEY,
    uuid             varchar(36)  NOT NULL,
    name             varchar(150) NOT NULL,
    slug             varchar(150) NOT NULL,
    password         varchar(60)  NOT NULL,
    email            varchar(254) NOT NULL,
    image            text,
    cover            text,
    bio              varchar(200),
    website          text,
    location         text,
    accessibility    text,
    status           varchar(150) NOT NULL DEFAULT 'active',
    language         varchar(6)   NOT NULL DEFAULT 'en_US',
    meta_title       varchar(150),
    meta_description varchar(200),
    last_login       timestamp,
    created_at       timestamp    NOT NULL,
    created_by       integer      NOT NULL,
    updated_at       timestamp,
    updated_by       integer
);