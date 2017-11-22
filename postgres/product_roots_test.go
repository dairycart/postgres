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

func setProductRootWithSKUPrefixExistenceQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, skuPrefix string, shouldExist bool, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productRootWithSKUPrefixExistenceQuery)

	mock.ExpectQuery(query).
		WithArgs(skuPrefix).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(shouldExist))).
		WillReturnError(err)
}

func TestProductRootWithSKUExists(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleSKUPrefix := "example"
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setProductRootWithSKUPrefixExistenceQueryExpectation(t, mock, exampleSKUPrefix, true, nil)
		actual, err := client.ProductRootWithSKUPrefixExists(mockDB, exampleSKUPrefix)

		require.Nil(t, err)
		require.True(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with no rows found", func(t *testing.T) {
		setProductRootWithSKUPrefixExistenceQueryExpectation(t, mock, exampleSKUPrefix, true, sql.ErrNoRows)
		actual, err := client.ProductRootWithSKUPrefixExists(mockDB, exampleSKUPrefix)

		require.Nil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with a database error", func(t *testing.T) {
		setProductRootWithSKUPrefixExistenceQueryExpectation(t, mock, exampleSKUPrefix, true, errors.New("pineapple on pizza"))
		actual, err := client.ProductRootWithSKUPrefixExists(mockDB, exampleSKUPrefix)

		require.NotNil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductRootExistenceQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, shouldExist bool, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productRootExistenceQuery)

	mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(shouldExist))).
		WillReturnError(err)
}

