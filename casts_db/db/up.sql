CREATE TABLE professions (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
CREATE TABLE casts (
    movie_id INT NOT NULL,
    person_id INT NOT NULL,
    profession_id INT REFERENCES professions(id) ON DELETE SET NULL ON UPDATE CASCADE,
    PRIMARY KEY(movie_id,person_id,profession_id)
);

GRANT SELECT ON casts TO casts_service;
GRANT SELECT ON professions TO casts_service;


