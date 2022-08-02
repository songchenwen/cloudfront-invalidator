package server

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itchyny/gojq"
	"github.com/songchenwen/cloudfront-invalidator/cf"
	"github.com/songchenwen/cloudfront-invalidator/config"
	"github.com/songchenwen/cloudfront-invalidator/utils"
)

const (
	jqTimeout = time.Second * 30
)

func handleInvalidate(c *gin.Context) {
	urls := collectUrls(c)
	id, err := cf.Invalidate(c.Param("distribution"), urls, c.Query("crawl") != "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"distribution":  c.Param("distribution"),
		"urls":          urls,
		"invalidate_id": id,
	})
}

func collectUrls(c *gin.Context) (urls []string) {
	if queryUrls, ok := c.GetQueryArray("url"); ok {
		urls = decodeB64AndAddString(urls, queryUrls...)
	}
	if postUrls, ok := c.GetPostFormArray("url"); ok {
		urls = decodeB64AndAddString(urls, postUrls...)
	}

	urls = append(urls, collectUrlsFromBody(c)...)
	urls = utils.Unique(urls)
	return
}

func collectUrlsFromBody(c *gin.Context) (urls []string) {
	if c.Request.Method != http.MethodPost {
		return
	}
	jqQuerysB64, ok := c.GetQueryArray("jq")
	if !ok || len(jqQuerysB64) == 0 {
		return
	}
	var jqQuerys []string
	jqQuerys, jqQuerysB64 = jqsFromConfig(jqQuerysB64)
	jqQuerys = decodeB64AndAddString(jqQuerys, jqQuerysB64...)
	jqQuerys = utils.Unique(jqQuerys)
	if len(jqQuerys) == 0 {
		return
	}
	if config.IsDebug() {
		log.Printf("jq: %v", jqQuerys)
	}
	var body interface{}
	err := c.ShouldBindJSON(&body)
	if err != nil && config.IsDebug() {
		log.Printf("body json decode err %v", err)
		return
	}
	if config.IsDebug() {
		log.Printf("body: %v", body)
	}
	ctx, cancel := context.WithTimeout(context.Background(), jqTimeout)
	defer cancel()
	for _, query := range jqQuerys {
		q, err := gojq.Parse(query)
		if err != nil && config.IsDebug() {
			log.Printf("cannot parse jq query %s %v\n", query, err)
			continue
		}
		iter := q.RunWithContext(ctx, body)
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}
			if u, ok := v.(string); ok {
				urls = append(urls, strings.TrimSpace(u))
				continue
			}
		}
	}
	return
}

func jqsFromConfig(jqB64 []string) (jqs []string, jqB64Out []string) {
	for _, j := range jqB64 {
		js := config.GetJqs(j)
		if len(js) > 0 {
			jqs = append(jqs, js...)
		} else {
			jqB64Out = append(jqB64Out, j)
		}
	}
	return
}

func decodeB64AndAddString(urls []string, add ...string) (result []string) {
	result = urls
	for _, a := range add {
		u, err := base64.RawURLEncoding.DecodeString(a)
		if err == nil {
			result = append(result, strings.TrimSpace(string(u)))
		} else if config.IsDebug() {
			log.Printf("cannot decode base64 %s %v\n", a, err)
		}
	}
	return
}
