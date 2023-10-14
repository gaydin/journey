-- name: create-table-roles-users

CREATE TABLE IF NOT EXISTS
    roles_users
(
    id      SERIAL PRIMARY KEY,
    role_id integer NOT NULL,
    user_id integer NOT NULL
);