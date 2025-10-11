CREATE TABLE IF NOT EXISTS exchange_averages(
    source VARCHAR(128)  NOT NULL,
    symbol VARCHAR(128)  NOT NULL,
    count INT NOT NULL,
    average_price  DOUBLE PRECISION NOT NULL,
    min_price   DOUBLE PRECISION NOT NULL, 
    max_price DOUBLE PRECISION NOT NULL,
    at_time BIGINT NOT NULL 
);

CREATE TABLE IF NOT EXISTS exchange_backup(
  source VARCHAR(128)  NOT NULL,
  symbol VARCHAR(128)  NOT NULL,
  price  DOUBLE PRECISION NOT NULL, 
  time_stamp BIGINT NOT NULL 
);

