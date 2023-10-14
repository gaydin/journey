-- name: create-table-settings

CREATE TABLE IF NOT EXISTS
    settings
(
    id         SERIAL PRIMARY KEY,
    uuid       varchar(36)  NOT NULL DEFAULT uuid_generate_v4(),
    key        varchar(150) NOT NULL,
    value      text,
    type       varchar(150) NOT NULL DEFAULT 'core',
    created_at timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by integer      NOT NULL,
    updated_at timestamp,
    updated_by integer
);

INSERT INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by)
VALUES (1, default, 'title', 'My Blog', 'blog', default, 1, default, 1);
INSERT INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by)
VALUES (2, default, 'description', 'Just another Blog', 'blog', default, 1, default, 1);
INSERT INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by)
VALUES (3, default, 'email', '', 'blog', default, 1, default, 1);
INSERT INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by)
VALUES (4, default, 'logo', '/public/images/blog-logo.jpg', 'blog', default, 1, default, 1);
INSERT INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by)
VALUES (5, default, 'cover', '/public/images/blog-cover.jpg', 'blog', default, 1, default, 1);
INSERT INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by)
VALUES (6, default, 'postsPerPage', 5, 'blog', default, 1, default, 1);
INSERT INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by)
VALUES (7, default, 'activeTheme', 'promenade', 'theme', default, 1, default, 1);
INSERT INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by)
VALUES (8, default, 'navigation', '[{"label":"Home", "url":"/"}]', 'blog', default, 1, default, 1);
