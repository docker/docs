package common

// ScopeSet is a set of scope strings.
type ScopeSet map[string]struct{}

// NewScopeSet makes a new set of the given values.
func NewScopeSet(vals ...string) ScopeSet {
	ss := make(ScopeSet, len(vals))
	for _, val := range vals {
		ss[val] = struct{}{}
	}

	return ss
}

// Contains checks for the existence of the given value in this scope set.
func (ss ScopeSet) Contains(val string) bool {
	_, hasVal := ss[val]
	return hasVal
}

// Add inserts the given value into this scope set.
func (ss ScopeSet) Add(val string) {
	ss[val] = struct{}{}
}

// HasAccess checks for the given access level in this scope set.
func (ss ScopeSet) HasAccess(level string) bool {
	return ss.Contains("*") || ss.Contains(level)
}

// Union creates a new scope set which is the union
// of this scope set and the given other scope set.
func (ss ScopeSet) Union(other ScopeSet) ScopeSet {
	// Short circuits.
	if ss.Contains("*") || other.Contains("*") {
		return NewScopeSet("*")
	}

	unioned := make(ScopeSet, len(ss)+len(other))

	// Add keys from this set.
	for key := range ss {
		unioned.Add(key)
	}

	// Add keys from the other set.
	for key := range other {
		unioned.Add(key)
	}

	return unioned
}

// Intersect returns a new scope set which is the intersection of this one and the other one.
func (ss ScopeSet) Intersect(other ScopeSet) ScopeSet {
	// Short circuits.
	if ss.Contains("*") {
		return other.Copy()
	}
	if other.Contains("*") {
		return ss.Copy()
	}

	// Ensure that `ss` is the smaller set.
	if len(ss) > len(other) {
		ss, other = other, ss
	}

	intersected := make(ScopeSet, len(ss))

	for key := range ss {
		if other.Contains(key) {
			intersected.Add(key)
		}
	}

	return intersected
}

// Copy makes a copy of this scope set.
func (ss ScopeSet) Copy() ScopeSet {
	// Short circuit.
	if ss.Contains("*") {
		return NewScopeSet("*")
	}

	copied := make(ScopeSet, len(ss))

	for key := range ss {
		copied.Add(key)
	}

	return copied
}
