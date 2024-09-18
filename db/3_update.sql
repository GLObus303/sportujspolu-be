CREATE TABLE levels (
    id SERIAL PRIMARY KEY,
    value VARCHAR(20) NOT NULL,
    label VARCHAR(50) NOT NULL
);

INSERT INTO levels (value, label)
VALUES
    ('beginner', 'Začátečník'),
    ('advanced', 'Pokročilý'),
    ('expert', 'Expert'),
    ('any', 'Pro každého');
