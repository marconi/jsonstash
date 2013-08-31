package bucket

import "errors"

/**
 * Bucket
 */

type Bucket struct {
	Values []string
	Index  map[string]int
}

func NewBucket() *Bucket {
	return &Bucket{Index: make(map[string]int)}
}

func (b *Bucket) Add(key string, val string) {
	// append the new value and remember its index
	values := append(b.Values, val)
	index := len(values) - 1
	b.Values = values

	// store the index
	b.Index[key] = index
}

func (b *Bucket) Get(key string) (string, error) {
	// check that the key exists and has valid index
	if index, ok := b.Index[key]; ok {
		if index < len(b.Values) {
			return b.Values[index], nil
		}
	}
	return "", errors.New("Invalid key")
}

func (b *Bucket) GetAll() []string {
	vals := make([]string, len(b.Values))
	copy(vals, b.Values)
	return vals
}

func (b *Bucket) Range(r *BucketRange) ([]string, error) {
	if r.Start < 0 {
		return nil, errors.New("Invalid start")
	}
	if r.Stop > len(b.Values) {
		return nil, errors.New("Invalid stop")
	} else if r.Stop == 0 {
		r.Stop = len(b.Values)
	}
	rVals := make([]string, (r.Stop - r.Start))
	copy(rVals, b.Values[r.Start:r.Stop])
	return rVals, nil
}

type BucketRange struct {
	Start int
	Stop  int
}

/*
 * Stash
 */

type Stash struct {
	Buckets map[string]*Bucket
}

func NewStash() *Stash {
	return &Stash{Buckets: make(map[string]*Bucket)}
}

func (s *Stash) Add(key string) *Bucket {
	// create new bucket and return it
	b := NewBucket()
	s.Buckets[key] = b
	return b
}

func (s *Stash) Get(key string) (*Bucket, error) {
	// check that the key is valid
	if b, ok := s.Buckets[key]; ok {
		return b, nil
	}
	return nil, errors.New("Invalid key")
}

func (s *Stash) GetBucketNames() []string {
	names := make([]string, 0, len(s.Buckets))
	for k, _ := range s.Buckets {
		names = append(names, k)
	}
	return names
}

func (s *Stash) GetBuckets() []*Bucket {
	buckets := make([]*Bucket, len(s.Buckets))
	for _, b := range s.Buckets {
		buckets = append(buckets, b)
	}
	return buckets
}
