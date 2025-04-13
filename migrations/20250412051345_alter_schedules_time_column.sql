-- +goose Up
ALTER TABLE schedules 
    RENAME COLUMN time_of_day TO scheduled_time;

ALTER TABLE schedules 
    ALTER COLUMN scheduled_time TYPE TIMESTAMP WITH TIME ZONE 
    USING (CURRENT_DATE + scheduled_time)::timestamp with time zone;

-- +goose Down
ALTER TABLE schedules 
    ALTER COLUMN scheduled_time TYPE TIME 
    USING scheduled_time::time;

ALTER TABLE schedules 
    RENAME COLUMN scheduled_time TO time_of_day; 