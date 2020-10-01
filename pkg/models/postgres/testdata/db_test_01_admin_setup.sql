-- must run as postgres user

-- database
CREATE DATABASE groupics_test;

-- user
CREATE USER groupics_test WITH ENCRYPTED PASSWORD 'pass';

GRANT ALL PRIVILEGES ON DATABASE groupics_test TO groupics_test;
