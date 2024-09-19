CREATE TABLE IF NOT EXISTS system (
	id         serial primary key,
	device_id       integer not null,
	temp_gpu	real not null,
	temp_cpu 	real not null,
	created_at 	timestamp not null default now()
);

CREATE TABLE IF NOT EXISTS devices (
	id serial primary key,
	name varchar(255) unique not null,
	created_at 	timestamp not null default now()
);
