-- name: create-table-posts

CREATE TABLE IF NOT EXISTS
	posts (
		id					SERIAL PRIMARY KEY,
		uuid				varchar(36) NOT NULL,
		title				varchar(150) NOT NULL,
		slug				varchar(150) NOT NULL,
		markdown			text,
		html				text,
		image				text,
		featured			bool NOT NULL DEFAULT false,
		page				bool NOT NULL DEFAULT false,
		status				varchar(150) NOT NULL DEFAULT 'draft',
		language			varchar(6) NOT NULL DEFAULT 'en_US',
		meta_title			varchar(150),
		meta_description	varchar(200),
		author_id			integer NOT NULL,
		created_at			timestamp NOT NULL,
		created_by			integer NOT NULL,
		updated_at			timestamp,
		updated_by			integer,
		published_at		timestamp,
		published_by		integer
	);




