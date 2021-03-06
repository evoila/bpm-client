package collections

type StringSet struct {
	set map[string]bool
}

func NewStringSet() *StringSet {
	return &StringSet{make(map[string]bool)}
}

func (set *StringSet) Add(s string) bool {
	_, found := set.set[s]
	set.set[s] = true
	return !found
}

func (set *StringSet) Get(s string) bool {
	_, found := set.set[s]
	return found
}

func (set *StringSet) Remove(s string) {
	delete(set.set, s)
}
