
DROP DATABASE IF EXISTS msds;
CREATE DATABASE msds;

DROP TABLE IF EXISTS MSDSCourseCatalog;

\c msds;

CREATE TABLE MSDSCourseCatalog (
    id SERIAL,
    cid VARCHAR(7) NOT NULL,
    cname VARCHAR(100),
    cprereq VARCHAR(100)
);
