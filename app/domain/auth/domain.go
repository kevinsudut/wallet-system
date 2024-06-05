package domainauth

import (
	"github.com/kevinsudut/wallet-system/pkg/lib/database"
	lrucache "github.com/kevinsudut/wallet-system/pkg/lib/lru-cache"
)

type domain struct {
	db    database.DatabaseItf
	cache lrucache.LRUCacheItf
}

func Init(db database.DatabaseItf) DomainItf {
	return &domain{
		db:    db,
		cache: lrucache.Init(),
	}
}
