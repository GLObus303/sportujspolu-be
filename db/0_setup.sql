CREATE TABLE `events` (
	`id` int PRIMARY KEY AUTO_INCREMENT,
	`name` varchar(30) NOT NULL,
	`sport` varchar(20) NOT NULL
);

INSERT INTO `events` (name, sport) VALUES
  ('Behani s libuskou', 'beh'),
  ('Patecni pinec', 'pingpong'),
  ('Zacatecnicky tenis', 'tenis');
