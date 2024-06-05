CREATE TABLE IF NOT EXISTS users (
  username VARCHAR PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE TABLE IF NOT EXISTS balances (
  username VARCHAR PRIMARY KEY,
  amount NUMERIC NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE TABLE IF NOT EXISTS histories (
  id CHAR(36) PRIMARY KEY,
  username VARCHAR NOT NULL,
  target_username VARCHAR NOT NULL,
  amount NUMERIC NOT NULL,
  "type" SMALLINT NOT NULL,
  notes VARCHAR NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX histories_username_created_at_desc_idx ON histories (username, created_at DESC);

CREATE TABLE IF NOT EXISTS history_summaries (
  id VARCHAR PRIMARY KEY,
  username VARCHAR,
  target_username VARCHAR,
  amount NUMERIC NOT NULL,
  "type" SMALLINT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX history_summaries_username_amount_desc_type ON history_summaries (username, amount DESC, type);
