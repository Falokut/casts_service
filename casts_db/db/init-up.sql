CREATE ROLE casts_service WITH
    LOGIN
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXzMGszsnnHexVwOU=';

CREATE TABLE professions (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);
CREATE TABLE casts (
    movie_id INT NOT NULL,
    actor_id INT NOT NULL,
    profession_id INT REFERENCES professions(id) ON DELETE SET NULL ON UPDATE CASCADE,
    PRIMARY KEY(movie_id,actor_id,profession_id)
);

GRANT SELECT ON casts TO casts_service;
GRANT SELECT ON professions TO casts_service;