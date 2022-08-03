package cf

import "fmt"

func getCacheKey(distributionId string, path string) string {
	return fmt.Sprintf("%s.%s", distributionId, path)
}

func addPathToCache(distributionId string, paths []string) {
	for _, p := range paths {
		invalidationCache.Set(getCacheKey(distributionId, p), 1)
	}
}

func checkPathsWithCache(distributionId string, paths []string) (notCached []string) {
	for _, p := range paths {
		k := getCacheKey(distributionId, p)
		if _, has := invalidationCache.Get(k); !has {
			notCached = append(notCached, p)
		}
	}
	return
}

func removePathFromCache(distributionId string, paths []string) {
	for _, p := range paths {
		invalidationCache.Remove(getCacheKey(distributionId, p))
	}
}
