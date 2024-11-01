\connect gostream

CREATE TABLE file_schema.files (
    id          UUID                    PRIMARY KEY,
    name        file_schema.string,
    bucket      file_schema.string,
    object      file_schema.string,
    
    CONSTRAINT  unique_filename UNIQUE (filename)
);

CREATE TABLE file_schema.files_meta (
    id          file_schema.uuid_key    PRIMARY KEY,
    uploaded_at file_schema.timestamp,
    file_id     UUID    REFERENCES file_schema.files (id) ON DELETE CASCADE
);
