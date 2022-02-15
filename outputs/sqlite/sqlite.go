/*
Copyright 2019 Cloudical Deutschland GmbH. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sqlite

import (
	"fmt"
	"path/filepath"

	"github.com/cloudical-io/ancientt/pkg/config"
	"github.com/cloudical-io/ancientt/pkg/util"
	"github.com/jmoiron/sqlx"

	// Include sqlite driver for sqlite output
	_ "github.com/mattn/go-sqlite3"

	"github.com/cloudical-io/ancientt/outputs"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// NameSQLite SQLite output name
const (
	NameSQLite = "sqlite"

	SQLiteIntType   = "BIGINT"
	SQLiteFloatType = "FLOAT"
	SQLiteBoolType  = "BOOLEAN"

	defaultNamePattern      = "ancientt-{{ .TestStartTime }}-{{ .Data.Tester }}.sqlite3"
	defaultTableNamePattern = "ancientt{{ .TestStartTime }}{{ .Data.Tester }}{{ .Data.ServerHost }}{{ .Data.ClientHost }}"

	createTableBeginQuery = "CREATE TABLE IF NOT EXISTS `%s` (\n"
	createTableEndQuery   = `);`
	insertDataBeginQuery  = "INSERT INTO %s VALUES ("
	insertDataEndQuery    = `);`
)

func init() {
	outputs.Factories[NameSQLite] = NewSQLiteOutput
}

// SQLite SQLite tester structure
type SQLite struct {
	outputs.Output
	logger *log.Entry
	config *config.SQLite
	dbCons map[string]*sqlx.DB
}

// NewSQLiteOutput return a new SQLite tester instance
func NewSQLiteOutput(cfg *config.Config, outCfg *config.Output) (outputs.Output, error) {
	if outCfg == nil {
		outCfg = &config.Output{
			SQLite: &config.SQLite{},
		}
	}
	s := SQLite{
		logger: log.WithFields(logrus.Fields{"output": NameSQLite}),
		config: outCfg.SQLite,
		dbCons: map[string]*sqlx.DB{},
	}
	if s.config.NamePattern == "" {
		s.config.NamePattern = defaultNamePattern
	}
	if s.config.TableNamePattern == "" {
		s.config.TableNamePattern = defaultTableNamePattern
	}

	return s, nil
}

// Do make SQLite outputs
func (s SQLite) Do(data outputs.Data) error {
	dataTable, ok := data.Data.(*outputs.Table)
	if !ok {
		return fmt.Errorf("data not in data table format for sqlite output")
	}

	filename, err := outputs.GetFilenameFromPattern(s.config.FilePath.NamePattern, "", data, nil)
	if err != nil {
		return err
	}

	tableName, err := outputs.GetFilenameFromPattern(s.config.TableNamePattern, "", data, nil)
	if err != nil {
		return err
	}

	var createTable bool

	outPath := filepath.Join(s.config.FilePath.FilePath, filename)
	db, ok := s.dbCons[outPath]
	if !ok {
		db, err = sqlx.Connect("sqlite3", outPath)
		if err != nil {
			return err
		}

		s.dbCons[outPath] = db
		createTable = true
	}

	if createTable {
		// Iterate over headers
		headers := []string{}
		for _, r := range dataTable.Headers {
			if r == nil {
				continue
			}
			headers = append(headers, util.CastToString(r.Value))
		}

		// Iterate over data columns to get the first row of data.
		// The first row of data is needed to set the types on the to be created SQLite table
		dataRows := []interface{}{}
		for _, row := range dataTable.Rows {
			for _, r := range row {
				if r == nil {
					continue
				}
				dataRows = append(dataRows, r.Value)
			}
			if len(dataRows) == 0 {
				continue
			}
			// Break after first round as we only need the first row!
			break
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("couldn't begin transaction in sqlite database. %+v", err)
		}
		tx.Exec(s.buildCreateTableQuery(tableName, headers, dataRows))
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("couldn't create table in sqlite database. %+v", err)
		}
	}

	// Iterate over data columns
	for _, row := range dataTable.Rows {
		dataRows := []interface{}{}
		for _, r := range row {
			dataRows = append(dataRows, r.Value)
		}
		if len(dataRows) == 0 {
			continue
		}

		query := s.buildInsertQuery(tableName, len(dataRows))
		if _, err := db.Exec(query, dataRows...); err != nil {
			return fmt.Errorf("couldn't insert data in sqlite database. %+v", err)
		}
	}

	return nil
}

func (s SQLite) buildCreateTableQuery(tableName string, columns []string, firstRow []interface{}) string {
	query := fmt.Sprintf(createTableBeginQuery, tableName)

	for i, c := range columns {
		cType := "TEXT"

		if len(firstRow) >= i+1 {
			switch firstRow[i].(type) {
			case bool:
				cType = SQLiteBoolType
			case float32:
				cType = SQLiteFloatType
			case float64:
				cType = SQLiteFloatType
			case int:
				cType = SQLiteIntType
			case int8:
				cType = SQLiteIntType
			case int16:
				cType = SQLiteIntType
			case int32:
				cType = SQLiteIntType
			case int64:
				cType = SQLiteIntType
			}
		}
		query += fmt.Sprintf("    `%s` %s", c, cType)
		if len(columns) != i+1 {
			query += ","
		}
		query += "\n"
	}

	query += createTableEndQuery

	return query
}

func (s SQLite) buildInsertQuery(tableName string, count int) string {
	query := fmt.Sprintf(insertDataBeginQuery, tableName)

	// Generate the placeholder `$1` and so on
	for i := 1; i <= count; i++ {
		query += fmt.Sprintf("$%d", i)
		if count >= i+1 {
			query += ", "
		}
	}

	query += insertDataEndQuery
	return query
}

// OutputFiles return a list of output files
func (s SQLite) OutputFiles() []string {
	list := []string{}
	for file := range s.dbCons {
		list = append(list, file)
	}
	return list
}

// Close closes all sqlite3 connections
func (s SQLite) Close() error {
	for name, db := range s.dbCons {
		s.logger.WithFields(logrus.Fields{"filepath": name}).Debug("closing db connection")
		if err := db.Close(); err != nil {
			s.logger.WithFields(logrus.Fields{"filepath": name}).Errorf("error closing db connection. %+v", err)
		}
	}

	return nil
}
