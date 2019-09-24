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

package mysql

import (
	"fmt"

	"github.com/cloudical-io/ancientt/pkg/config"
	"github.com/cloudical-io/ancientt/pkg/util"
	"github.com/jmoiron/sqlx"

	// Include MySQL driver for mysql output
	_ "github.com/go-sql-driver/mysql"

	"github.com/cloudical-io/ancientt/outputs"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// NameMySQL MySQL output name
const (
	NameMySQL      = "mysql"
	MySQLIntType   = "BIGINT"
	MySQLFloatType = "FLOAT"
	MySQLBoolType  = "BOOLEAN"
)

func init() {
	outputs.Factories[NameMySQL] = NewMySQLOutput
}

// MySQL MySQL tester structure
type MySQL struct {
	outputs.Output
	logger *log.Entry
	config *config.MySQL
	dbCons map[string]*sqlx.DB
}

const (
	defaultTableNamePattern = "ancientt{{ .TestStartTime }}{{ .Data.Tester }}{{ .Data.ServerHost }}{{ .Data.ClientHost }}"

	checkIfTableExistsQuery = "SELECT 1 FROM `%s` LIMIT 1;"
	createTableBeginQuery   = "CREATE TABLE IF NOT EXISTS `%s` (\n"
	createTableEndQuery     = `);`
	insertDataBeginQuery    = "INSERT INTO %s VALUES ("
	insertDataEndQuery      = `);`
)

// NewMySQLOutput return a new MySQL tester instance
func NewMySQLOutput(cfg *config.Config, outCfg *config.Output) (outputs.Output, error) {
	if outCfg == nil {
		outCfg = &config.Output{
			MySQL: &config.MySQL{
				AutoCreateTables: util.BoolPointer(true),
			},
		}
	}
	if outCfg.MySQL.AutoCreateTables == nil {
		outCfg.MySQL.AutoCreateTables = util.BoolPointer(true)
	}
	m := MySQL{
		logger: log.WithFields(logrus.Fields{"output": NameMySQL}),
		config: outCfg.MySQL,
		dbCons: map[string]*sqlx.DB{},
	}
	if m.config.DSN == "" {
		return nil, fmt.Errorf("no DSN for mysql connection given")
	}
	if m.config.TableNamePattern == "" {
		m.config.TableNamePattern = defaultTableNamePattern
	}

	return m, nil
}

// Do make MySQL outputs
func (m MySQL) Do(data outputs.Data) error {
	dataTable, ok := data.Data.(outputs.Table)
	if !ok {
		return fmt.Errorf("data not in table for mysql output")
	}

	tableName, err := outputs.GetFilenameFromPattern(m.config.TableNamePattern, "", data, nil)
	if err != nil {
		return err
	}

	dbPath := fmt.Sprintf("%s-%s", m.config.DSN, tableName)

	db, ok := m.dbCons[dbPath]
	if !ok {
		db, err = sqlx.Connect("mysql", dbPath)
		if err != nil {
			return err
		}

		m.dbCons[dbPath] = db
	}

	if err := m.createTable(db, dataTable, tableName); err != nil {
		return err
	}

	// Iterate over data columns
	for _, column := range dataTable.Columns {
		dataColumn := []interface{}{}
		for _, row := range column.Rows {
			dataColumn = append(dataColumn, row.Value)
		}
		if len(dataColumn) == 0 {
			continue
		}

		query := m.buildInsertQuery(tableName, len(dataColumn))
		if _, err := db.Exec(query, dataColumn...); err != nil {
			return fmt.Errorf("couldn't insert data in mysql database. %+v", err)
		}
	}

	return nil
}

func (m MySQL) createTable(db *sqlx.DB, dataTable outputs.Table, tableName string) error {
	// Iterate over headers
	headerColumns := []string{}
	for _, column := range dataTable.Headers {
		for _, row := range column.Rows {
			headerColumns = append(headerColumns, util.CastToString(row.Value))
		}
		if len(headerColumns) == 0 {
			continue
		}

	}

	// Iterate over data columns to get the first row of data.
	// The first row of data is needed to set the types on the to be created MySQL table
	dataColumn := []interface{}{}
	for _, column := range dataTable.Columns {
		for _, row := range column.Rows {
			dataColumn = append(dataColumn, row.Value)
		}
		if len(dataColumn) == 0 {
			continue
		}
		// Break after first round as we only need the first row!
		break
	}

	// The error should not return an error when the table exists, try to create the database
	if _, err := db.Exec(fmt.Sprintf(checkIfTableExistsQuery, tableName)); err != nil {
		// Only auto create tables when enabled
		if *m.config.AutoCreateTables {
			// Start transaction, exec the CREATE TABLE query and commit the result
			tx, err := db.Begin()
			if err != nil {
				return fmt.Errorf("couldn't begin transaction in mysql database. %+v", err)
			}
			tx.Exec(m.buildCreateTableQuery(tableName, headerColumns, dataColumn))
			if err := tx.Commit(); err != nil {
				return fmt.Errorf("couldn't create table in mysql database. %+v", err)
			}
		} else {
			return fmt.Errorf("table %s doesn't exist in mysql database and AutoCreateTables is false", tableName)
		}
	}

	return nil
}

func (m MySQL) buildCreateTableQuery(tableName string, columns []string, firstRow []interface{}) string {
	query := fmt.Sprintf(createTableBeginQuery, tableName)

	for i, c := range columns {
		cType := "TEXT"

		if len(firstRow) >= i+1 {
			switch firstRow[i].(type) {
			case bool:
				cType = MySQLBoolType
			case float32:
				cType = MySQLFloatType
			case float64:
				cType = MySQLFloatType
			case int:
				cType = MySQLIntType
			case int8:
				cType = MySQLIntType
			case int16:
				cType = MySQLIntType
			case int32:
				cType = MySQLIntType
			case int64:
				cType = MySQLIntType
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

func (m MySQL) buildInsertQuery(tableName string, count int) string {
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
func (m MySQL) OutputFiles() []string {
	return []string{}
}

// Close closes all MySQL connections
func (m MySQL) Close() error {
	for name, db := range m.dbCons {
		m.logger.WithFields(logrus.Fields{"filepath": name}).Debug("closing db connection")
		if err := db.Close(); err != nil {
			m.logger.WithFields(logrus.Fields{"filepath": name}).Errorf("error closing db connection. %+v", err)
		}
	}

	return nil
}
