CREATE TABLE "users" (
    "id" serial PRIMARY KEY,
    "first_name" character varying NOT NULL,
    "last_name" character varying NOT NULL,
    "middle_name" character varying,
    "email" character varying NOT NULL UNIQUE,
    "phone_number" varchar UNIQUE,
    "birthday" timestamp with time zone,
    "password" character varying NOT NULL,
    "created_at" timestamp with time zone DEFAULT current_timestamp,
    "updated_at" timestamp with time zone DEFAULT current_timestamp
);