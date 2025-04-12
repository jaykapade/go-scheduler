-- +goose Up
CREATE TABLE schedules (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    greeting_id UUID NOT NULL,
    start_date TIMESTAMP NOT NULL,
    frequency_type TEXT NOT NULL,
    frequency_interval INTEGER NOT NULL,
    time_of_day TIME NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_sent_at TIMESTAMP,
    latest_error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (greeting_id) REFERENCES greetings(id)
);

-- +goose Down
DROP TABLE schedules;
