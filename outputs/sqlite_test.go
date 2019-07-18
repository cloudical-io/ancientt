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

package outputs

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cloudical-io/acntt/pkg/config"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLite(t *testing.T) {
	// Create mock database
	db, mock, err := sqlmock.New()
	require.Nil(t, err)

	tempDir := os.TempDir()
	outName := fmt.Sprintf("acntt-test-%s.sqlite3", t.Name())
	outCfg := &config.Output{
		SQLite: &config.SQLite{
			FilePath:    tempDir,
			NamePattern: outName,
		},
	}

	// Generate mock Data with Table data
	data := generateMockTableData(5)

	filename, err := getFilenameFromPattern(outCfg.SQLite.NamePattern, "", data, nil)
	require.Nil(t, err)

	tableName, err := getFilenameFromPattern(outCfg.SQLite.TableNamePattern, "", data, nil)
	require.Nil(t, err)

	// Because the db driver already exists, the "CREATE TABLE" query is not triggered
	// Match the two inserts
	mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", tableName))
	mock.ExpectClose()

	dbx := sqlx.NewDb(db, "sqlmock")

	m, err := NewSQLiteOutput(nil, outCfg)
	assert.Nil(t, err)

	// Cast the outputs.Output to the SQLite so we can manipulate the object
	ms, ok := m.(SQLite)
	require.True(t, ok)

	outPath := filepath.Join(outCfg.SQLite.FilePath, filename)
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
