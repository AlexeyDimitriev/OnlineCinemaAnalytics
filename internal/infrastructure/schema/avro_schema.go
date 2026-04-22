package schema

const MovieEventAvroSchema = `
{
  "type": "record",
  "name": "MovieEvent",
  "namespace": "com.movie.analytics",
  "doc": "Event of user action in online cinema",
  "fields": [
    {
      "name": "event_id",
      "type": "string",
      "doc": "Уникальный идентификатор события (UUID)"
    },
    {
      "name": "user_id",
      "type": "string",
      "doc": "Идентификатор пользователя"
    },
    {
      "name": "movie_id",
      "type": "string",
      "doc": "Идентификатор фильма"
    },
    {
      "name": "event_type",
      "type": {
        "type": "enum",
        "name": "EventType",
        "symbols": [
          "VIEW_STARTED",
          "VIEW_FINISHED",
          "VIEW_PAUSED",
          "VIEW_RESUMED",
          "LIKED",
          "SEARCHED"
        ],
        "doc": "Type of event: started, finished, paused or resumed view, like or searched"
      },
      "doc": "Тип произошедшего события"
    },
    {
      "name": "timestamp",
      "type": {
        "type": "long",
        "logicalType": "timestamp-millis"
      },
      "doc": "Время события в миллисекундах (UTC)"
    },
    {
      "name": "device_type",
      "type": {
        "type": "enum",
        "name": "DeviceType",
        "symbols": [
          "MOBILE",
          "DESKTOP",
          "TV",
          "TABLET"
        ],
        "doc": "Type of device"
      },
      "doc": "Device used for action"
    },
    {
      "name": "session_id",
      "type": "string",
      "doc": "Идентификатор сессии"
    },
    {
      "name": "progress_seconds",
      "type": "int",
      "doc": "Прогресс просмотра в секундах"
    }
  ]
}
`

// TopicName for Kafka events
const TopicName = "movie-events"

// Subject for Schema Registry in Kafka
const SchemaSubject = "movie-events-value"