func TestProductRootExists(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("existing", func(t *testing.T) {
		setProductRootExistenceQueryExpectation(t, mock, exampleID, true, nil)
		actual, err := client.ProductRootExists(mockDB, exampleID)

		require.Nil(t, err)
		require.True(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with no rows found", func(t *testing.T) {
		setProductRootExistenceQueryExpectation(t, mock, exampleID, true, sql.ErrNoRows)
		actual, err := client.ProductRootExists(mockDB, exampleID)

		require.Nil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with a database error", func(t *testing.T) {
		setProductRootExistenceQueryExpectation(t, mock, exampleID, true, errors.New("pineapple on pizza"))
		actual, err := client.ProductRootExists(mockDB, exampleID)

		require.NotNil(t, err)
		require.False(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductRootReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, toReturn *models.ProductRoot, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productRootSelectionQuery)

	exampleRows := sqlmock.NewRows([]string{
		"id",
		"name",
		"subtitle",
		"description",
		"sku_prefix",
		"manufacturer",
		"brand",
		"taxable",
		"cost",
		"product_weight",
		"product_height",
		"product_width",
		"product_length",
		"package_weight",
		"package_height",
		"package_width",
		"package_length",
		"quantity_per_package",
		"available_on",
		"created_on",
		"updated_on",
		"archived_on",
	}).AddRow(
		toReturn.ID,
		toReturn.Name,
		toReturn.Subtitle,
		toReturn.Description,
		toReturn.SKUPrefix,
		toReturn.Manufacturer,
		toReturn.Brand,
		toReturn.Taxable,
		toReturn.Cost,
		toReturn.ProductWeight,
		toReturn.ProductHeight,
		toReturn.ProductWidth,
		toReturn.ProductLength,
		toReturn.PackageWeight,
		toReturn.PackageHeight,
		toReturn.PackageWidth,
		toReturn.PackageLength,
		toReturn.QuantityPerPackage,
		toReturn.AvailableOn,
		toReturn.CreatedOn,
		toReturn.UpdatedOn,
		toReturn.ArchivedOn,
	)
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestGetProductRoot(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	expected := &models.ProductRoot{ID: exampleID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductRootReadQueryExpectation(t, mock, exampleID, expected, nil)
		actual, err := client.GetProductRoot(mockDB, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected productroot did not match actual productroot")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductRootListReadQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, qf *models.QueryFilter, example *models.ProductRoot, rowErr error, err error) {
	exampleRows := sqlmock.NewRows([]string{
		"id",
		"name",
		"subtitle",
		"description",
		"sku_prefix",
		"manufacturer",
		"brand",
		"taxable",
		"cost",
		"product_weight",
		"product_height",
		"product_width",
		"product_length",
		"package_weight",
		"package_height",
		"package_width",
		"package_length",
		"quantity_per_package",
		"available_on",
		"created_on",
		"updated_on",
		"archived_on",
	}).AddRow(
		example.ID,
		example.Name,
		example.Subtitle,
		example.Description,
		example.SKUPrefix,
		example.Manufacturer,
		example.Brand,
		example.Taxable,
		example.Cost,
		example.ProductWeight,
		example.ProductHeight,
		example.ProductWidth,
		example.ProductLength,
		example.PackageWeight,
		example.PackageHeight,
		example.PackageWidth,
		example.PackageLength,
		example.QuantityPerPackage,
		example.AvailableOn,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).AddRow(
		example.ID,
		example.Name,
		example.Subtitle,
		example.Description,
		example.SKUPrefix,
		example.Manufacturer,
		example.Brand,
		example.Taxable,
		example.Cost,
		example.ProductWeight,
		example.ProductHeight,
		example.ProductWidth,
		example.ProductLength,
		example.PackageWeight,
		example.PackageHeight,
		example.PackageWidth,
		example.PackageLength,
		example.QuantityPerPackage,
		example.AvailableOn,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).AddRow(
		example.ID,
		example.Name,
		example.Subtitle,
		example.Description,
		example.SKUPrefix,
		example.Manufacturer,
		example.Brand,
		example.Taxable,
		example.Cost,
		example.ProductWeight,
		example.ProductHeight,
		example.ProductWidth,
		example.ProductLength,
		example.PackageWeight,
		example.PackageHeight,
		example.PackageWidth,
		example.PackageLength,
		example.QuantityPerPackage,
		example.AvailableOn,
		example.CreatedOn,
		example.UpdatedOn,
		example.ArchivedOn,
	).RowError(1, rowErr)

	query, _ := buildProductRootListRetrievalQuery(qf)

	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestGetProductRootList(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	example := &models.ProductRoot{ID: exampleID}
	client := NewPostgres()
	exampleQF := &models.QueryFilter{
		Limit: 25,
		Page:  1,
	}

	t.Run("optimal behavior", func(t *testing.T) {
		setProductRootListReadQueryExpectation(t, mock, exampleQF, example, nil, nil)
		actual, err := client.GetProductRootList(mockDB, exampleQF)

		require.Nil(t, err)
		require.NotEmpty(t, actual, "list retrieval method should not return an empty slice")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with error executing query", func(t *testing.T) {
		setProductRootListReadQueryExpectation(t, mock, exampleQF, example, nil, errors.New("pineapple on pizza"))
		actual, err := client.GetProductRootList(mockDB, exampleQF)

		require.NotNil(t, err)
		require.Nil(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with error scanning values", func(t *testing.T) {
		exampleRows := sqlmock.NewRows([]string{"things"}).AddRow("stuff")
		query, _ := buildProductRootListRetrievalQuery(exampleQF)
		mock.ExpectQuery(formatQueryForSQLMock(query)).
			WillReturnRows(exampleRows)

		actual, err := client.GetProductRootList(mockDB, exampleQF)

		require.NotNil(t, err)
		require.Nil(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
	t.Run("with with row errors", func(t *testing.T) {
		setProductRootListReadQueryExpectation(t, mock, exampleQF, example, errors.New("pineapple on pizza"), nil)
		actual, err := client.GetProductRootList(mockDB, exampleQF)

		require.NotNil(t, err)
		require.Nil(t, actual)
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductRootCreationQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toCreate *models.ProductRoot, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productrootCreationQuery)
	exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(uint64(1), generateExampleTimeForTests(t))
	mock.ExpectQuery(query).
		WithArgs(
			toCreate.Name,
			toCreate.Subtitle,
			toCreate.Description,
			toCreate.SKUPrefix,
			toCreate.Manufacturer,
			toCreate.Brand,
			toCreate.Taxable,
			toCreate.Cost,
			toCreate.ProductWeight,
			toCreate.ProductHeight,
			toCreate.ProductWidth,
			toCreate.ProductLength,
			toCreate.PackageWeight,
			toCreate.PackageHeight,
			toCreate.PackageWidth,
			toCreate.PackageLength,
			toCreate.QuantityPerPackage,
			toCreate.AvailableOn,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestCreateProductRoot(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	expectedID := uint64(1)
	exampleInput := &models.ProductRoot{ID: expectedID}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductRootCreationQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests(t)
		actualID, actualCreationDate, err := client.CreateProductRoot(mockDB, exampleInput)

		require.Nil(t, err)
		require.Equal(t, expectedID, actualID, "expected and actual IDs don't match")
		require.Equal(t, expected, actualCreationDate, "expected creation time did not match actual creation time")

		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductRootUpdateQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, toUpdate *models.ProductRoot, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productRootUpdateQuery)
	exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(generateExampleTimeForTests(t))
	mock.ExpectQuery(query).
		WithArgs(
			toUpdate.Name,
			toUpdate.Subtitle,
			toUpdate.Description,
			toUpdate.SKUPrefix,
			toUpdate.Manufacturer,
			toUpdate.Brand,
			toUpdate.Taxable,
			toUpdate.Cost,
			toUpdate.ProductWeight,
			toUpdate.ProductHeight,
			toUpdate.ProductWidth,
			toUpdate.ProductLength,
			toUpdate.PackageWeight,
			toUpdate.PackageHeight,
			toUpdate.PackageWidth,
			toUpdate.PackageLength,
			toUpdate.QuantityPerPackage,
			toUpdate.AvailableOn,
			toUpdate.ID,
		).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestUpdateProductRootByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleInput := &models.ProductRoot{ID: uint64(1)}
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductRootUpdateQueryExpectation(t, mock, exampleInput, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.UpdateProductRoot(mockDB, exampleInput)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func setProductRootDeletionQueryExpectation(t *testing.T, mock sqlmock.Sqlmock, id uint64, err error) {
	t.Helper()
	query := formatQueryForSQLMock(productRootDeletionQuery)
	exampleRows := sqlmock.NewRows([]string{"archived_on"}).AddRow(generateExampleTimeForTests(t))
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestDeleteProductRootByID(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer mockDB.Close()
	exampleID := uint64(1)
	client := NewPostgres()

	t.Run("optimal behavior", func(t *testing.T) {
		setProductRootDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		actual, err := client.DeleteProductRoot(mockDB, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})

	t.Run("with transaction", func(t *testing.T) {
		mock.ExpectBegin()
		setProductRootDeletionQueryExpectation(t, mock, exampleID, nil)
		expected := generateExampleTimeForTests(t)
		tx, err := mockDB.Begin()
		require.Nil(t, err, "no error should be returned setting up a transaction in the mock DB")
		actual, err := client.DeleteProductRoot(tx, exampleID)

		require.Nil(t, err)
		require.Equal(t, expected, actual, "expected deletion time did not match actual deletion time")
		require.Nil(t, mock.ExpectationsWereMet(), "not all database expectations were met")
	})
}
