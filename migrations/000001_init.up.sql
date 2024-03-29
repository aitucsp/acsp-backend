CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY,
    email      VARCHAR     NOT NULL,
    name       VARCHAR     NOT NULL,
    password   VARCHAR     NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    is_admin   BOOL        NOT NULL DEFAULT FALSE,
    roles      VARCHAR[]            DEFAULT ARRAY ['user'],
    image_url  VARCHAR     NOT NULL DEFAULT '/default'
);

CREATE TABLE user_details
(
    id             BIGSERIAL PRIMARY KEY,
    user_id        BIGINT      NOT NULL,
    first_name     VARCHAR     NOT NULL,
    last_name      VARCHAR     NOT NULL,
    email          VARCHAR     NOT NULL,
    phone_number   VARCHAR     NOT NULL,
    specialization VARCHAR     NOT NULL,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
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
    image_url   VARCHAR     NOT NULL DEFAULT '/default',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE scholar_materials
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
    upvote     BIGINT      NOT NULL DEFAULT 0,
    downvote   BIGINT      NOT NULL DEFAULT 0,
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

CREATE TABLE scholar_article_comment_votes
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL,
    comment_id BIGINT NOT NULL,
    vote_type  INT    NOT NULL,
    FOREIGN KEY (comment_id) REFERENCES scholar_comments (id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE CASCADE ON UPDATE CASCADE
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
    status     VARCHAR     NOT NULL DEFAULT ('NOT ANSWERED'),
    feedback   VARCHAR     NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY (inviter_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (card_id) REFERENCES code_connection_cards (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE contests
(
    id           BIGSERIAL   NOT NULL PRIMARY KEY,
    contest_name VARCHAR     NOT NULL,
    description  VARCHAR     NOT NULL,
    link         VARCHAR     NOT NULL,
    start_date   TIMESTAMPTZ NOT NULL DEFAULT (now()),
    end_date     TIMESTAMPTZ NOT NULL DEFAULT (now()),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE coding_lab_disciplines
(
    id          BIGSERIAL   NOT NULL PRIMARY KEY,
    title       VARCHAR     NOT NULL,
    description VARCHAR     NOT NULL,
    image_url   VARCHAR     NOT NULL DEFAULT ('/default'),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE coding_lab_projects
(
    id            BIGSERIAL   NOT NULL PRIMARY KEY,
    discipline_id INT         NOT NULL,
    title         VARCHAR     NOT NULL,
    description   VARCHAR     NOT NULL,
    level         VARCHAR     NOT NULL,
    image_url     VARCHAR     NOT NULL DEFAULT ('/default'),
    work_hours    INT         NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY (discipline_id) REFERENCES coding_lab_disciplines (id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE coding_lab_project_modules
(
    id            BIGSERIAL   NOT NULL PRIMARY KEY,
    project_id    INT         NOT NULL,
    title         VARCHAR     NOT NULL,
    description   VARCHAR     NOT NULL,
    reference_url VARCHAR     NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY (project_id) REFERENCES coding_lab_projects (id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE coding_lab_project_enrollments
(
    id         BIGSERIAL   NOT NULL PRIMARY KEY,
    project_id INT         NOT NULL,
    user_id    INT         NOT NULL,
    status     VARCHAR     NOT NULL DEFAULT ('PENDING'),
    feedback   VARCHAR     NOT NULL DEFAULT (''),
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY (project_id) REFERENCES coding_lab_projects (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE courses
(
    id          BIGSERIAL   NOT NULL PRIMARY KEY,
    author_id   INT         NOT NULL,
    title       VARCHAR     NOT NULL,
    description VARCHAR     NOT NULL,
    rating      INT         NOT NULL DEFAULT 0,
    image_url   VARCHAR     NOT NULL DEFAULT ('/default'),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY (author_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE course_modules
(
    id              BIGSERIAL   NOT NULL PRIMARY KEY,
    course_id       INT         NOT NULL,
    title           VARCHAR     NOT NULL,
    expected_result VARCHAR     NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE course_module_lessons
(
    id            BIGSERIAL   NOT NULL PRIMARY KEY,
    module_id     INT         NOT NULL,
    title         VARCHAR     NOT NULL,
    description   VARCHAR     NOT NULL,
    reference_url VARCHAR     NOT NULL default '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY (module_id) REFERENCES course_modules (id) ON DELETE CASCADE ON UPDATE CASCADE
);

INSERT INTO coding_lab_disciplines
    (title, description, image_url) VALUES ('Java',
                                            'Java is a high-level, class-based, object-oriented programming language that is designed to have as few implementation dependencies as possible.',
                                            '/java');

INSERT INTO coding_lab_disciplines
    (title, description, image_url) VALUES ('Python',
                                            'Python is an interpreted, high-level and general-purpose programming language. Python design philosophy emphasizes code readability with its notable use of significant indentation.',
                                            '/python');

INSERT INTO coding_lab_disciplines
    (title, description, image_url) VALUES ('React',
                                            'React is an open-source front-end JavaScript library for building user interfaces or UI components. It is maintained by Facebook and a community of individual developers and companies.',
                                            '/react');

CREATE TABLE course_lesson_comments
(
    id         BIGSERIAL   NOT NULL PRIMARY KEY,
    lesson_id  INT         NOT NULL,
    user_id    INT         NOT NULL,
    text    VARCHAR     NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (lesson_id) REFERENCES course_module_lessons (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE course_lesson_comment_answers
(
    id         BIGSERIAL   NOT NULL PRIMARY KEY,
    comment_id INT         NOT NULL,
    answer     VARCHAR     NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY (comment_id) REFERENCES course_lesson_comments (id) ON DELETE CASCADE ON UPDATE CASCADE
);

INSERT INTO roles (id, name)
VALUES (1, 'user');
INSERT INTO roles (id, name)
VALUES (2, 'admin');

INSERT INTO contests (contest_name, description, link, start_date, end_date)
VALUES ('LeetCode Weekly Contest 348',
        'This LeetCode contest is sponsored by LeetCode.', 'https://leetcode.com/contest/weekly-contest-348/',
        '2023-06-03 02:30:00',
        '2023-06-03 04:30:00');

INSERT INTO contests (contest_name, description, link, start_date, end_date)
VALUES ('LeetCode Biweekly Contest 106',
        'This LeetCode contest is sponsored by LeetCode.', 'https://leetcode.com/contest/biweekly-contest-106/',
        '2023-06-03 14:30:00',
        '2023-06-10 16:30:00');

