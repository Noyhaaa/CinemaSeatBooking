DROP TABLE IF EXISTS movies;
CREATE TABLE movies (
  id         INT AUTO_INCREMENT NOT NULL,
  title      VARCHAR(128) NOT NULL,
  realisator     VARCHAR(128) NOT NULL,
  leftticket     SMALLINT,
  maxticket      SMALLINT,
  price          DECIMAL(5,2),
  room           SMALLINT, 
  PRIMARY KEY (`id`)
);

INSERT INTO movies
  (title, realisator, leftticket, maxticket, price, room)
VALUES
  ('Elemental', 'Peter Sohn', 259, 259, 11.9, 1),
  ('Asteroid City', 'Wes Anderson', 457, 405, 11.99, 2),
  ('Spiderman Across the spiderverse', 'Joaquim Dos Santos', 0, 525, 11.99, 3),
  ('Indiana Jones', 'George Lucas', 179, 179, 11.99, 4);