-- name: create-table-roles

CREATE TABLE IF NOT EXISTS
    roles
(
    id          SERIAL PRIMARY KEY,
    uuid        varchar(36)  NOT NULL DEFAULT uuid_generate_v4(),
    name        varchar(150) NOT NULL,
    description varchar(200),
    created_at  timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by  integer      NOT NULL,
    updated_at  timestamp,
    updated_by  integer
);


INSERT INTO roles (id, uuid, name, description, created_at, created_by, updated_at, updated_by)
VALUES (1, default, 'Administrator', 'Administrators', default, 1, default, 1);
INSERT INTO roles (id, uuid, name, description, created_at, created_by, updated_at, updated_by)
VALUES (2, default, 'Editor', 'Editors', default, 1, default, 1);
INSERT INTO roles (id, uuid, name, description, created_at, created_by, updated_at, updated_by)
VALUES (3, default, 'Author', 'Authors', default, 1, default, 1);
INSERT INTO roles (id, uuid, name, description, created_at, created_by, updated_at, updated_by)
VALUES (4, default, 'Owner', 'Blog Owner', default, 1, default, 1);
