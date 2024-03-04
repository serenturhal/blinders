package explore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMatchKeyWithUserID(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "pass",
			args: args{userID: "userID"},
			want: "match:userID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, CreateMatchKeyWithUserID(tt.args.userID), "CreateMatchKeyWithUserID(%v)", tt.args.userID)
		})
	}
}
