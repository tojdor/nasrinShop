CREATE TABLE materials(
    id SERIAL PRIMARY KEY,
    categorie_id INT NOT NULL REFERENCES categories(id),
    price NUMERIC(10, 2) NOT NULL,
    image_url TEXT NOT NULL,
);