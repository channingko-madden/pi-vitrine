CREATE TABLE IF NOT EXISTS system (
	id         serial primary key,
	device_id       integer not null,
	temp_gpu	real not null default 0,
	temp_cpu 	real not null default 0,
	created_at 	timestamp not null default now()
);

CREATE TABLE IF NOT EXISTS devices (
	id serial primary key,
	name varchar(255) unique not null,
	location varchar(255) not null default '',
	created_at 	timestamp not null default now()
);

CREATE TABLE IF NOT EXISTS indoor_climate (
	id serial primary key,
	device_id       integer not null,
	air_temp	real not null default 0, -- Kelvin
	pressure 	real not null default 0, -- Pascal
	relative_humidity real not null default 0, -- 0 - 100%
	created_at 	timestamp not null default now()
);
