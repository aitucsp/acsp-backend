CREATE TABLE users
(
    id         bigserial PRIMARY KEY,
    email      varchar     NOT NULL,
    name       varchar     NOT NULL,
    password   varchar     NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now()),
    updated_at timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE articles
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    topic       VARCHAR     NOT NULL,
    description VARCHAR     NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT (now())
);

ALTER TABLE articles
    ADD CONSTRAINT fk_articles_users FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE CASCADE;