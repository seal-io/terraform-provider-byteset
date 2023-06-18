package byteset

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/seal-io/terraform-provider-byteset/utils/strx"
	"github.com/seal-io/terraform-provider-byteset/utils/testx"
)

func TestAccResourcePipeline_raw_to_mysql(t *testing.T) {
	// Start Database.
	var (
		database = "byteset"
		password = strx.String(10)
	)

	ctx := context.TODO()
	c := dockerContainer{
		Name:  "mysql",
		Image: "mysql:8",
		Env: []string{
			"MYSQL_DATABASE=" + database,
			"MYSQL_ROOT_PASSWORD=" + password,
		},
		Port: []string{
			"3306:3306",
		},
	}

	err := c.Start(t, ctx)
	if err != nil {
		t.Fatalf("failed to start MySQL container: %v", err)
	}

	defer func() { _ = c.Stop(t, ctx) }()

	// Test pipeline.
	var (
		resourceName = "byteset_pipeline.test"

		basicSrc = strings.ReplaceAll(`raw://
-- company table
DROP TABLE IF EXISTS company;
CREATE TABLE company
(
    id      INTEGER PRIMARY KEY AUTO_INCREMENT,
    name    TEXT NOT NULL,
    age     INT  NOT NULL,
    address CHAR(50),
    salary  NUMERIC
);

-- company data
INSERT INTO company (name, age, address, salary) VALUES ('Paul', 32, 'California', 20000.00);
INSERT INTO company (name, age, address, salary) VALUES ('Allen', 25, 'Texas', 15000.00);
INSERT INTO company (name, age, address, salary) VALUES ('Teddy', 23, 'Norway', 20000.00);
INSERT INTO company (name, age, address, salary) VALUES ('Mark', 25, 'Rich-Mond ', 65000.00);
INSERT INTO company (name, age, address, salary) VALUES ('David', 27, 'Texas', 85000.00);
INSERT INTO company (name, age, address, salary) VALUES ('Kim', 22, 'South-Hall', 45000.00);
INSERT INTO company (name, age, address, salary) VALUES ('James', 24, 'Houston', 10000.00);`, "\n", "\\n")
		basicDst = fmt.Sprintf("mysql://root:%s@tcp(127.0.0.1:3306)/%s", password, database)
	)

	resource.Test(t, resource.TestCase{
		IDRefreshName:            resourceName,
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigOfSourceFile(basicSrc, basicDst, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName,
						"source.address",
						strings.ReplaceAll(basicSrc, "\\n", "\n"),
					),
					resource.TestCheckResourceAttr(resourceName, "destination.address", basicDst),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "destination.cost"),
				),
			},
		},
	})
}

func TestAccResourcePipeline_raw_base64_to_mysql(t *testing.T) {
	// Start Database.
	var (
		database = "byteset"
		password = strx.String(10)
	)

	ctx := context.TODO()
	c := dockerContainer{
		Name:  "mysql",
		Image: "mysql:8",
		Env: []string{
			"MYSQL_DATABASE=" + database,
			"MYSQL_ROOT_PASSWORD=" + password,
		},
		Port: []string{
			"3306:3306",
		},
	}

	err := c.Start(t, ctx)
	if err != nil {
		t.Fatalf("failed to start MySQL container: %v", err)
	}

	defer func() { _ = c.Stop(t, ctx) }()

	// Test pipeline.
	var (
		resourceName = "byteset_pipeline.test"

		basicSrc = fmt.Sprintf("raw+base64://%s", strx.EncodeBase64(`
-- company table
DROP TABLE IF EXISTS company;
CREATE TABLE company
(
    id      INTEGER PRIMARY KEY AUTO_INCREMENT,
    name    TEXT NOT NULL,
    age     INT  NOT NULL,
    address CHAR(50),
    salary  NUMERIC
);

-- company data
INSERT INTO company (name, age, address, salary) VALUES ('Paul', 32, 'California', 20000.00);
INSERT INTO company (name, age, address, salary) VALUES ('Allen', 25, 'Texas', 15000.00);
INSERT INTO company (name, age, address, salary) VALUES ('Teddy', 23, 'Norway', 20000.00);
INSERT INTO company (name, age, address, salary) VALUES ('Mark', 25, 'Rich-Mond ', 65000.00);
INSERT INTO company (name, age, address, salary) VALUES ('David', 27, 'Texas', 85000.00);
INSERT INTO company (name, age, address, salary) VALUES ('Kim', 22, 'South-Hall', 45000.00);
INSERT INTO company (name, age, address, salary) VALUES ('James', 24, 'Houston', 10000.00);`))
		basicDst = fmt.Sprintf("mysql://root:%s@tcp(127.0.0.1:3306)/%s", password, database)
	)

	resource.Test(t, resource.TestCase{
		IDRefreshName:            resourceName,
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigOfSourceFile(basicSrc, basicDst, 10),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source.address", basicSrc),
					resource.TestCheckResourceAttr(resourceName, "destination.address", basicDst),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "destination.cost"),
				),
			},
		},
	})
}

