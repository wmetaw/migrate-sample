package main

import (
	gspanner "cloud.google.com/go/spanner"
	sdb "cloud.google.com/go/spanner/admin/database/apiv1"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/spanner"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"io/ioutil"
)

var client *gspanner.Client
var adminClient *sdb.DatabaseAdminClient

func main() {
	ctx := context.Background()

	option := createOption(ctx)

	db := "projects/<PROJECT_ID>/instances/<INSTANCE_ID>/databases/<DB_NAME>"

	// create client
	client, _ := gspanner.NewClient(ctx, db, option)

	// create adminClient
	adminClient, _ := sdb.NewDatabaseAdminClient(ctx, option)

	//
	// I cannot create 'DB instance'. Because 'admin' and 'data' field are private.
	//
	driver, _ := spanner.WithInstance(&spanner.DB{
		admin: adminClient,
		data:  client,
	}, &spanner.Config{
		MigrationsTable: spanner.DefaultMigrationsTable,
		DatabaseName:    db,
	})

	m, err := migrate.NewWithDatabaseInstance("file://path/to/your", "<DB_NAME>", driver)

	err = m.Up()
	fmt.Println(err)
}

func createOption(ctx context.Context) option.ClientOption {
	data, _ := ioutil.ReadFile("your-credential.json")

	conf, _ := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/spanner.admin", "https://www.googleapis.com/auth/spanner.data")
	option := option.WithTokenSource(conf.TokenSource(ctx))
	return option
}
