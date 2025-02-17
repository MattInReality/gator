-- +goose Up
CREATE TABLE feed_follows (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  user_id UUID NOT NULL,
  feed_id UUID NOT NULL,
  UNIQUE (user_id, feed_id),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_feed FOREIGN KEY (feed_id)  REFERENCES feeds(id) ON DELETE CASCADE
);
-- +goose Down
DROP TABLE feed_follows;
