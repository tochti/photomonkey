CREATE TABLE photos (
	id varchar(255) PRIMARY KEY,
	hash char(40) NOT NULL,
	caption TEXT,
	create_time TIMESTAMP DEFAULT current_timestamp
);
