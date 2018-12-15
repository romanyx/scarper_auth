package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Migrate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	t.Log("with given database connection.")
	{
		m, err := newMigration(db)
		assert.Nil(t, err)

		t.Log("\ttest:0\tshould up schema.")
		{
			err := Migrate(db)
			assert.Nil(t, err)
		}

		t.Log("\ttest:1\tshould down schema.")
		{
			err := m.Down()
			assert.Nil(t, err)
		}
	}
}
