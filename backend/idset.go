package backend

import "sort"

// IDSet is a set of IDs.
type IDSet map[ID]struct{}

// NewIDSet returns a new IDSet, populated with ids.
func NewIDSet(ids ...ID) IDSet {
	m := make(IDSet)
	for _, id := range ids {
		m[id] = struct{}{}
	}

	return m
}

// Has returns true iff id is contained in the set.
func (s IDSet) Has(id ID) bool {
	_, ok := s[id]
	return ok
}

// Insert adds id to the set.
func (s IDSet) Insert(id ID) {
	s[id] = struct{}{}
}

// Delete removes id from the set.
func (s IDSet) Delete(id ID) {
	delete(s, id)
}

// List returns a slice of all IDs in the set.
func (s IDSet) List() IDs {
	list := make(IDs, 0, len(s))
	for id := range s {
		list = append(list, id)
	}

	sort.Sort(list)

	return list
}

// Equals returns true iff s equals other.
func (s IDSet) Equals(other IDSet) bool {
	if len(s) != len(other) {
		return false
	}

	for id := range s {
		if _, ok := other[id]; !ok {
			return false
		}
	}

	for id := range other {
		if _, ok := s[id]; !ok {
			return false
		}
	}

	return true
}

func (s IDSet) String() string {
	str := s.List().String()
	if len(str) < 2 {
		return "{}"
	}

	return "{" + str[1:len(str)-1] + "}"
}
