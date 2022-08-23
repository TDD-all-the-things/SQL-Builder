package builder_test

import (
	"testing"

	"github.com/TDD-all-the-things/SQL-Builder/builder"

	"github.com/stretchr/testify/assert"
)

func TestInsertStmt(t *testing.T) {

	testcases := map[string]struct {
		input      any
		wantQuery  string
		wantValues []any
		wantErr    error
	}{
		"nil": {
			input:      nil,
			wantQuery:  "",
			wantValues: nil,
			wantErr:    nil,
		},
	}
	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			query, values, err := builder.InsertStmt(tc.input)
			assert.Equal(t, tc.wantQuery, query)
			assert.Equal(t, tc.wantValues, values)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
