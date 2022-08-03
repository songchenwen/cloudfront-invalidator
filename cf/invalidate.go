package cf

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/google/uuid"
	"github.com/songchenwen/cloudfront-invalidator/config"
	"github.com/songchenwen/cloudfront-invalidator/utils"
)

const (
	httpTimeout = time.Second * 30
	waitTimeout = time.Minute * 10
)

func Invalidate(distribution string, urls []string, crawl bool) (id string, err error) {
	if client == nil {
		err = errors.New("cf client not initialized")
		return
	}

	paths := urls2paths(urls)
	paths = utils.Unique(paths)
	paths = checkPathsWithCache(distribution, paths)
	quantity := int32(len(paths))
	if quantity == 0 {
		return "No Need To Invalidate", nil
	}
	addPathToCache(distribution, paths)
	reference := uuid.New().String()
	req := &cloudfront.CreateInvalidationInput{
		DistributionId: &distribution,
		InvalidationBatch: &types.InvalidationBatch{
			CallerReference: &reference,
			Paths: &types.Paths{
				Quantity: &quantity,
				Items:    paths,
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()
	res, err := client.CreateInvalidation(ctx, req)
	if err != nil {
		log.Printf("invalidate failed %s ref %s, %v", distribution, reference, paths)
		return
	}
	id = *res.Invalidation.Id
	log.Printf("invalidate created %s ref %s id %s, %v", distribution, reference, id, paths)
	go func() {
		begin := time.Now()
		err = waiter.Wait(context.Background(), &cloudfront.GetInvalidationInput{
			DistributionId: &distribution,
			Id:             &id,
		}, waitTimeout)
		duration := time.Since(begin)
		if err != nil {
			log.Printf("invalidate wait err %s in %v %v", id, duration, err)
			return
		}
		log.Printf("invalidate complete %s in %v, %v", id, duration, paths)
		removePathFromCache(distribution, paths)
		if crawl {
			crawlUrls(urls)
		}
	}()
	return
}

func crawlUrls(urls []string) {
	if config.IsDebug() {
		log.Printf("craw urls %v", urls)
	}
	for _, u := range urls {
		go func(u string) {
			_, err := url.Parse(u)
			if err != nil && config.IsDebug() {
				log.Printf("url error %s, %v", u, err)
				return
			}
			_, err = http.Get(u)
			log.Printf("crawled url %s %v", u, err)
		}(u)
	}
}
