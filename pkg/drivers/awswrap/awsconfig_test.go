package awswrap

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/swizzleio/swiz/mocks/ext/aws"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetDefaultConfig(t *testing.T) {
	testCases := []struct {
		name           string
		getCallerIDErr error
		listAliasesErr error
		accountAliases []string
		expected       *AwsConfig
		expectErr      bool
	}{
		{
			name:           "returns default config",
			getCallerIDErr: nil,
			listAliasesErr: nil,
			accountAliases: []string{"coolalias"},
			expected: &AwsConfig{
				Profile:   "coolalias",
				AccountId: "test-account-id",
				Region:    "us-somewhere-1",
			},
			expectErr: false,
		},
		{
			name:           "returns config without alias",
			getCallerIDErr: nil,
			listAliasesErr: nil,
			accountAliases: []string{},
			expected: &AwsConfig{
				Profile:   "dev",
				AccountId: "test-account-id",
				Region:    "us-somewhere-1",
			},
			expectErr: false,
		},
		{
			name:           "returns error when GetCallerIdentity fails",
			getCallerIDErr: errors.New("error"),
			listAliasesErr: nil,
			expected:       nil,
			expectErr:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock AWS SDK calls
			stsMock := &mockaws.Stser{}
			iamMock := &mockaws.Iamer{}
			stsMock.On("GetCallerIdentity", mock.Anything, &sts.GetCallerIdentityInput{}).Return(&sts.GetCallerIdentityOutput{
				Account: aws.String("test-account-id"),
			}, tc.getCallerIDErr)
			iamMock.On("ListAccountAliases", mock.Anything, &iam.ListAccountAliasesInput{}).Return(&iam.ListAccountAliasesOutput{
				AccountAliases: tc.accountAliases,
			}, tc.listAliasesErr)

			manage := &AwsConfigManage{
				cfg: aws.Config{
					Region: "us-somewhere-1",
				},
				iam: iamMock,
				sts: stsMock,
			}

			result, err := manage.GetDefaultConfig()
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestGetAllOrgAccounts(t *testing.T) {
	testCases := []struct {
		name            string
		listAccountsErr error
		expected        []AwsConfig
		expectErr       bool
	}{
		{
			name:            "returns all organization accounts",
			listAccountsErr: nil,
			expected: []AwsConfig{
				{
					Profile:   "name",
					AccountId: "id",
					Region:    "us-somewhere-1",
				},
			},
			expectErr: false,
		},
		{
			name:            "returns error when ListAccounts fails",
			listAccountsErr: errors.New("error"),
			expected:        nil,
			expectErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock AWS SDK calls
			orgMock := &mockaws.Orger{}
			orgMock.On("ListAccounts", mock.Anything, &organizations.ListAccountsInput{}).Return(&organizations.ListAccountsOutput{
				Accounts: []types.Account{{Id: aws.String("id"), Name: aws.String("name")}},
			}, tc.listAccountsErr)

			manage := &AwsConfigManage{
				cfg: aws.Config{Region: "us-somewhere-1"},
				org: orgMock,
			}

			result, err := manage.GetAllOrgAccounts()
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestGenerateConfig(t *testing.T) {
	testCases := []struct {
		name        string
		configInput AwsConfig
		expectPanic bool
	}{
		{
			name: "generates AWS config",
			configInput: AwsConfig{
				Profile:   "test-profile",
				AccountId: "test-account-id",
				Region:    "test-region",
				Endpoint:  "",
			},
			expectPanic: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectPanic {
				assert.Panics(t, func() { tc.configInput.GenerateConfig() })
			} else {
				assert.NotPanics(t, func() { tc.configInput.GenerateConfig() })
			}
		})
	}
}
