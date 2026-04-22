CREATE TABLE IF NOT EXISTS movie_events_kafka (
    event_id String,
    user_id String,
    movie_id String,
    event_type String,
    timestamp DateTime64(3, 'UTC'),
    device_type String,
    session_id String,
    progress_seconds Int32
) ENGINE = Kafka()
SETTINGS
    kafka_broker_list = 'kafka:29092',
    kafka_topic_list = 'movie-events',
    kafka_group_name = 'clickhouse_consumer_group',
    kafka_format = 'JSONEachRow',
    kafka_skip_broken_messages = 1;

CREATE TABLE IF NOT EXISTS movie_events (
    event_id String,
    user_id String,
    movie_id String,
    event_type Enum8(
        'VIEW_STARTED' = 1,
        'VIEW_FINISHED' = 2,
        'VIEW_PAUSED' = 3,
        'VIEW_RESUMED' = 4,
        'LIKED' = 5,
        'SEARCHED' = 6
    ),
    timestamp DateTime64(3, 'UTC'),
    device_type Enum8(
        'MOBILE' = 1,
        'DESKTOP' = 2,
        'TV' = 3,
        'TABLET' = 4
    ),
    session_id String,
    progress_seconds Int32,
    
    INDEX idx_event_type event_type TYPE set(0) GRANULARITY 4,
    INDEX idx_user_id user_id TYPE bloom_filter GRANULARITY 4
) ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (user_id, toStartOfDay(timestamp), timestamp)
SETTINGS
    index_granularity = 8192,
    allow_nullable_key = 1;

CREATE MATERIALIZED VIEW IF NOT EXISTS movie_events_mv
TO movie_events
AS SELECT
    event_id,
    user_id,
    movie_id,
    event_type,
    timestamp,
    device_type,
    session_id,
    progress_seconds
FROM movie_events_kafka;

SYSTEM SYNC REPLICA movie_events;
