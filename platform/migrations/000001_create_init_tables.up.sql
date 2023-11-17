-- ./platform/migrations/000001_create_init_tables.up.sql

-- Add UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set timezone
-- For more information, please visit:
-- https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
SET TIMEZONE="Europe/Moscow";

-- Create customers table
CREATE TABLE customers (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    name VARCHAR (255) NOT NULL,
    balance INT NOT NULL
);


-- Create campaigns table
CREATE TABLE campaigns (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    customer_id UUID REFERENCES customers(id),
    creative VARCHAR (255) NOT NULL,
    views INT NOT NULL DEFAULT 0
);

-- Create countries table
CREATE TABLE countries (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    NAME VARCHAR (255) NOT NULL
);

-- Create devices table
CREATE TABLE devices (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    NAME VARCHAR (255) NOT NULL
);


-- Create campaign countries table
CREATE TABLE campaign_countries (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    campaign_id UUID REFERENCES campaigns(id),
	country_id UUID REFERENCES countries(id)
);

-- Create campaign devices table
CREATE TABLE campaign_devices (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    campaign_id UUID REFERENCES campaigns(id),
	device_id UUID REFERENCES devices(id)
);



