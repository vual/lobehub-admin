package model

import (
	"net/url"
	"os"
	"testing"

	"github.com/QuantumNous/new-api/common"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizePostgresDSNUsesIsolatedSchemaByDefault(t *testing.T) {
	dsn, schema, err := normalizePostgresDSN(
		"postgresql://root:pass@localhost:5432/app?sslmode=disable",
		"",
	)

	require.NoError(t, err)
	assert.Equal(t, defaultPostgresSchema, schema)
	parsed, err := url.Parse(dsn)
	require.NoError(t, err)
	assert.Equal(t, defaultPostgresSchema, parsed.Query().Get("search_path"))
	assert.Equal(t, "disable", parsed.Query().Get("sslmode"))
}

func TestNormalizePostgresDSNHonorsConfiguredSchema(t *testing.T) {
	dsn, schema, err := normalizePostgresDSN(
		"postgres://root:pass@localhost:5432/app?search_path=from_dsn",
		"tenant_gateway",
	)

	require.NoError(t, err)
	assert.Equal(t, "tenant_gateway", schema)
	parsed, err := url.Parse(dsn)
	require.NoError(t, err)
	assert.Equal(t, "tenant_gateway", parsed.Query().Get("search_path"))
}

func TestNormalizePostgresDSNHonorsDSNSchema(t *testing.T) {
	dsn, schema, err := normalizePostgresDSN(
		"postgres://root:pass@localhost:5432/app?search_path=existing_schema",
		"",
	)

	require.NoError(t, err)
	assert.Equal(t, "existing_schema", schema)
	parsed, err := url.Parse(dsn)
	require.NoError(t, err)
	assert.Equal(t, "existing_schema", parsed.Query().Get("search_path"))
}

func TestNormalizePostgresDSNRejectsUnsafeSchema(t *testing.T) {
	_, _, err := normalizePostgresDSN(
		"postgres://root:pass@localhost:5432/app",
		"new_api,public",
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid PostgreSQL schema name")
}

func TestChooseDBCreatesConfiguredPostgresSchema(t *testing.T) {
	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("TEST_POSTGRES_DSN is not configured")
	}

	const schema = "new_api_choose_db_schema_integration_test"
	adminDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), newGormConfig(false))
	require.NoError(t, err)
	var existed bool
	require.NoError(t, adminDB.Raw(
		"SELECT EXISTS (SELECT 1 FROM information_schema.schemata WHERE schema_name = ?)",
		schema,
	).Scan(&existed).Error)
	if existed {
		t.Skip(schema + " already exists")
	}

	t.Setenv("SQL_DSN", dsn)
	t.Setenv("SQL_SCHEMA", schema)
	db, dbType, err := chooseDB("SQL_DSN", false)
	require.NoError(t, err)
	assert.Equal(t, common.DatabaseTypePostgreSQL, dbType)
	t.Cleanup(func() {
		sqlDB, closeErr := db.DB()
		if closeErr == nil {
			require.NoError(t, sqlDB.Close())
		}
		require.NoError(t, adminDB.Exec(`DROP SCHEMA "`+schema+`" CASCADE`).Error)
		adminSQLDB, closeErr := adminDB.DB()
		if closeErr == nil {
			require.NoError(t, adminSQLDB.Close())
		}
	})

	var currentSchema string
	require.NoError(t, db.Raw("SELECT current_schema()").Scan(&currentSchema).Error)
	assert.Equal(t, schema, currentSchema)
}
