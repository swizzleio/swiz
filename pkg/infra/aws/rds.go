package aws

import (
	"context"
	"fmt"
	"getswizzle.io/swiz/pkg/common"
	"getswizzle.io/swiz/pkg/infra/model"
	"getswizzle.io/swiz/pkg/network"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type Rds struct {
	client *rds.Client
}

func NewRds(cfg aws.Config) Rds {
	return Rds{
		client: rds.NewFromConfig(cfg),
	}
}

// ListInstances lists all the RDS instances and returns a mapping by unique id
func (e Rds) ListInstances() (map[string]model.TargetInstance, error) {
	instances := map[string]model.TargetInstance{}

	params := &rds.DescribeDBClustersInput{}
	paginator := rds.NewDescribeDBClustersPaginator(e.client, params)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		for _, cluster := range output.DBClusters {
			id := strOrEmpty(cluster.DBClusterArn)
			instances[id] = model.TargetInstance{
				Id:   id,
				Name: strOrEmpty(cluster.DBClusterIdentifier),
				Os:   e.getOs(cluster.Engine),
				Endpoints: []network.Endpoint{
					{
						Host: *cluster.Endpoint,
						Port: int(*cluster.Port),
						User: *cluster.MasterUsername,
					},
				},
			}
		}
	}

	return instances, nil
}

// getOs returns the database type
func (e Rds) getOs(engine *string) string {
	switch *engine {
	case "aurora-postgresql", "postgres":
		return common.OsPgSql
	case "aurora-mysql", "mysql":
		return common.OsMySql
	case "mariadb":
		return common.OsMariaDb
	default:
		return fmt.Sprintf("UNKNOWN:%s", *engine)
	}
}
