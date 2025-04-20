BEGIN;
CREATE TABLE IF NOT EXISTS public.files (
    id bigint GENERATED ALWAYS AS IDENTITY NOT NULL,
    "file_id" bigint NOT NULL,
    "name" varchar(255) NOT NULL,
    "data" bytea NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT files_pk PRIMARY KEY (id),
    CONSTRAINT fk_files FOREIGN KEY(file_id) REFERENCES meta(id)
);
COMMIT;