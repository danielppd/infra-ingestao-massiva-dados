CREATE TABLE IF NOT EXISTS sensor_readings (
    id          SERIAL PRIMARY KEY,
    device_id   VARCHAR(100)   NOT NULL,
    timestamp   TIMESTAMPTZ    NOT NULL,
    sensor_type VARCHAR(50)    NOT NULL,
    reading_type VARCHAR(20)   NOT NULL CHECK (reading_type IN ('analog', 'discrete')),
    value       NUMERIC(12, 4) NOT NULL,
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sensor_readings_device_id   ON sensor_readings(device_id);
CREATE INDEX IF NOT EXISTS idx_sensor_readings_sensor_type ON sensor_readings(sensor_type);
CREATE INDEX IF NOT EXISTS idx_sensor_readings_timestamp   ON sensor_readings(timestamp);