CREATE TABLE system (
	temp_gpu	real not null,
	temp_cpu 	real not null,
	created_at 	timestamp not null default now()
);
