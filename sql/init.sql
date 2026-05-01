CREATE EXTENSION IF NOT EXISTS "pgcrypto";


DROP TABLE IF EXISTS booking;
DROP TABLE IF EXISTS schedule_day_of_week;
DROP TABLE IF EXISTS schedule;
DROP TABLE IF EXISTS slot;
DROP TABLE IF EXISTS room;
DROP TABLE IF EXISTS app_user;

CREATE TABLE app_user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    user_role VARCHAR(10) NOT NULL,
    created_at DATE NOT NULL
);

-- INSERT INTO app_user(email, user_role, created_at) VALUES ("user", "user", NOW()), ("admin", "admin", NOW());



CREATE TABLE room (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL UNIQUE,
    description VARCHAR,
    capacity INTEGER NOT NULL,
    created_at DATE NOT NULL
);



CREATE TABLE schedule (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID REFERENCES room(id),
    created_at DATE NOT NULL
);



CREATE TABLE schedule_day_of_week (
    id SERIAL PRIMARY KEY,
    schedule_id UUID REFERENCES schedule(id),
    day_of_week INT CHECK(day_of_week >= 0 AND day_of_week <= 6)
);



CREATE TABLE slot (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID REFERENCES room(id),
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL
);



CREATE TABLE booking (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slot_id UUID REFERENCES slot(id),
    user_id UUID REFERENCES app_user(id),
    status VARCHAR(255)
);