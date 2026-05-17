--liquibase formatted sql

--changeset app:001
CREATE TABLE IF NOT EXISTS tasks (
    id          UUID PRIMARY KEY,
    title       VARCHAR(200) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status      VARCHAR(20) NOT NULL DEFAULT 'todo',
    priority    VARCHAR(10) NOT NULL DEFAULT 'medium',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_status ON tasks (status);
CREATE INDEX idx_tasks_priority ON tasks (priority);
--rollback DROP INDEX IF EXISTS idx_tasks_priority;
--rollback DROP INDEX IF EXISTS idx_tasks_status;
--rollback DROP TABLE IF EXISTS tasks;
