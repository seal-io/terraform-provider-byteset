/* Setting */
ALTER TABLE IF EXISTS members
    DROP CONSTRAINT fk_teams;

-- teams table
DROP TABLE IF EXISTS teams;
CREATE TABLE teams
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100)
);

-- members table
DROP TABLE IF EXISTS members;
CREATE TABLE members
(
    id         SERIAL PRIMARY KEY,
    last_name  VARCHAR(100) NOT NULL,
    first_name VARCHAR(100),
    team_id    INTEGER
);


-- members data
INSERT INTO members (id, last_name, first_name, team_id)
VALUES (1, 'Lucy', 'Li', 1);
INSERT INTO members (id, last_name, first_name, team_id)
VALUES (2, 'Lily', 'Zhang', 1);
INSERT INTO members (id, last_name, first_name, team_id)
VALUES (3, 'Stephen', 'Chen', 2);
INSERT INTO members (id, last_name, first_name, team_id)
VALUES (4, 'Frank', NULL, 2);

-- teams data
INSERT INTO teams (id, name)
VALUES (1, 'Finance');
INSERT INTO teams (id, name)
VALUES (2, 'Development');

/* Setting */
ALTER TABLE members
    ADD CONSTRAINT fk_teams FOREIGN KEY (team_id) REFERENCES teams (id) ON DELETE RESTRICT;
