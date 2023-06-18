package sqlx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type (
		input struct {
			drv string
			sql string
		}
		output struct {
			insert DMLInsert
			ok     bool
		}
	)

	tc := []struct {
		given    input
		expected output
	}{
		{
			given: input{
				drv: MySQLDialect,
				sql: "INSERT INTO `city` (`ID`, `Name`, `CountryCode`, `District`, `Population`) VALUES \n" +
					"(79,'Lanús','ARG','Buenos Aires',469735),\n" +
					"(80,'Merlo','ARG','Buenos Aires',463846);",
			},
			expected: output{
				insert: DMLInsert{
					Prefix: "INSERT INTO `city` (`ID`, `Name`, `CountryCode`, `District`, `Population`) ",
					Values: []string{
						"(79, 'Lanús', 'ARG', 'Buenos Aires', 469735)",
						"(80, 'Merlo', 'ARG', 'Buenos Aires', 463846)",
					},
				},
				ok: true,
			},
		},
		{
			given: input{
				drv: PostgresDialect,
				sql: `INSERT INTO public.customers (customer_id, company_name, contact_name, contact_title, address, city, region, postal_code, country, phone, fax) VALUES 
('ISLAT', 'Island Trading', 'Helen Bennett', 'Marketing Manager', 'Garden House Crowther Way', 'Cowes', 'Isle of Wight', 'PO31 7PJ', 'UK', '(198) 555-8888', NULL),
('KOENE', 'Königlich Essen', 'Philip Cramer', 'Sales Associate', 'Maubelstr. 90', 'Brandenburg', NULL, '14776', 'Germany', '0555-09876', NULL);`,
			},
			expected: output{
				insert: DMLInsert{
					Prefix: "INSERT INTO public.customers (customer_id, company_name, contact_name, contact_title, address, city, region, postal_code, country, phone, fax) ",
					Values: []string{
						"('ISLAT', 'Island Trading', 'Helen Bennett', 'Marketing Manager', 'Garden House Crowther Way', 'Cowes', 'Isle of Wight', 'PO31 7PJ', 'UK', '(198) 555-8888', NULL)",
						"('KOENE', e'K\\u00F6niglich Essen', 'Philip Cramer', 'Sales Associate', 'Maubelstr. 90', 'Brandenburg', NULL, '14776', 'Germany', '0555-09876', NULL)",
					},
				},
				ok: true,
			},
		},
	}

	for _, c := range tc {
		p, err := Parse(c.given.drv, c.given.sql)

		if assert.NoError(t, err) {
			var actual output
			actual.insert, actual.ok = p.AsDMLInsert()

			if c.expected.ok && assert.True(t, actual.ok) {
				assert.Equal(t, c.expected, actual)
			}
		}
	}
}
