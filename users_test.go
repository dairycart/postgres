package postgres

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"strconv"
	"testing"

	// internal dependencies
	"github.com/dairycart/dairymodels/v1"

	// external dependencies
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func setUserReadQueryExpectationByUsername(t *testing.T, mock sqlmock.Sqlmock, username string, toReturn *models.User, err error) {
	t.Helper()
	query := formatQueryForSQLMock(userQueryByUsername)
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
	mock.ExpectQuery(query).WithArgs(username).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestGetUserByUsername(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	client := NewPostgres()

	exampleUsername := "username"
	expected := &models.User{Username: exampleUsername}

	t.Run("optimal behavior", func(t *testing.T) {
		setUserReadQueryExpectationByUsername(t, mock, exampleUsername, expected, nil)
		actual, err := client.GetUserByUsername(mockDB, exampleUsername)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected user did not match actual user")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setUserWithUsernameExistenceQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, username string, shouldExist bool, err error) {
	t.Helper()
	query := formatQueryForSQLMock(userWithUsernameExistenceQuery)

	mock.ExpectQuery(query).
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(shouldExist))).
		WillReturnError(err)
}

func TestUserWithUsernameExists(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleUsername := "username"
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setUserWithUsernameExistenceQueryExpectation(t, mock, exampleUsername, true, nil)
		actual, err := client.UserWithUsernameExists(mockDB, exampleUsername)

		assert.NoError(t, err)
		assert.True(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with no rows found", func(t *testing.T) {
		setUserWithUsernameExistenceQueryExpectation(t, mock, exampleUsername, true, sql.ErrNoRows)
		actual, err := client.UserWithUsernameExists(mockDB, exampleUsername)

		assert.NoError(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with a database error", func(t *testing.T) {
		setUserWithUsernameExistenceQueryExpectation(t, mock, exampleUsername, true, errors.New("pineapple on pizza"))
		actual, err := client.UserWithUsernameExists(mockDB, exampleUsername)

		assert.NotNil(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

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
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setUserExistenceQueryExpectation(t, mock, exampleID, true, nil)
		actual, err := client.UserExists(mockDB, exampleID)

		assert.NoError(t, err)
		assert.True(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with no rows found", func(t *testing.T) {
		setUserExistenceQueryExpectation(t, mock, exampleID, true, sql.ErrNoRows)
		actual, err := client.UserExists(mockDB, exampleID)

		assert.NoError(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with a database error", func(t *testing.T) {
		setUserExistenceQueryExpectation(t, mock, exampleID, true, errors.New("pineapple on pizza"))
		actual, err := client.UserExists(mockDB, exampleID)

		assert.NotNil(t, err)
		assert.False(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
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

func TestGetUser(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	expected := &models.User{ID: exampleID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setUserReadQueryExpectation(t, mock, exampleID, expected, nil)
		actual, err := client.GetUser(mockDB, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected user did not match actual user")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setUserListReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, example *models.User, rowErr error, err error) {
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
		example.ID,
		example.FirstName,
		example.LastName,
		example.Username,
		example.Email,
		example.Password,
		example.Salt,
		example.IsAdmin,
		example.PasswordLastChangedOn,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).AddRow(
		example.ID,
		example.FirstName,
		example.LastName,
		example.Username,
		example.Email,
		example.Password,
		example.Salt,
		example.IsAdmin,
		example.PasswordLastChangedOn,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).AddRow(
		example.ID,
		example.FirstName,
		example.LastName,
		example.Username,
		example.Email,
		example.Password,
		example.Salt,
		example.IsAdmin,
		example.PasswordLastChangedOn,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).RowError(1, rowErr)

	query, _ := buildUserListRetrievalQuery(qf)

	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestGetUserList(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	example := &models.User{ID: exampleID}
	client := NewPostgres()
	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}

	t.Run("optimal behavior", func(t *testing.T) {
		setUserListReadQueryExpectation(t, mock, exampleQF, example, nil, nil)
		actual, err := client.GetUserList(mockDB, exampleQF)

		assert.NoError(t, err)
		assert.NotEmpty(t, actual, "list retrieval method should not return an empty slice")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error executing query", func(t *testing.T) {
		setUserListReadQueryExpectation(t, mock, exampleQF, example, nil, errors.New("pineapple on pizza"))
		actual, err := client.GetUserList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with error scanning values", func(t *testing.T) {
		exampleRows := sqlmock.NewRows([]string{"things"}).AddRow("stuff")
		query, _ := buildUserListRetrievalQuery(exampleQF)
		mock.ExpectQuery(formatQueryForSQLMock(query)).
			WillReturnRows(exampleRows)

		actual, err := client.GetUserList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with with row errors", func(t *testing.T) {
		setUserListReadQueryExpectation(t, mock, exampleQF, example, errors.New("pineapple on pizza"), nil)
		actual, err := client.GetUserList(mockDB, exampleQF)

		assert.NotNil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestBuildUserCountRetrievalQuery(t *testing.T) {
	t.Parallel()

	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}
	expected := `SELECT count(id) FROM users WHERE archived_on IS NULL LIMIT 25`
	actual, _ := buildUserCountRetrievalQuery(exampleQF)

	assert.Equal(t, expected, actual, "expected and actual queries should match")
}

func setUserCountRetrievalQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, count uint64, err error) {
	t.Helper()
	query, args := buildUserCountRetrievalQuery(qf)
	query = formatQueryForSQLMock(query)

	var argsToExpect []driver.Value
	for _, x := range args {
		argsToExpect = append(argsToExpect, x)
	}

	exampleRow := sqlmock.NewRows([]string{"count"}).AddRow(count)
	mock.ExpectQuery(query).WithArgs(argsToExpect...).WillReturnRows(exampleRow).WillReturnError(err)
}

func TestGetUserCount(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	client := NewPostgres()
	expected := uint64(123)
	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}

	t.Run("optimal behavior", func(t *testing.T) {
		setUserCountRetrievalQueryExpectation(t, mock, exampleQF, expected, nil)
		actual, err := client.GetUserCount(mockDB, exampleQF)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "count retrieval method should return the expected value")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setUserCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate *models.User, err error) {
	t.Helper()
	query := formatQueryForSQLMock(userCreationQuery)
	tt := buildTestTime(t)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(uint64(1), tt)
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
	assert.NoError(t, err)
	defer mockDB.Close()
	expectedID := uint64(1)
	exampleInput := &models.User{ID: expectedID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setUserCreationQueryExpectation(t, mock, exampleInput, nil)
		expectedCreatedOn := buildTestTime(t)

		actualID, actualCreatedOn, err := client.CreateUser(mockDB, exampleInput)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		assert.Equal(t, expectedCreatedOn, actualCreatedOn, "expected creation time did not match actual creation time")

		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setUserUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate *models.User, err error) {
	t.Helper()
	query := formatQueryForSQLMock(userUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(buildTestTime(t))
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
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleInput := &models.User{ID: uint64(1)}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setUserUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := buildTestTime(t)
		actual, err := client.UpdateUser(mockDB, exampleInput)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setUserDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(userDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(buildTestTime(t))
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeleteUserByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setUserDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := buildTestTime(t)
		actual, err := client.DeleteUser(mockDB, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with transaction", func(t *testing.T) {
		mock.ExpectBegin()
		setUserDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := buildTestTime(t)
		tx, err := mockDB.Begin()
		assert.NoError(t, err, "no error should be returned setting up a transaction in the mock DB")
		actual, err := client.DeleteUser(tx, exampleID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		assert.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
