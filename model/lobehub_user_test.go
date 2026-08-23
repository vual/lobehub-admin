/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
package model

import (
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLobeHubSchemaDefaultsToPublic(t *testing.T) {
	t.Setenv("LOBEHUB_DB_SCHEMA", "")
	schema, err := lobeHubSchema()
	require.NoError(t, err)
	assert.Equal(t, "public", schema)
}

func TestLobeHubSchemaRejectsUnsafeIdentifier(t *testing.T) {
	t.Setenv("LOBEHUB_DB_SCHEMA", `public"; DROP SCHEMA admin; --`)
	_, err := lobeHubSchema()
	assert.ErrorIs(t, err, ErrLobeHubSchemaUnavailable)
}

func TestEscapeLobeHubSearchEscapesLikeWildcards(t *testing.T) {
	assert.Equal(t, `50\%\_off\\today`, escapeLobeHubSearch(`50%_off\today`))
}

func TestNormalizeLobeHubWriteErrorMapsUniqueViolation(t *testing.T) {
	err := normalizeLobeHubWriteError(&pgconn.PgError{Code: "23505"})
	assert.ErrorIs(t, err, ErrLobeHubConflict)
}
