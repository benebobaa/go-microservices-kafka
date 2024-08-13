-- Workflow Types
CREATE TABLE workflow_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Workflows
CREATE TABLE workflows (
    id SERIAL PRIMARY KEY,
    workflow_type_id INTEGER NOT NULL REFERENCES workflow_types(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Steps
CREATE TABLE steps (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Workflow Steps
CREATE TABLE workflow_steps (
    id SERIAL PRIMARY KEY,
    workflow_id INTEGER NOT NULL REFERENCES workflows(id),
    step_id INTEGER NOT NULL REFERENCES steps(id),
    order_index INTEGER NOT NULL,
    UNIQUE (workflow_id, order_index)
);

-- Workflow Instances
CREATE TABLE workflow_instances (
    id SERIAL PRIMARY KEY,
    workflow_id INTEGER NOT NULL REFERENCES workflows(id),
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Workflow Instance Steps
CREATE TABLE workflow_instance_steps (
    id SERIAL PRIMARY KEY,
    workflow_instance_id INTEGER NOT NULL REFERENCES workflow_instances(id),
    workflow_step_id INTEGER NOT NULL REFERENCES workflow_steps(id),
    status VARCHAR(20) NOT NULL,
    event_message TEXT,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE
);
