package benchmarks_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/brianvoe/gofakeit"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"github.com/lucaskatayama/benchmark-db/benchmarks"
	"github.com/rocketlaunchr/dbq/v2"
)

var db *sql.DB
var g *gorm.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("mysql", "root:secret@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}
	g, err = gorm.Open("mysql", db)
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func setup() {
	var err error
	// Create table
	createQ := `
	CREATE TABLE tests (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		name varchar(50) NOT NULL DEFAULT '',
		email varchar(150) NOT NULL DEFAULT '',
		PRIMARY KEY (id)
	)`

	_, err = db.Exec(createQ)
	if err != nil {
		panic(err)
	}

	// Add 10,000 fake entries
	entries := []interface{}{}
	for i := 0; i < 10000; i++ {
		entry := []interface{}{
			i + 1,
			gofakeit.Name(),  // Fake name
			gofakeit.Email(), // Fake email
		}
		entries = append(entries, entry)
	}
	stmt := dbq.INSERTStmt("tests", []string{"id", "name", "email"}, len(entries))
	_, err = dbq.E(context.Background(), db, stmt, nil, entries)
	if err != nil {
		panic(err)
	}
}

func cleanup() {
	// Delete table
	_, err := db.Exec(`DROP TABLE IF EXISTS tests`)
	if err != nil {
		panic(err)
	}
}

func BenchmarkAll(b *testing.B) {
	cleanup()
	setup()
	defer cleanup()

	limits := []int{
		5,
		50,
		500,
		10000,
	}
	b.ReportAllocs()
	for _, lim := range limits { // Fetch varying number of rows
		lim := lim

		// Benchmark dbq
		b.Run(fmt.Sprintf("dbq  limit:%d", lim), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				q := fmt.Sprintf("SELECT id, name, email FROM tests ORDER BY id LIMIT %d", lim)
				_, err := dbq.Qs(context.Background(), db, q, benchmarks.Model{}, nil)
				if err != nil {
					b.Fatal(err)
				}
				//b.Log(len(res.([]*benchmarks.Model)))
			}
		})

		// Benchmark sqlx
		b.Run(fmt.Sprintf("sqlx limit:%d", lim), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				db := sqlx.NewDb(db, "mysql")
				q := fmt.Sprintf("SELECT id, name, email FROM tests ORDER BY id LIMIT %d", lim)

				res := []benchmarks.Model{}
				err := db.Select(&res, q)
				if err != nil {
					b.Fatal(err)
				}
				//b.Log(len(res))
			}
		})

		// Benchmark gorm
		b.Run(fmt.Sprintf("gorm limit:%d", lim), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var res = []benchmarks.Model{}

				err := g.Order("id").Limit(lim).Find(&res).Error
				if err != nil {
					b.Fatal(err)
				}
				//b.Log(len(res))
			}
		})

		b.Run(fmt.Sprintf("std  limit:%d", lim), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				q := fmt.Sprintf("SELECT id, name, email FROM tests ORDER BY id LIMIT %d", lim)
				rows, err := db.Query(q)
				if err != nil {
					b.Fatal(err)
				}
				var res []benchmarks.Model
				for rows.Next() {
					var row benchmarks.Model
					if err := rows.Scan(&row.ID, &row.Name, &row.Email); err != nil {
						b.Fatal(err)
					}
					res = append(res, row)
				}
				//b.Log(len(res))
			}
		})
		fmt.Println("============")
	}

}
