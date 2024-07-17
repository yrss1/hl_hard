
CREATE TABLE IF NOT EXISTS users (
                                     created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     id SERIAL PRIMARY KEY,
                                     full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(50) NOT NULL
    );

CREATE TABLE IF NOT EXISTS projects (
                                        created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        id SERIAL PRIMARY KEY,
                                        title VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    start_date DATE NOT NULL,
    end_date DATE,
    manager_id INT NOT NULL,
    FOREIGN KEY (manager_id) REFERENCES users(id) ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS tasks (
                                     created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     id SERIAL PRIMARY KEY,
                                     title VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    priority VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    assignee_id INT,
    project_id INT NOT NULL,
    completed_at DATE,
    FOREIGN KEY (assignee_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
    );

INSERT INTO users (full_name, email, role) VALUES
                                               ('Alenov Abay', 'onelvay@google.com', 'developer'),
                                               ('Rick Sanchez', 'WubbaLubbadubdub@c137.com', 'scientist'),
                                               ('Morty Smith', 'theonetruemorty@c137.com', 'student'),
                                               ('Who u are', 'whouare@whoami.com', 'developer');
INSERT INTO projects (title, description, start_date, end_date, manager_id) VALUES
                                                                               ('Alpha', 'Save Morty', '2023-01-01', '2023-06-30', 2),
                                                                               ('Beta', 'Bla-bla-bla', '2023-02-01', '2023-07-31', 2),
                                                                               ('Gamma', 'Third project description', '2023-03-01', NULL, 1);
INSERT INTO tasks (title, description, priority, status, assignee_id, project_id, completed_at) VALUES
                                                                                                    ('Design Homepage', 'Create a responsive homepage design', 'Medium', 'Active', 4, 1, NULL),
                                                                                                    ('Implement Login', 'Develop the login functionality', 'Low', 'Active', 2, 1, NULL),
                                                                                                    ('Database Schema', 'Define the database schema', 'Medium', 'Done', 3, 1, '2023-04-15 10:30:00'),
                                                                                                    ('Where', 'Find out where Morty''s been taken', 'High', 'Done', 2, 2, '2023-01-01'),
                                                                                                    ('Rescue', 'Rescue', 'High', 'Active', 2, 2, NULL);
