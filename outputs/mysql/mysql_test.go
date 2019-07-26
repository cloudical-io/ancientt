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
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cloudical-io/acntt/outputs"
	"github.com/cloudical-io/acntt/outputs/tests"
	"github.com/cloudical-io/acntt/pkg/config"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMySQL(t *testing.T) {
	// Create mock database
	db, mock, err := sqlmock.New()
	require.Nil(t, err)

	outCfg := &config.Output{
		MySQL: &config.MySQL{
			// Set a non-empty DSN, otherwise we get an error
			DSN: "username:password@127.0.0.1/mydb",
		},
	}

	// Generate mock Data with Table data
	data := tests.GenerateMockTableData(5)

	tableName, err := outputs.GetFilenameFromPattern(defaultTableNamePattern, "", data, nil)
	require.Nil(t, err)

	// Because the db driver already exists, the "CREATE TABLE" query is not triggered
	// Match the two inserts
	mock.ExpectExec(fmt.Sprintf("SELECT 1 FROM `%s` LIMIT 1", tableName)).WillReturnError(fmt.Errorf("table does not exist fake error"))
	mock.ExpectBegin()
	mock.ExpectExec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`", tableName))
	mock.ExpectCommit()
	mock.ExpectExec(fmt.Sprintf("INSERT INTO %s .*", tableName))
	mock.ExpectClose()

	dbx := sqlx.NewDb(db, "sqlmock")

	m, err := NewMySQLOutput(nil, outCfg)
	assert.Nil(t, err)

	// Cast the outputs.Output to the SQLite so we can manipulate the object
	ms, ok := m.(MySQL)
	require.True(t, ok)

	outPath := fmt.Sprintf("%s-%s", outCfg.MySQL.DSN, tableName)
	ms.dbCons[outPath] = dbx

	// Do() and Close() to run the database flow
	err = m.Do(data)
	assert.NotNil(t, err)
	err = m.Close()
	assert.Nil(t, err)

	// Check if all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)

	// TODO Verify data written to database
}
