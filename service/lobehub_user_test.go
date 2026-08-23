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
package service

import (
	"encoding/hex"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashBetterAuthPasswordMatchesBetterAuthScryptFormat(t *testing.T) {
	salt, err := hex.DecodeString("00112233445566778899aabbccddeeff")
	require.NoError(t, err)

	hash, err := hashBetterAuthPasswordWithSalt("Pa\u0308ssword!", salt)
	require.NoError(t, err)
	assert.Equal(t,
		"00112233445566778899aabbccddeeff:ed6138f755167c0c66f5d3807435c572d558c0ca490d25eb39d91e4c52d08c5047613a95e61e98e5c2875fde74bb01c6b1316cc9a9deca960c7f3f08d6865e56",
		hash,
	)
}

func TestGenerateLobeHubTemporaryPasswordContract(t *testing.T) {
	password, err := generateLobeHubTemporaryPassword()
	require.NoError(t, err)
	assert.Len(t, password, 24)
	assert.True(t, containsRuneClass(password, unicode.IsUpper))
	assert.True(t, containsRuneClass(password, unicode.IsLower))
	assert.True(t, containsRuneClass(password, unicode.IsDigit))
	assert.Contains(t, password, "!")
}

func TestNormalizeLobeHubTextTurnsBlankIntoNull(t *testing.T) {
	blank := "  "
	normalized, err := normalizeLobeHubText(&blank)
	require.NoError(t, err)
	assert.Nil(t, normalized)
}

func TestNormalizeLobeHubEmailCanonicalizesCase(t *testing.T) {
	email := " User@Example.COM "
	normalized, err := normalizeLobeHubEmail(&email)
	require.NoError(t, err)
	require.NotNil(t, normalized)
	assert.Equal(t, "user@example.com", *normalized)
}

func containsRuneClass(value string, match func(rune) bool) bool {
	for _, character := range value {
		if match(character) {
			return true
		}
	}
	return false
}
