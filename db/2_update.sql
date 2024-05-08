CREATE TABLE email_requests (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    text TEXT NOT NULL,
    event_id varchar(12) NOT NULL,
    event_owner_id varchar(12) NOT NULL,
    requester_id varchar(12) NOT NULL,
    approved BOOLEAN DEFAULT NULL,
    approved_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE events ALTER COLUMN name TYPE VARCHAR(100);
