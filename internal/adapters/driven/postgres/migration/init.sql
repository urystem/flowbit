-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--  SELECT uuid_generate_v4();

CREATE TABLE IF NOT EXISTS users (
  session_id UUID PRIMARY KEY /*UNIQUE*/,
  name VARCHAR(150),
  avatar_url TEXT--will be null if it is in ricky and morty
);
-------------------------------------------

CREATE TABLE IF NOT EXISTS posts (
  post_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_id UUID REFERENCES users(session_id) ON DELETE CASCADE NOT NULL,
  author VARCHAR(150)  NOT NULL,
  title VARCHAR(150)  NOT NULL,
  post_content TEXT,
  has_image BOOLEAN NOT NULL,
  archived BOOLEAN DEFAULT FALSE NOT NULL,
  post_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Быстрый доступ к только активным или архивным
CREATE INDEX idx_posts_archived ON posts(archived);

-- -- Комбинированный для главной (если часто WHERE archived=false ORDER BY post_time DESC)
-- CREATE INDEX idx_posts_main_page ON posts(archived, post_time DESC);

------------------------------------------
-- ) PARTITION BY LIST (archived);

-- -- Активные посты (на главной)
-- CREATE TABLE IF NOT EXISTS posts_active
--   PARTITION OF posts FOR VALUES IN (false);

-- -- Архивированные посты
-- CREATE TABLE IF NOT EXISTS posts_archived
--   PARTITION OF posts FOR VALUES IN (true);


-------------------------------------------------------

CREATE TABLE IF NOT EXISTS comments (
  comment_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  post_id INT REFERENCES posts(post_id) ON DELETE CASCADE NOT NULL,
  user_id UUID REFERENCES users(session_id) ON DELETE CASCADE NOT NULL,
  comment_content TEXT,
  parent_comment_id BIGINT REFERENCES comments(comment_id) ON DELETE CASCADE,--DEFAULT NULL
  has_image BOOLEAN,
  comment_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
