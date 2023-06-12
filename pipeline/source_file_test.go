package pipeline

import (
	"bufio"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/seal-io/terraform-provider-byteset/utils/testx"
)

func TestSource_srcFile_split(t *testing.T) {
	f, err := testx.File("testdata/complex.sql")
	if err != nil {
		panic(err)
	}

	defer func() { _ = f.Close() }()

	var actual []string
	ss := bufio.NewScanner(f)
	ss.Split(split)

	for ss.Scan() {
		actual = append(actual, ss.Text())
	}

	expected := []string{
		`-- Comment 1`,

		`--`,
		`-- Comment 2`,
		`--`,

		`-- /// Comment 3`,

		`-- /* Comment 4 */`,

		`/* -- Comment 5 */`,

		`/* Comment 6 */`,

		`/* Comment 7 */;`,

		`/*
    Comment 8
 */;`,

		`/*
    Comment 9
 */ ;`,

		`/*
    Comment 10;
*/`,

		`;`,

		`;;`,

		`;;;`,

		`/*!40014 SET @OLD_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS = 0 */;`,

		`DROP TABLE IF EXISTS test;`,

		`CREATE TABLE test
(
    id  INTEGER PRIMARY KEY AUTO_INCREMENT,
    val REAL
);`,

		`INSERT INTO test (val)
VALUES ('Test 1');`,

		`INSERT INTO test (val) VALUES ('Test 2');`,
	}

	assert.Equal(t, expected, actual)
}
