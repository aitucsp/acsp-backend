CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY,
    email      VARCHAR     NOT NULL,
    name       VARCHAR     NOT NULL,
    password   VARCHAR     NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    roles      VARCHAR[]
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

CREATE TABLE roles
(
    id   INT          NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULl
);

CREATE TABLE user_roles
(
    id      BIGSERIAL NOT NULL PRIMARY KEY,
    user_id INT       NOT NULL,
    role_id INT       NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE ,
    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE ON UPDATE CASCADE
);

INSERT INTO roles (id, name)
VALUES (1, 'user');
INSERT INTO roles (id, name)
VALUES (1, 'admin');

ALTER TABLE articles
    ADD CONSTRAINT fk_articles_users FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE CASCADE ON UPDATE CASCADE;