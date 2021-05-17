# Create teams table
CREATE TABLE IF NOT EXISTS teams (
     id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
     name VARCHAR(100) NOT NULL,
     UNIQUE (name)
);

# Create people table
CREATE TABLE IF NOT EXISTS people (
      id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
      first_name VARCHAR(100) NOT NULL,
      last_name VARCHAR(100) NOT NULL,
      email VARCHAR(80) NOT NULL,
      team_id BIGINT NOT NULL,
      UNIQUE(email),
      FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
);

# Create turns table
CREATE TABLE IF NOT EXISTS turns (
     id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
     person_id BIGINT NOT NULL,
     team_id BIGINT NOT NULL,
     date DATE NOT NULL,
     created_at DATETIME DEFAULT NOW(),
     UNIQUE(team_id, date),
     FOREIGN KEY (person_id) REFERENCES people(id) ON DELETE CASCADE
);

