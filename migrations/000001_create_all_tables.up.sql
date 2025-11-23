CREATE TABLE teams (
    team_name VARCHAR(255) PRIMARY KEY
);

CREATE TABLE users (
    user_id   VARCHAR(255) PRIMARY KEY,
    team_name VARCHAR(255) REFERENCES teams(team_name) ON DELETE SET NULL,
    username  VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_users_team_name ON users(team_name);

CREATE TABLE pull_requests (
    pull_request_id   VARCHAR(255) PRIMARY KEY,
    pull_request_name VARCHAR(255) NOT NULL,
    author_id         VARCHAR(255) REFERENCES users(user_id) ON DELETE CASCADE,
    status            VARCHAR(16) NOT NULL,
    created_at        TIMESTAMPTZ NULL,
    merged_at         TIMESTAMPTZ NULL
);

CREATE INDEX idx_pull_requests_author_id ON pull_requests(author_id);

CREATE TABLE pull_request_reviewers (
    pull_request_id VARCHAR(255) REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    reviewer_id     VARCHAR(255) REFERENCES users(user_id) ON DELETE CASCADE,
    PRIMARY KEY (pull_request_id, reviewer_id)
);

CREATE INDEX idx_pr_reviewers_reviewer_id ON pull_request_reviewers(reviewer_id);
