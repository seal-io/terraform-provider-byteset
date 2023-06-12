package byteset

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/seal-io/terraform-provider-byteset/utils/strx"
	"github.com/seal-io/terraform-provider-byteset/utils/testx"
)

func TestAccResourcePipeline_file_to_sqlite(t *testing.T) {
	// Test pipeline.
	var (
		testdataPath = testx.AbsolutePath("testdata")

		resourceName = "byteset_pipeline.test"

		basicSrc = fmt.Sprintf("file://%s/sqlite.sql", testdataPath)
		basicDst = "sqlite:///tmp/sqlite.db"

		fkSrc = fmt.Sprintf("file://%s/sqlite-fk.sql", testdataPath)
		fkDst = "sqlite:///tmp/sqlite.db?_pragma=foreign_keys(1)"
	)

	resource.Test(t, resource.TestCase{
		IDRefreshName:            resourceName,
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			// Basic.
			{
				Config: testConfigOfSourceFile(basicSrc, basicDst, 1, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source.address", basicSrc),
					resource.TestCheckResourceAttr(resourceName, "destination.address", basicDst),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max_open", "1"),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max_idle", "1"),
				),
			},
			// Foreign Key.
			{
				Config: testConfigOfSourceFile(fkSrc, fkDst, 1, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source.address", fkSrc),
					resource.TestCheckResourceAttr(resourceName, "destination.address", fkDst),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max_open", "1"),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max_idle", "1"),
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
				Config: testConfigOfSourceFile(basicSrc, basicDst, 5, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source.address", basicSrc),
					resource.TestCheckResourceAttr(resourceName, "destination.address", basicDst),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max_open", "5"),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max_idle", "5"),
				),
			},
			{
				Config: testConfigOfSourceFile(fkSrc, fkDst, 5, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source.address", fkSrc),
					resource.TestCheckResourceAttr(resourceName, "destination.address", fkDst),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max_open", "5"),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max_idle", "5"),
				),
			},
			{
				Config: testConfigOfSourceFile(largeSrc, largeDst, 5, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source.address", largeSrc),
					resource.TestCheckResourceAttr(resourceName, "destination.address", largeDst),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max_open", "5"),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max_idle", "5"),
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
				Config: testConfigOfSourceFile(basicSrc, basicDst, 5, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source.address", basicSrc),
					resource.TestCheckResourceAttr(resourceName, "destination.address", basicDst),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max_open", "5"),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max_idle", "5"),
				),
			},
			{
				Config: testConfigOfSourceFile(fkSrc, fkDst, 5, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source.address", fkSrc),
					resource.TestCheckResourceAttr(resourceName, "destination.address", fkDst),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max_open", "5"),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max_idle", "5"),
				),
			},
			{
				Config: testConfigOfSourceFile(largeSrc, largeDst, 5, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source.address", largeSrc),
					resource.TestCheckResourceAttr(resourceName, "destination.address", largeDst),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max_open", "5"),
					resource.TestCheckResourceAttr(resourceName, "destination.conn_max_idle", "5"),
				),
			},
		},
	})
}

func testConfigOfSourceFile(src, dst string, dstMaxOpen, dstMaxIdle int) string {
	const tmpl = `
resource "byteset_pipeline" "test" {
  source = {
	address = "{{ .Src }}"
  }
  destination = {
    address = "{{ .Dst }}"
    conn_max_open = {{ .DstMaxOpen }}
    conn_max_idle = {{ .DstMaxIdle }}
  }
}`

	return renderConfigTemplate(tmpl,
		"Src", src,
		"Dst", dst,
		"DstMaxOpen", dstMaxOpen,
		"DstMaxIdle", dstMaxIdle)
}
