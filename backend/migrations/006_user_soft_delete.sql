-- Add soft delete support for users
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMPTZ;
