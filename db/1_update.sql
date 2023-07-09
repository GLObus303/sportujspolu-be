CREATE TABLE `event_levels` (`id` int PRIMARY KEY, `name` varchar(20));
INSERT INTO `event_levels` (`id`, `name`) VALUES (1, 'zacatecnik'), (2, 'pokrocily'), (3, 'expert');


CREATE TABLE `sports` (`id` int PRIMARY KEY, `name` varchar(20));
INSERT INTO `sports` (`id`, `name`) VALUES (1, 'basketbal'), (2, 'florbal'), (3, 'fotbal');


ALTER TABLE `events`
ADD COLUMN `date` DATE,
ADD COLUMN `location` varchar(50) NOT NULL,
ADD COLUMN `price` decimal(10,2) NOT NULL,
ADD COLUMN `description` TEXT NOT NULL,
ADD COLUMN `level` varchar(30) NOT NULL;

CREATE TABLE `users` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(30) NOT NULL,
  `email` varchar(30) NOT NULL,
  `password` varchar(30) NOT NULL,
  `rating` int NOT NULL
);

INSERT INTO `events` (name, sport, date, location, price, description, level) VALUES
  ('Basketball Match at Park', 'Basketball', '2023-07-10', 'Central Park', 0.00, 'Looking for players for a friendly basketball match at the park. All skill levels are welcome!', 'Any'),
  ('Football Training Session', 'Football', '2023-07-15', 'City Stadium', 5.00, 'Organizing a football training session to improve skills and have fun. Intermediate to advanced level players preferred.', 'Intermediate'),
  ('Tennis Doubles Tournament', 'Tennis', '2023-07-20', 'Community Tennis Club', 10.00, 'Join our doubles tournament at the local tennis club. Players with previous tennis experience are encouraged to participate.', 'Intermediate'),
  ('Morning Jogging Group', 'Running', '2023-07-12', 'City Park', 0.00, 'Starting a morning jogging group for fitness enthusiasts. We''ll meet at the park and run together.', 'Any'),
  ('Swimming Lessons for Beginners', 'Swimming', '2023-07-18', 'Community Pool', 15.00, 'Offering swimming lessons for beginners. Get started with the basics of swimming in a supportive environment.', 'Beginner'),
  ('Cycling to the Countryside', 'Cycling', '2023-07-25', 'Meeting Point: Bike Shop', 0.00, 'Embark on a cycling adventure to the beautiful countryside. Intermediate level cyclists preferred.', 'Intermediate'),
  ('Badminton Open Play', 'Badminton', '2023-07-13', 'Local Sports Hall', 3.50, 'Open play session for badminton enthusiasts. Join us for some friendly matches and improve your skills.', 'Any'),
  ('Table Tennis Tournament', 'Table Tennis', '2023-07-22', 'Community Center', 8.00, 'Compete in our table tennis tournament and showcase your skills. Intermediate to advanced level players are invited.', 'Intermediate'),
  ('Group Hiking Trip', 'Hiking', '2023-07-17', 'Mountain Range', 0.00, 'Planning a group hiking trip to explore the scenic trails. All levels of hikers are welcome!', 'Any'),
  ('Golf Practice Session', 'Golf', '2023-07-30', 'Local Golf Course', 20.00, 'Join our golf practice session to improve your swing and technique. All skill levels are welcome!', 'Any');
