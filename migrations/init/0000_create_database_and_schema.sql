CREATE DATABASE gostream;

\connect gostream

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