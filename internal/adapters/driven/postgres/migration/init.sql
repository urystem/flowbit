CREATE TABLE
  IF NOT EXISTS exchange_averages (
    source VARCHAR(128) NOT NULL CHECK (char_length(trim(source)) > 0),
    symbol VARCHAR(128) NOT NULL CHECK (char_length(trim(symbol)) > 0),
    count INT NOT NULL CHECK (count > 0),
    average_price DOUBLE PRECISION NOT NULL CHECK (average_price > 0),
    min_price DOUBLE PRECISION NOT NULL CHECK (min_price > 0),
    max_price DOUBLE PRECISION NOT NULL CHECK (max_price > 0),
    at_time TIMESTAMPTZ (3) NOT NULL
  );

CREATE TABLE
  IF NOT EXISTS exchange_backup (
    source VARCHAR(128) NOT NULL CHECK (char_length(trim(source)) > 0),
    symbol VARCHAR(128) NOT NULL CHECK (char_length(trim(symbol)) > 0),
    price DOUBLE PRECISION NOT NULL CHECK (price > 0),
    time_stamp BIGINT NOT NULL,
    CONSTRAINT exchange_unique UNIQUE (source, symbol, time_stamp)
  );