package cleanpath

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

/*
RemoveGoPath makes a path relative to one of the src directories in the $GOPATH
environment variable. If $GOPATH is empty or the input path is not contained
within any of the src directories in $GOPATH, the original path is returned. If
the input path is contained within multiple of the src directories in $GOPATH,
it is made relative to the longest one of them.
*/
func RemoveGoPath(path string) string {
	dirs := filepath.SplitList(os.Getenv("GOPATH"))
	// Sort in decreasing order by length so the longest matching prefix is removed
	sort.Stable(longestFirst(dirs))
	for _, dir := range dirs {
		srcdir := filepath.Join(dir, "src")
		rel, err := filepath.Rel(srcdir, path)
		// filepath.Rel can traverse parent directories, don't want those
		if err == nil && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return rel
		}
	}
	return path
}

type longestFirst []string

func (strs longestFirst) Len() int           { return len(strs) }
func (strs longestFirst) Less(i, j int) bool { return len(strs[i]) > len(strs[j]) }
func (strs longestFirst) Swap(i, j int)      { strs[i], strs[j] = strs[j], strs[i] }
