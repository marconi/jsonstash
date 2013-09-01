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

func (b *Bucket) Add(key string, val string) error {
	if _, ok := b.Index[key]; ok {
		return errors.New("Value with that key already exists.")
	}

	// append the new value and remember its index
	values := append(b.Values, val)
	index := len(values) - 1
	b.Values = values

	// store the index
	b.Index[key] = index
	return nil
}

func (b *Bucket) Get(key string) (string, error) {
	// check that the key exists and has valid index
	if index, ok := b.Index[key]; ok {
		if index < len(b.Values) {
			return b.Values[index], nil
		}
	}
	return "", errors.New("Invalid value")
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

func (b *Bucket) Update(key string, val string) error {
	// check that the key exists and has valid index
	if index, ok := b.Index[key]; ok {
		if index < len(b.Values) {
			b.Values[index] = val
			return nil
		}
	}
	return errors.New("Invalid value")
}

func (b *Bucket) Delete(key string) error {
	// check that the key exists and has valid index
	if i, ok := b.Index[key]; ok {
		if i < len(b.Values) {
			b.Values = append(b.Values[:i], b.Values[i+1:]...)
			delete(b.Index, key)
			return nil
		}
	}
	return errors.New("Invalid value")
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

func (s *Stash) Add(key string) (*Bucket, error) {
	// check that this is a new bucket
	if _, ok := s.Buckets[key]; ok {
		return nil, errors.New("Bucket already exists.")
	}

	// create new bucket and return it
	b := NewBucket()
	s.Buckets[key] = b
	return b, nil
}

func (s *Stash) Get(key string) (*Bucket, error) {
	// check that the key is valid
	if b, ok := s.Buckets[key]; ok {
		return b, nil
	}
	return nil, errors.New("Invalid bucket")
}

func (s *Stash) Delete(key string) error {
	// check that the key is valid
	if _, ok := s.Buckets[key]; ok {
		delete(s.Buckets, key)
		return nil
	}
	return errors.New("Invalid bucket")
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
