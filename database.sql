CREATE TABLE IF NOT EXISTS users (
  id CHAR(36) PRIMARY KEY,
  username VARCHAR,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE UNIQUE INDEX users_username_unq ON users (username); 

CREATE TABLE IF NOT EXISTS balances (
  user_id CHAR(36) PRIMARY KEY,
  amount NUMERIC NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE TABLE IF NOT EXISTS histories (
  id CHAR(36) PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  target_user_id CHAR(36) NOT NULL,
  amount NUMERIC NOT NULL,
  "type" SMALLINT NOT NULL,
  notes VARCHAR NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX histories_user_id_created_at_desc_idx ON histories (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS history_summaries (
  id VARCHAR PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  target_user_id CHAR(36),
  amount NUMERIC NOT NULL,
  "type" SMALLINT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX history_summaries_username_amount_desc_type ON history_summaries (user_id, amount DESC, type);
