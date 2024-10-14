\connect :DB_NAME

CREATE SCHEMA file_schema;

CREATE DOMAIN file_schema.uuid_key AS UUID 
DEFAULT gen_random_uuid() 
NOT NULL;

CREATE DOMAIN file_schema.string AS VARCHAR(256) NOT NULL;

CREATE DOMAIN file_schema.positive_int AS INTEGER
CHECK (
    VALUE > 0
);

CREATE DOMAIN file_schema.timestamp AS TIMESTAMP WITH TIME ZONE
DEFAULT (current_timestamp AT TIME ZONE 'UTC');

CREATE TABLE file_schema.files (
    id          file_schema.uuid_key,
    filename    file_schema.string,
    PRIMARY KEY (id)
);

CREATE TABLE file_schema.files_meta (
    id          file_schema.uuid_key,
    size        file_schema.positive_int,
    uploaded_at file_schema.timestamp,
    file_id     UUID    REFERENCES file_schema.files (id)
);

CREATE TABLE file_schema.files_vid (
    id          file_schema.uuid_key,
    location    file_schema.string,
    file_id     UUID    REFERENCES file_schema.files (id)
);

CREATE TABLE file_schema.files_manifest (
    id          file_schema.uuid_key,
    location    file_schema.string,
    file_id     UUID    REFERENCES file_schema.files (id)
);

CREATE TABLE file_schema.files_ts (
    id          file_schema.uuid_key,
    location    file_schema.string,
    file_id     UUID    REFERENCES file_schema.files (id)
);