func TestAccResourcePipeline_file_to_mysql(t *testing.T) {
	// Start Database.
	var (
		database = "byteset"
		password = strx.String(10)
	)

	ctx := context.TODO()
	c := dockerContainer{
		Name:  "mysql",
		Image: "mysql:8",
		Env: []string{
			"MYSQL_DATABASE=" + database,
			"MYSQL_ROOT_PASSWORD=" + password,
		},
		Port: []string{
			"3306:3306",
		},
	}

	err := c.Start(t, ctx)
	if err != nil {
		t.Fatalf("failed to start MySQL container: %v", err)
	}

	defer func() { _ = c.Stop(t, ctx) }()

	// Test pipeline.
	var (
		testdataPath = testx.AbsolutePath("testdata")
		resourceName = "byteset_pipeline.test"

		basicSrc = fmt.Sprintf("file://%s/mysql.sql", testdataPath)
		basicDst = fmt.Sprintf("mysql://root:%s@tcp(127.0.0.1:3306)/%s", password, database)

		fkSrc = fmt.Sprintf("file://%s/mysql-fk.sql", testdataPath)
		fkDst = fmt.Sprintf("mysql://root:%s@tcp(127.0.0.1)/%s", password, database)

		largeSrc = "https://raw.githubusercontent.com/seal-io/terraform-provider-byteset/main/byteset/testdata/mysql-lg.sql"
		largeDst = fmt.Sprintf("mysql://root:%s@tcp/", password)
	)

	resource.Test(t, resource.TestCase{
		IDRefreshName:            resourceName,
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigOfSourceFile(basicSrc, basicDst, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source.address", basicSrc),
					resource.TestCheckResourceAttr(resourceName, "destination.address", basicDst),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "destination.cost"),
				),
			},
			{
				Config: testConfigOfSourceFile(fkSrc, fkDst, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source.address", fkSrc),
					resource.TestCheckResourceAttr(resourceName, "destination.address", fkDst),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "destination.cost"),
				),
			},
			{
				Config: testConfigOfSourceFile(largeSrc, largeDst, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source.address", largeSrc),
					resource.TestCheckResourceAttr(resourceName, "destination.address", largeDst),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "destination.cost"),
				),
			},
		},
	})
}

func TestAccResourcePipeline_file_to_postgres(t *testing.T) {
	// Start Database.
	var (
		database = "byteset"
		password = strx.String(10)
	)

	ctx := context.TODO()
	c := dockerContainer{
		Name:  "postgres",
		Image: "postgres:13",
		Env: []string{
			"POSTGRES_DB=" + database,
			"POSTGRES_USER=root", // Rename superuser from postgres to root.
			"POSTGRES_PASSWORD=" + password,
		},
		Port: []string{
			"5432:5432",
		},
	}

	err := c.Start(t, ctx)
	if err != nil {
		t.Fatalf("failed to start Postgres container: %v", err)
	}

	defer func() { _ = c.Stop(t, ctx) }()

	// Test pipeline.
	var (
		testdataPath = testx.AbsolutePath("testdata")
		resourceName = "byteset_pipeline.test"

		basicSrc = fmt.Sprintf("file://%s/postgres.sql", testdataPath)
		basicDst = fmt.Sprintf("postgresql://root:%s@127.0.0.1:5432/%s?sslmode=disable", password, database)

		fkSrc = fmt.Sprintf("file://%s/postgres-fk.sql", testdataPath)
		fkDst = fmt.Sprintf("postgres://root:%s@127.0.0.1/%s?sslmode=disable", password, database)

		largeSrc = "https://raw.githubusercontent.com/seal-io/terraform-provider-byteset/main/byteset/testdata/postgres-lg.sql"
		largeDst = fmt.Sprintf("postgresql://root:%s@127.0.0.1/%s?sslmode=disable", password, database)
	)

	resource.Test(t, resource.TestCase{
		IDRefreshName:            resourceName,
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigOfSourceFile(basicSrc, basicDst, 10),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source.address", basicSrc),
					resource.TestCheckResourceAttr(resourceName, "destination.address", basicDst),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "destination.cost"),
				),
			},
			{
				Config: testConfigOfSourceFile(fkSrc, fkDst, 10),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source.address", fkSrc),
					resource.TestCheckResourceAttr(resourceName, "destination.address", fkDst),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "destination.cost"),
				),
			},
			{
				Config: testConfigOfSourceFile(largeSrc, largeDst, 10),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source.address", largeSrc),
					resource.TestCheckResourceAttr(resourceName, "destination.address", largeDst),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "destination.cost"),
				),
			},
		},
	})
}

func testConfigOfSourceFile(src, dst string, dstConnMax int) string {
	const tmpl = `
resource "byteset_pipeline" "test" {
  source = {
	address = "{{ .Src -}}"
  }
  destination = {
    address = "{{ .Dst }}"
    conn_max = {{ .DstConnMax }}
  }
}`

	return renderConfigTemplate(tmpl,
		"Src", src,
		"Dst", dst,
		"DstConnMax", dstConnMax)
}
