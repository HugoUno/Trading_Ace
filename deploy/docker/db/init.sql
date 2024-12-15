CREATE DATABASE trading_ace;
\c trading_ace;

-- 複製 scripts/migrations/001_init_schema.sql 的內容
CREATE TABLE campaigns (
    id BIGSERIAL PRIMARY KEY,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL
);

CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,
    campaign_id BIGINT REFERENCES campaigns(id),
    type VARCHAR(20) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    pool_address VARCHAR(42),
    points_pool BIGINT NOT NULL,
    CONSTRAINT valid_task_type CHECK (type IN ('ONBOARDING', 'SHARE_POOL'))
);

CREATE TABLE user_tasks (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(42) NOT NULL,
    task_id BIGINT REFERENCES tasks(id),
    status VARCHAR(20) NOT NULL,
    amount DECIMAL(20,8) NOT NULL DEFAULT 0,
    points BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_status CHECK (status IN ('PENDING', 'COMPLETED'))
);

CREATE INDEX idx_user_tasks_user_id ON user_tasks(user_id);
CREATE INDEX idx_user_tasks_task_id ON user_tasks(task_id); 