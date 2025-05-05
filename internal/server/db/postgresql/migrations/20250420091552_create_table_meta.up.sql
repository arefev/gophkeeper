BEGIN;
CREATE TABLE IF NOT EXISTS public.meta (
    id bigint GENERATED ALWAYS AS IDENTITY NOT NULL,
    "uuid" uuid DEFAULT gen_random_uuid(),
    "user_id" bigint NOT NULL,
    "type" smallint NOT NULL,
    "name" varchar(255) NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT meta_pk PRIMARY KEY (id),
    CONSTRAINT meta_unique UNIQUE (uuid),
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);
COMMIT;