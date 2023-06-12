-- company table
DROP TABLE IF EXISTS company;
CREATE TABLE company
(
    id      SERIAL PRIMARY KEY,
    name    TEXT NOT NULL,
    age     INT  NOT NULL,
    address CHAR(50),
    salary  REAL
);


-- company data
INSERT INTO company (name, age, address, salary)
VALUES ('Paul', 32, 'California', 20000.00);
INSERT INTO company (name, age, address, salary)
VALUES ('Allen', 25, 'Texas', 15000.00);
INSERT INTO company (name, age, address, salary)
VALUES ('Teddy', 23, 'Norway', 20000.00);
INSERT INTO company (name, age, address, salary)
VALUES ('Mark', 25, 'Rich-Mond ', 65000.00);
INSERT INTO company (name, age, address, salary)
VALUES ('David', 27, 'Texas', 85000.00);
INSERT INTO company (name, age, address, salary)
VALUES ('Kim', 22, 'South-Hall', 45000.00);
INSERT INTO company (name, age, address, salary)
VALUES ('James', 24, 'Houston', 10000.00);
