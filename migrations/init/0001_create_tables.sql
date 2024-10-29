\connect gostream

CREATE TABLE file_schema.files (
    id          UUID                    PRIMARY KEY,
    filename    file_schema.string,
    CONSTRAINT  unique_filename UNIQUE (filename)
);

CREATE TABLE file_schema.files_meta (
    id          file_schema.uuid_key    PRIMARY KEY,
    size        file_schema.positive_int,
    uploaded_at file_schema.timestamp,
    file_id     UUID    REFERENCES file_schema.files (id)
);

CREATE TABLE file_schema.files_vid (
    id          file_schema.uuid_key    PRIMARY KEY,
    location    file_schema.string,
    file_id     UUID    REFERENCES file_schema.files (id)
);

CREATE TABLE file_schema.files_manifest (
    id          file_schema.uuid_key    PRIMARY KEY,
    location    file_schema.string,
    file_id     UUID    REFERENCES file_schema.files (id)
);

CREATE TABLE file_schema.files_ts (
    id          file_schema.uuid_key    PRIMARY KEY,
    location    file_schema.string,
    file_id     UUID    REFERENCES file_schema.files (id)
);