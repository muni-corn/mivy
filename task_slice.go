package mivy

// TaskSlice is a slice of tasks that can be sorted
type TaskSlice []Task

func (s TaskSlice) Len() int {
    return len(s)
}

func (s TaskSlice) Less(i, j int) bool {
    return s[i].IsLessThan(s[j])
}

func (s TaskSlice) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

func (s TaskSlice) Groups() []string {
    type void struct{}
    var yeet void

    set := make(map[string]void)

    for _, t := range s {
        set[t.Group] = yeet
    }

    keys := make([]string, len(set))
    i := 0
    for k := range set {
        keys[i] = k
    }

    return keys
}

