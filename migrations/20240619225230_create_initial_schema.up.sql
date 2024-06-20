START TRANSACTION;

CREATE TABLE users
(
    id            INT AUTO_INCREMENT PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password      VARCHAR(255) NOT NULL,
    name          VARCHAR(255),
    gender        VARCHAR(10),
    date_of_birth DATE,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    location      VARCHAR(255) NOT NULL
);

CREATE TABLE swipes
(
    id           INT AUTO_INCREMENT PRIMARY KEY,
    user_id      INT  NOT NULL,
    candidate_id INT  NOT NULL,
    likes        BOOL NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY idx_swipes_swiper_swiped (user_id, candidate_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (candidate_id) REFERENCES users (id)
);

CREATE TABLE sessions
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    user_id    INT          NOT NULL,
    token      VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP    NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE user_preferences
(
    id              INT AUTO_INCREMENT PRIMARY KEY,
    user_id         INT NOT NULL,
    wants_children  BOOLEAN,
    enjoys_travel   BOOLEAN,
    education_level VARCHAR(50),
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    min_age         INT,
    max_age         INT,
    genders         VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

COMMIT;