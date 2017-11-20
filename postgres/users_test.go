package postgres

import (
	"database/sql"
	"errors"
	"strconv"
	"testing"

	// internal dependencies
	"github.com/dairycart/dairycart/api/storage/models"

	// external dependencies
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func setUserExistenceQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, shouldExist bool, err error) {
	t.Helper()
	query := formatQueryForSQLMock(userExistenceQuery)

	mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(shouldExist))).
		WillReturnError(err)
}

func TestUserExists(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)

	t.Run("existing", func(t *testing.T) {
		setUserExistenceQueryExpectation(t, mock, exampleID, true, nil)
		client := Postgres{DB: mockDB}
		actual, err := client.UserExists(exampleID)

		require.Nil(t, err)
		require.True(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with no rows found", func(t *testing.T) {
		setUserExistenceQueryExpectation(t, mock, exampleID, true, sql.ErrNoRows)
		client := Postgres{DB: mockDB}
		actual, err := client.UserExists(exampleID)

		require.Nil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with a database error", func(t *testing.T) {
		setUserExistenceQueryExpectation(t, mock, exampleID, true, errors.New("pineapple on pizza"))
		client := Postgres{DB: mockDB}
		actual, err := client.UserExists(exampleID)

		require.NotNil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setUserReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, toReturn *models.User, err error) {
	t.Helper()
	query := formatQueryForSQLMock(userSelectionQuery)

	exampleRows := sqlmock.NewRows([]string{
		"id",
		"first_name",
		"last_name",
		"username",
		"email",
		"password",
		"salt",
		"is_admin",
		"password_last_changed_on",
		"created_on",
		"updated_on",
		"archived_on",
	}).AddRow(
		toReturn.ID,
		toReturn.FirstName,
		toReturn.LastName,
		toReturn.Username,
		toReturn.Email,
		toReturn.Password,
		toReturn.Salt,
		toReturn.IsAdmin,
		toReturn.PasswordLastChangedOn,
		toReturn.CreatedOn,
		toReturn.UpdatedOn,
		toReturn.ArchivedOn,
	)
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestGetUserByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()

	exampleID := uint64(1)
	expected := &models.User{ID: exampleID}

	t.Run("optimal behavior", func(t *testing.T) {
		setUserReadQueryExpectation(t, mock, exampleID, expected, nil)
		client := Postgres{DB: mockDB}
		actual, err := client.GetUser(exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected user did not match actual user")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setUserCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate *models.User, err error) {
	t.Helper()
	query := formatQueryForSQLMock(userCreationQuery)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(uint64(1), generateExampleTimeForTests())
	mock.ExpectQuery(query).
		WithArgs(
			toCreate.FirstName,
			toCreate.LastName,
			toCreate.Username,
			toCreate.Email,
			toCreate.Password,
			toCreate.Salt,
			toCreate.IsAdmin,
			toCreate.PasswordLastChangedOn,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestCreateUser(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	expectedID := uint64(1)
	exampleInput := &models.User{ID: expectedID}

	t.Run("optimal behavior", func(t *testing.T) {
		setUserCreationQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actualID, actualCreationDate, err := client.CreateUser(exampleInput)

		require.Nil(t, err)
		require.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		require.Equal(t, expected, actualCreationDate, "expected creation time did not match actual creation time")

		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setUserUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate *models.User, err error) {
	t.Helper()
	query := formatQueryForSQLMock(userUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests())
	mock.ExpectQuery(query).
		WithArgs(
			toUpdate.FirstName,
			toUpdate.LastName,
			toUpdate.Username,
			toUpdate.Email,
			toUpdate.Password,
			toUpdate.Salt,
			toUpdate.IsAdmin,
			toUpdate.PasswordLastChangedOn,
			toUpdate.ID,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestUpdateUserByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleInput := &models.User{ID: uint64(1)}

	t.Run("optimal behavior", func(t *testing.T) {
		setUserUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actual, err := client.UpdateUser(exampleInput)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setUserDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(userDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(generateExampleTimeForTests())
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeleteUserByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)

	t.Run("optimal behavior", func(t *testing.T) {
		setUserDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests()
		client := Postgres{DB: mockDB}
		actual, err := client.DeleteUser(exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
