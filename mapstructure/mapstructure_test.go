package mapstructure

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type dummy struct {
	Name     string        `mapstructure:"name"`
	Count    int           `mapstructure:"count"`
	Duration time.Duration `mapstructure:"duration"`
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		output  interface{}
		want    interface{}
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			input:   map[string]interface{}{"name": "foo", "count": "10", "duration": "5s"},
			output:  &dummy{},
			want:    &dummy{Name: "foo", Count: 10, Duration: 5 * time.Second},
			wantErr: assert.NoError,
		},
		{
			name:    "failedCreateDecoder",
			input:   map[string]interface{}{"name": "foo", "count": "10", "duration": "5s"},
			output:  dummy{},
			want:    dummy{},
			wantErr: assert.Error,
		},
		{
			name:    "failedDecode",
			input:   map[string]interface{}{"name": []string{}},
			output:  &dummy{},
			want:    &dummy{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Decode(tt.input, tt.output)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, tt.output)
		})
	}
}
