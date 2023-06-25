CREATE TABLE `events` (
	`id` int PRIMARY KEY AUTO_INCREMENT,
	`name` varchar(30) NOT NULL,
	`sport` varchar(20) NOT NULL
);

CREATE TABLE `users` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `email` varchar(30) NOT NULL UNIQUE, 
  `name` varchar(30) NOT NULL,
  `password` varchar(30) NOT NULL
);

INSERT INTO `events` (name, sport) VALUES
  ('Behani s libuskou', 'beh'),
  ('Patecni pinec', 'pingpong'),
  ('Zacatecnicky tenis', 'tenis');
