package database

import (
	"fmt"
	"github.com/anjanashankar9/random-learning/go-sqlite/database/client"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"

	"git.soma.salesforce.com/Infrastructure-Security/go-common/pkg/logging"
)

// Read config as required
// Create constants for all database statements

const (
	logSource = "database.go"
	ALL       = "All"

	// SQL Queries
	CreateTable = `DROP TABLE IF EXISTS eip_new;
				CREATE TABLE eip_new(id INTEGER PRIMARY KEY, 
				fi TEXT,
				fd TEXT,
				env TEXT,
				service TEXT,
				externalIP TEXT,
				used_for_fd_egress TEXT, 
				tags TEXT,
				INDEX env_index (env));`

	CreateTableWithIndex = `BEGIN;
							DROP TABLE IF EXISTS eip_new;
							CREATE TABLE 
								eip_new(id INTEGER PRIMARY KEY, 
								fi TEXT,
								fd TEXT,
								env TEXT,
								service TEXT,
								externalIP TEXT,
								used_for_fd_egress TEXT, 
								tags TEXT
							);
							CREATE INDEX env_idx ON eip_new (env);
							COMMIT;`
	CreateTableWithIndex2 = `DROP TABLE IF EXISTS eip_new;
							CREATE TABLE
								eip_new(id INTEGER PRIMARY KEY,
								fi TEXT,
								fd TEXT,
								env TEXT,
								service TEXT,
								externalIP TEXT,
								used_for_fd_egress TEXT,
								tags TEXT,
								INDEX env_idx ON eip_new (env)
							);`

	InsertIntoEIP = `INSERT INTO eip_new(fi, fd, env, service, externalIP, used_for_fd_egress, tags) 
					VALUES(:fi, :fd, :env, :service, :externalIP, :used_for_fd_egress, :tags)`
	RenameTable = `ALTER TABLE eip_new RENAME TO eip;`
	DropTable   = `DROP TABLE IF EXISTS eip;`

	SelectAllQuery             = `SELECT * FROM eip`
	SelectEIPId                = `SELECT * FROM eip WHERE id = ?`
	SelectQueryWithLimitOffset = `SELECT * FROM eip LIMIT ?,?`

	SelectQueryWithEgressTag = `SELECT * FROM eip where used_for_fd_egress = ? COLLATE NOCASE and fi = "aws-dev1-uswest" order by fi DESC`

	// Count queries
	SelectCountAll = `SELECT COUNT(*) FROM eip`
	SelectCountEnv = `SELECT COUNT(*) FROM eip where env = ?`
)

type EIP struct {
	Id              int64  `db:"id"`
	FI              string `db:"fi"`
	FD              string `db:"fd"`
	Env             string `db:"env"`
	Service         string `db:"service"`
	ExternalIP      string `db:"externalIP"`
	UsedForFdEgress string `db:"used_for_fd_egress"`
	Tags            string `db:"tags"`
}

type DatabaseClient struct {
	client client.Database
	log    *logrus.Entry
}

func New(db client.Database) DatabaseClient {
	return DatabaseClient{
		client: db,
		log:    logging.New(logSource),
	}
}

func (db *DatabaseClient) CreateTable() error {
	_, err := db.client.Exec(CreateTableWithIndex)

	if err != nil {
		db.log.Errorf("unable to create table: %s", err)
		return err
	}

	return nil
}

func (db *DatabaseClient) RenameTable() error {
	_, err := db.client.Exec(DropTable)
	if err != nil {
		db.log.Errorf("unable to drop table: %s", err)
		return err
	}
	_, err = db.client.Exec(RenameTable)

	if err != nil {
		db.log.Errorf("unable to create table: %s", err)
		return err
	}

	return nil
}

func (db *DatabaseClient) InsertIntoTable(eip EIP) error {
	_, err := db.client.NamedExec(InsertIntoEIP, eip)
	if err != nil {
		db.log.WithError(err).Error("Error in insert query")
		return err
	}
	return nil
}

func (db *DatabaseClient) GetAll() (resEvents []EIP, err error) {
	err = db.client.Select(&resEvents, SelectAllQuery)
	if err != nil {
		db.log.WithError(err).Error("Error in select query")
		return resEvents, err
	}
	return resEvents, nil
}

func (db *DatabaseClient) GetWithId(i int) (resEvents []EIP, err error) {
	err = db.client.Select(&resEvents, SelectEIPId, 1, 1)
	if err != nil {
		db.log.WithError(err).Error("Error in select query")
		return resEvents, err
	}
	return resEvents, nil
}

func (db *DatabaseClient) GetWithLimitOffset(limit, offset int) (resEvents []EIP, err error) {
	err = db.client.Select(&resEvents, SelectQueryWithLimitOffset, limit, offset)
	if err != nil {
		db.log.WithError(err).Error("Error in select query")
		return resEvents, err
	}
	return resEvents, nil
}

func (db *DatabaseClient) GetWithTag(value string) (resEvents []EIP, err error) {
	err = db.client.Select(&resEvents, SelectQueryWithEgressTag, value)
	if err != nil {
		db.log.WithError(err).Error("Error in select query")
		return resEvents, err
	}
	return resEvents, nil
}

func (db *DatabaseClient) GetCount(env string) (id int, err error) {
	if strings.EqualFold(env, ALL) {
		err = db.client.Get(&id, SelectCountAll)
	} else {
		err = db.client.Get(&id, SelectCountEnv, env)
	}

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *DatabaseClient) ReadInBatches() error {
	wg := &sync.WaitGroup{}
	ch := make(chan []byte)
	wg.Add(2)
	f, ferr := os.Create("eip.csv")
	if ferr != nil {
		fmt.Errorf("Unable to open file for writing %s", ferr)
	}

	defer f.Close()

	go db.readDatabase(wg, ch)
	go db.writeToCSV(wg, ch, f)
	wg.Wait()
	return nil
}

func (db *DatabaseClient) readDatabase(wg *sync.WaitGroup, ch chan []byte) {
	defer wg.Done()
	fmt.Println("Inside read")
	for i := 0; i < 1000; i++ {
		ch <- []byte("Hello")
	}
	close(ch)
}

func (db *DatabaseClient) writeToCSV(wg *sync.WaitGroup, ch chan []byte, f *os.File) {
	defer wg.Done()
	fmt.Println("Inside write")
	for msg := range ch {
		s := fmt.Sprintf("%s", msg)
		_, err := f.WriteString(s)
		if err != nil {
			fmt.Errorf("Cannot write %s to file %s", s, err)
		}
	}
}
