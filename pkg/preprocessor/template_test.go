package preprocessor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTemplateTokens(t *testing.T) {
	tests := []struct {
		name       string
		template   string
		replaceIdx map[string]string
		want       string
	}{
		{
			name:     "Test 1",
			template: "{{env_name:32}}-{{stack_name:32}}",
			replaceIdx: map[string]string{
				"env_name":   "MyFavoriteEnv",
				"stack_name": "CoolestStack",
			},
			want: "MyFavoriteEnv-CoolestStack",
		},
		{
			name:     "Test 1",
			template: "{{env_name:32}}-{{stack_name:32}}",
			replaceIdx: map[string]string{
				"env_name":   "MyFavoriteEnvMyFavoriteEnvMyFavoriteEnvMyFavoriteEnvMyFavoriteEnvMyFavoriteEnv",
				"stack_name": "CoolestStack",
			},
			want: "MyFavoriteEnvMyFavoriteEnvMyFavo-CoolestStack",
		},
		{
			name:     "Test 1",
			template: "{{env_name:32}}-{{env_name:32}}!!-{{stack_name:8}}",
			replaceIdx: map[string]string{
				"env_name":   "MyFavoriteEnv",
				"stack_name": "CoolestStack",
			},
			want: "MyFavoriteEnv-MyFavoriteEnv!!-CoolestS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTemplateTokens(tt.template, tt.replaceIdx)
			assert.Equalf(t, tt.want, got, "ParseTemplateTokens(%v, %v)", tt.template, tt.replaceIdx)
		})
	}
}
