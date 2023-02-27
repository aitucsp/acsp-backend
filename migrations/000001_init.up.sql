CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY,
    email      VARCHAR     NOT NULL,
    name       VARCHAR     NOT NULL,
    password   VARCHAR     NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    is_admin   BOOL        NOT NULL DEFAULT FALSE,
    roles      VARCHAR[]            DEFAULT ARRAY ['user']
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
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE scholar_articles
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    topic       VARCHAR     NOT NULL,
    description VARCHAR     NOT NULL,
    upvote      BIGINT      NOT NULL DEFAULT 0,
    downvote    BIGINT      NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE scholar_comments
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    article_id BIGINT      NOT NULL,
    parent_id  BIGINT,
    text       VARCHAR     NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (article_id) REFERENCES scholar_articles (id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT "Comment_parent_id_fkey"
        FOREIGN KEY (parent_id) REFERENCES scholar_comments (id)
            ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE code_connection_cards
(
    id          BIGSERIAL   NOT NULL PRIMARY KEY,
    user_id     INT         NOT NULL,
    position    VARCHAR     NOT NULL,
    skills      VARCHAR[]   NOT NULL,
    description VARCHAR     NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE code_connection_invitations
(
    id         BIGSERIAL   NOT NULL PRIMARY KEY,
    card_id    INT         NOT NULL,
    inviter_id INT         NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY (inviter_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (card_id) REFERENCES code_connection_cards (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE code_connection_responses
(
    id            BIGSERIAL   NOT NULL PRIMARY KEY,
    invitation_id INT         NOT NULL,
    status        VARCHAR     NOT NULL DEFAULT ('NOT ANSWERED'),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY (invitation_id) REFERENCES code_connection_invitations (id) ON DELETE CASCADE ON UPDATE CASCADE
);

INSERT INTO roles (id, name)
VALUES (1, 'user');
INSERT INTO roles (id, name)
VALUES (2, 'admin');

