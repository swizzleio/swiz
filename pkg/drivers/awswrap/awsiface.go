package awswrap

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

//go:generate mockery --name Iam --filename iam_mock.go --output ../../../mocks/ext/aws --outpkg mockaws
type Iamer interface {
	iam.ListAccountAliasesAPIClient
}

//go:generate mockery --name Sts --filename sts_mock.go --output ../../../mocks/ext/aws --outpkg mockaws
type Stser interface {
	GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

//go:generate mockery --name Org --filename org_mock.go --output ../../../mocks/ext/aws --outpkg mockaws
type Orger interface {
	organizations.ListAccountsAPIClient
}
