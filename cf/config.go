package cf

import (
	"context"
	"net/url"
	"path/filepath"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/songchenwen/cloudfront-invalidator/config"
)

var client *cloudfront.Client
var waiter *cloudfront.InvalidationCompletedWaiter

func Init() (err error) {
	c, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.AwsKeyId(), config.AwsKeySec(), "")),
		awsconfig.WithRegion(config.AwsRegion()),
	)
	if err != nil {
		return
	}
	client = cloudfront.NewFromConfig(c)
	waiter = cloudfront.NewInvalidationCompletedWaiter(client)
	return
}

func urls2paths(urls []string) (paths []string) {
	for _, u := range urls {
		parsed, err := url.Parse(u)
		if err == nil {
			u = parsed.Path
		}
		if filepath.IsAbs(u) {
			paths = append(paths, u)
		}
	}
	return
}
