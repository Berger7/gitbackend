package gitbackend

import (
	"github.com/libgit2/git2go"
)

type Commit struct {
	repo *git.Repository
}
