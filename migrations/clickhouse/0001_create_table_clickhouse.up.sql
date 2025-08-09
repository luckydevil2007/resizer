
CREATE TABLE  IF NOT EXISTS events (image_id UInt32, timestamp DateTime,transform_type UInt32, user_id UInt32) 
ENGINE = MergeTree() ORDER BY (timestamp, image_id) PRIMARY KEY (timestamp, image_id)
