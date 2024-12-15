-- 活動表
CREATE TABLE campaigns (
    id BIGSERIAL PRIMARY KEY,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL
);

-- 任務表
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

-- 用戶任務表
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

-- 索引
CREATE INDEX idx_user_tasks_user_id ON user_tasks(user_id);
CREATE INDEX idx_user_tasks_task_id ON user_tasks(task_id); 