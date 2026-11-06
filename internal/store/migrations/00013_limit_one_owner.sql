-- +goose Up
-- At most one row may have role = 'owner'.
CREATE UNIQUE INDEX users_single_owner_idx ON users ((role)) WHERE role = 'owner';

-- +goose Down
DROP INDEX users_single_owner_idx;
