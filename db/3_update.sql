ALTER TABLE events ADD COLUMN public_id varchar(12),
ADD UNIQUE KEY idx_public_id (public_id);