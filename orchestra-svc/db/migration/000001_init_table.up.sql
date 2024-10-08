
-- Workflows
CREATE TABLE workflows (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Steps
CREATE TABLE steps (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(255) NOT NULL,
    service VARCHAR(100) NOT NULL,
    topic VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE payload_keys (
    id SERIAL PRIMARY KEY,
    step_id INTEGER NOT NULL REFERENCES steps(id),
    key VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- State_Actions
CREATE TABLE state_actions (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL REFERENCES workflows(type),
    state VARCHAR(255) NOT NULL,
    step_id INTEGER NOT NULL REFERENCES steps(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Workflow Instances
CREATE TABLE workflow_instances (
    id VARCHAR PRIMARY KEY,
    workflow_id INTEGER NOT NULL REFERENCES workflows(id),
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Workflow Instance Steps
CREATE TABLE workflow_instance_steps (
    id SERIAL PRIMARY KEY,
    event_id VARCHAR NOT NULL,
    status_code INTEGER,
    response VARCHAR,
    workflow_instance_id VARCHAR NOT NULL REFERENCES workflow_instances(id),
    step_id INTEGER NOT NULL REFERENCES steps(id),
    status VARCHAR(20) NOT NULL,
    event_message TEXT,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE
);

-- Process Log
CREATE TABLE process_logs (
    id SERIAL PRIMARY KEY,
    event_id VARCHAR NOT NULL,
    workflow_instance_id VARCHAR NOT NULL,
    state VARCHAR(255) NOT NULL,
    status_code INTEGER,
    status VARCHAR(20) NOT NULL,
    event_message TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
