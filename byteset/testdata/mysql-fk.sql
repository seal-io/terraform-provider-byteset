/* Setting */
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS = 0 */;

-- teams table
DROP TABLE IF EXISTS teams;
CREATE TABLE teams
(
    id   INTEGER PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100)
);

-- members table
DROP TABLE IF EXISTS members;
CREATE TABLE members
(
    id         INTEGER PRIMARY KEY AUTO_INCREMENT,
    last_name  VARCHAR(100) NOT NULL,
    first_name VARCHAR(100),
    team_id    INTEGER,
    CONSTRAINT fk_teams FOREIGN KEY (team_id) REFERENCES teams (id)
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
/*!40014 SET FOREIGN_KEY_CHECKS = @OLD_FOREIGN_KEY_CHECKS */;