

-- Comment 1

--
-- Comment 2
--

-- /// Comment 3

-- /* Comment 4 */

/* -- Comment 5 */

/* Comment 6 */

/* Comment 7 */;

/*
    Comment 8
 */;

/*
    Comment 9
 */ ;

/*
    Comment 10;
*/

;
;;
;;;

/*!40014 SET @OLD_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS = 0 */;

DROP TABLE IF EXISTS test;
CREATE TABLE test
(
    id  INTEGER PRIMARY KEY AUTO_INCREMENT,
    val REAL
);

INSERT INTO test (val)
VALUES ('Test 1');

INSERT INTO test (val) VALUES ('Test 2');
