package rbac

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

const failedMockDB = "Failed to create mock database connection: %v"
const mockCreateQuery = "CREATE TABLE IF NOT EXISTS rule_engine_casbin_rule"
const mockSelectQuery = "SELECT p_type,v0,v1,v2,v3,v4,v5 FROM rule_engine_casbin_rule"

func TestNewRolesManager(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(failedMockDB, err)
	}
	defer db.Close()

	// Set up expectations for the mock database
	mock.ExpectExec(mockCreateQuery).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(mockSelectQuery).WillReturnRows(sqlmock.NewRows([]string{"p_type", "v0", "v1", "v2", "v3", "v4", "v5"}))

	// Call the NewRolesManager function
	rolesManager := NewRolesManager(db)

	// Assert that the rolesManager is created correctly
	assert.NotNil(t, rolesManager)
	assert.IsType(t, rolesManager, rolesManager)

	// Assert that the database connection is set correctly
	assert.Equal(t, db, rolesManager.db)

	// Verify that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRolesManagerEnforce(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(failedMockDB, err)
	}
	defer db.Close()

	// Set up expectations for the mock database
	mock.ExpectExec(mockCreateQuery).WillReturnResult(sqlmock.NewResult(0, 0))

	rows := sqlmock.NewRows([]string{"p_type", "v0", "v1", "v2", "v3", "v4", "v5"}).
		AddRow("p", "user1", "data1", "read", "", "", "")

	mock.ExpectQuery(mockSelectQuery).WillReturnRows(rows)

	// Create a new RolesManager instance
	rolesManager := NewRolesManager(db)

	// Call the Enforce function
	valid, _ := rolesManager.Enforce("user1", "data1", "read")

	// Assert the result
	assert.True(t, valid)

	// Call the Enforce function
	valid2, _ := rolesManager.Enforce("user", "data1", "read")

	// Assert the result
	assert.False(t, valid2)

	// Verify that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func TestRolesManagerUpdatePermissionsForRole(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf(failedMockDB, err)
	}
	defer db.Close()

	// Set up expectations for the mock database
	//nolint:lll
	cq := "CREATE TABLE IF NOT EXISTS rule_engine_casbin_rule( p_type VARCHAR(32) DEFAULT '' NOT NULL, v0 VARCHAR(255) DEFAULT '' NOT NULL, v1 VARCHAR(255) DEFAULT '' NOT NULL, v2 VARCHAR(255) DEFAULT '' NOT NULL, v3 VARCHAR(255) DEFAULT '' NOT NULL, v4 VARCHAR(255) DEFAULT '' NOT NULL, v5 VARCHAR(255) DEFAULT '' NOT NULL, INDEX idx_rule_engine_casbin_rule (p_type,v0,v1) ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;"
	mock.ExpectExec(cq).WillReturnResult(sqlmock.NewResult(0, 0))
	rows := sqlmock.NewRows([]string{"p_type", "v0", "v1", "v2", "v3", "v4", "v5"}).
		AddRow("p", "role1", "data1", "read", "", "", "")
	mock.ExpectQuery(mockSelectQuery).WillReturnRows(rows)
	mock.ExpectExec("DELETE FROM rule_engine_casbin_rule WHERE p_type=? AND v0=?").
		WithArgs("p", "role1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO rule_engine_casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES (?,?,?,?,?,?,?)")
	mock.ExpectExec("INSERT INTO rule_engine_casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES (?,?,?,?,?,?,?)").
		WithArgs("p", "role1", "data1", "read", "", "", "").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	// Create a new RolesManager instance
	rolesManager := NewRolesManager(db)

	// Call the UpdatePermissionsForRole function
	_, err = rolesManager.UpdatePermissionsForRole("role1", [][]string{{"data1", "read"}})

	// Assert that no error occurred
	assert.NoError(t, err)

	// Verify that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRolesManagerGetAllRole(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(failedMockDB, err)
	}
	defer db.Close()

	// Set up expectations for the mock database
	mock.ExpectExec(mockCreateQuery).WillReturnResult(sqlmock.NewResult(0, 0))

	rows := sqlmock.NewRows([]string{"p_type", "v0", "v1", "v2", "v3", "v4", "v5"}).
		AddRow("p", "role1", "obj1", "read", "", "", "").
		AddRow("p", "role2", "obj2", "read", "", "", "")

	mock.ExpectQuery(mockSelectQuery).WillReturnRows(rows)

	// Create a new RolesManager instance
	rolesManager := NewRolesManager(db)

	// Call the GetAllRole function
	roles, _ := rolesManager.GetAllRole()

	// Assert the expected roles
	expectedRoles := []string{"role1", "role2"}
	assert.ElementsMatch(t, expectedRoles, roles)

	// Verify that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRolesManagerGetPermissionsForRole(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(failedMockDB, err)
	}
	defer db.Close()

	// Set up expectations for the mock database
	mock.ExpectExec(mockCreateQuery).WillReturnResult(sqlmock.NewResult(0, 0))

	rows := sqlmock.NewRows([]string{"p_type", "v0", "v1", "v2", "v3", "v4", "v5"}).
		AddRow("p", "role1", "data1", "read", "", "", "").
		AddRow("p", "role1", "data2", "write", "", "", "")

	mock.ExpectQuery(mockSelectQuery).WillReturnRows(rows)

	// Create a new RolesManager instance
	rolesManager := NewRolesManager(db)

	// Call the GetPermissionsForRole function
	permissions, err := rolesManager.GetPermissionsForRole("role1")

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert the expected permissions
	expectedPermissions := [][]string{{"role1", "data1", "read"}, {"role1", "data2", "write"}}
	assert.Equal(t, expectedPermissions, permissions)

	// Verify that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
