CREATE TABLE IF NOT EXISTS users (
    id                   UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email                TEXT        NOT NULL,
    name                 TEXT        NOT NULL,
    password_hash        TEXT        NOT NULL,
    country              CHAR(2)     NOT NULL,
    profile_picture_url  TEXT,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT users_email_unique     UNIQUE (email),
    CONSTRAINT users_email_lowercase  CHECK (email = lower(email)),
    CONSTRAINT users_name_not_blank   CHECK (length(trim(name)) > 0),
    CONSTRAINT users_country_iso      CHECK (country ~ '^[A-Z]{2}$')
);