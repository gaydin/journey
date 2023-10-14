-- name: create-table-posts_tags

CREATE TABLE IF NOT EXISTS
    posts_tags
(
    id      SERIAL PRIMARY KEY,
    post_id integer NOT NULL,
    tag_id  integer NOT NULL
);