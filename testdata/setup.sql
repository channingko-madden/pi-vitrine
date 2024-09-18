CREATE TABLE IF NOT EXISTS system (
	id         serial primary key,
	device_id       integer not null,
	temp_gpu	real not null,
	temp_cpu 	real not null,
	created_at 	timestamp not null default now(),
);

CREATE TABLE IF NOT EXISTS devices (
	id serial primary key,
	name varchar(255) unique not null,
	mac_addr varchar(50) unique not null,
	hardware varchar(255),
	revision varchar(255), 
	serial varchar(255), 
	model varchar(255),
	created_at 	timestamp not null default now()
);
