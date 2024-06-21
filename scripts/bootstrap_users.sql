-- Insert users
INSERT INTO users (name, email, password, location, gender, date_of_birth) VALUES
('Alice', 'alice@example.com', '$2a$04$.n3FbKaXHF.sZAQTJIModOsmA6J0trqBy0vUOcC4ff6DZnPqoESqS', '37.7749,-122.4194', 'Female', '1990-01-01'),
('Bob', 'bob@example.com', '$2a$04$.n3FbKaXHF.sZAQTJIModOsmA6J0trqBy0vUOcC4ff6DZnPqoESqS', '34.0522,-118.2437', 'Male', '1988-05-15'),
('Charlie', 'charlie@example.com', '$2a$04$.n3FbKaXHF.sZAQTJIModOsmA6J0trqBy0vUOcC4ff6DZnPqoESqS', '40.7128,-74.0060', 'Male', '1992-09-23'),
('David', 'david@example.com', '$2a$04$.n3FbKaXHF.sZAQTJIModOsmA6J0trqBy0vUOcC4ff6DZnPqoESqS', '41.8781,-87.6298', 'Male', '1985-12-02'),
('Eve', 'eve@example.com', '$2a$04$.n3FbKaXHF.sZAQTJIModOsmA6J0trqBy0vUOcC4ff6DZnPqoESqS', '34.0522,-118.2437', 'Female', '1991-04-18'),
('Frank', 'frank@example.com', '$2a$04$.n3FbKaXHF.sZAQTJIModOsmA6J0trqBy0vUOcC4ff6DZnPqoESqS', '29.7604,-95.3698', 'Male', '1987-07-30'),
('Grace', 'grace@example.com', '$2a$04$.n3FbKaXHF.sZAQTJIModOsmA6J0trqBy0vUOcC4ff6DZnPqoESqS', '39.7392,-104.9903', 'Female', '1993-03-12'),
('Hank', 'hank@example.com', '$2a$04$.n3FbKaXHF.sZAQTJIModOsmA6J0trqBy0vUOcC4ff6DZnPqoESqS', '47.6062,-122.3321', 'Male', '1990-08-24'),
('Ivy', 'ivy@example.com', '$2a$04$.n3FbKaXHF.sZAQTJIModOsmA6J0trqBy0vUOcC4ff6DZnPqoESqS', '25.7617,-80.1918', 'Female', '1989-11-11'),
('Jack', 'jack@example.com', '$2a$04$.n3FbKaXHF.sZAQTJIModOsmA6J0trqBy0vUOcC4ff6DZnPqoESqS', '32.7767,-96.7970', 'Male', '1986-02-20');

-- Insert user preferences
INSERT INTO user_preferences (user_id, wants_children, enjoys_travel, education_level, min_age, max_age, genders)
VALUES (1, TRUE, TRUE, 'BSCH', 25, 35, 'Male,Female');
INSERT INTO user_preferences (user_id, wants_children, enjoys_travel, education_level, min_age, max_age, genders)
VALUES (2, FALSE, TRUE, 'MSCH', 20, 30, 'Female');
INSERT INTO user_preferences (user_id, wants_children, enjoys_travel, education_level, min_age, max_age, genders)
VALUES (3, TRUE, FALSE, 'HS', 30, 40, 'Male');
INSERT INTO user_preferences (user_id, wants_children, enjoys_travel, education_level, min_age, max_age, genders)
VALUES (4, FALSE, FALSE, 'PHD', 25, 35, 'Female');
INSERT INTO user_preferences (user_id, wants_children, enjoys_travel, education_level, min_age, max_age, genders)
VALUES (5, TRUE, TRUE, 'ASC', 22, 32, 'Male,Female');
INSERT INTO user_preferences (user_id, wants_children, enjoys_travel, education_level, min_age, max_age, genders)
VALUES (6, FALSE, TRUE, 'BSCH', 28, 38, 'Female');
INSERT INTO user_preferences (user_id, wants_children, enjoys_travel, education_level, min_age, max_age, genders)
VALUES (7, TRUE, TRUE, 'MSCH', 26, 36, 'Male');
INSERT INTO user_preferences (user_id, wants_children, enjoys_travel, education_level, min_age, max_age, genders)
VALUES (8, FALSE, FALSE, 'HS', 23, 33, 'Female');
INSERT INTO user_preferences (user_id, wants_children, enjoys_travel, education_level, min_age, max_age, genders)
VALUES (9, TRUE, TRUE, 'BSCH', 27, 37, 'Male,Female');
INSERT INTO user_preferences (user_id, wants_children, enjoys_travel, education_level, min_age, max_age, genders)
VALUES (10, FALSE, TRUE, 'PHD', 24, 34, 'Male,Female');

-- Insert swipes
INSERT INTO swipes (user_id, candidate_id, likes)
VALUES (1, 2, TRUE),
       (1, 3, TRUE),
       (1, 4, FALSE),
       (2, 1, TRUE),
       (2, 3, FALSE),
       (2, 5, TRUE),
       (3, 1, FALSE),
       (3, 2, TRUE),
       (3, 6, TRUE),
       (4, 1, TRUE),
       (4, 5, FALSE),
       (4, 6, TRUE),
       (5, 1, TRUE),
       (5, 2, TRUE),
       (5, 3, FALSE),
       (6, 1, TRUE),
       (6, 4, FALSE),
       (6, 5, TRUE),
       (7, 2, TRUE),
       (7, 3, TRUE),
       (7, 8, FALSE),
       (8, 1, TRUE),
       (8, 2, FALSE),
       (8, 7, TRUE),
       (9, 3, TRUE),
       (9, 4, FALSE),
       (9, 5, TRUE),
       (10, 1, TRUE),
       (10, 6, TRUE),
       (10, 8, FALSE);