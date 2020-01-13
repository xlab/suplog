package blob

import (
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

// NewBlobID returns a pseudo-randomly generated ULID -
// Universally Unique Lexicographically Sortable Identifier
// - see https://github.com/ulid/spec
func NewBlobID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), globalRand).String()
}

//nolint:gochecknoglobals
var globalRand = rand.New(&lockedSource{
	src: rand.NewSource(time.Now().UnixNano()),
})

// lockedSource provides rand.Source with a mutex to avoid races.
type lockedSource struct {
	lk  sync.Mutex
	src rand.Source
}

func (r *lockedSource) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()

	return
}

func (r *lockedSource) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}
