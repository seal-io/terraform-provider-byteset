/* Setting */
PRAGMA foreign_keys = OFF;

-- teams table
DROP TABLE IF EXISTS teams;
CREATE TABLE teams
(
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR
);

-- members table
DROP TABLE IF EXISTS members;
CREATE TABLE members
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    last_name  VARCHAR NOT NULL,
    first_name VARCHAR,
    team_id    INTEGER,
    CONSTRAINT fk_teams FOREIGN KEY (team_id) REFERENCES teams (id)
);

-- members data
;;
;;
INSERT INTO members (id, last_name, first_name, team_id)
VALUES (1, 'Lucy', 'Li', 1);
INSERT INTO members (id, last_name, first_name, team_id)
VALUES (2, 'Lily', 'Zhang', 1);
INSERT INTO members (id, last_name, first_name, team_id)
VALUES (3, 'Stephen', 'Chen', 2);
INSERT INTO members (id, last_name, first_name, team_id)
VALUES (4, 'Frank', NULL, 2);

-- teams data
;;
;;
INSERT INTO teams (id, name)
VALUES (1, 'Finance');
INSERT INTO teams (id, name)
VALUES (2, 'Development');

/* Setting */
PRAGMA foreign_keys = ON;
