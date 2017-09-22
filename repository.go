package gitbackend

import (
	"github.com/libgit2/git2go"
	"strings"
	"fmt"
	"regexp"
)

type RawRepository struct {
	Repository 		*git.Repository
	RepoPath 		string
	RepoName		string
}

const emptyBranchName = ""

var RepositoryClient RawRepository

// Init create raw repository client
// full repo path must be the absolute _bare_ path, e.g.
// /path/to/your-repo.git
func Init(fullRepoPath string) (*RawRepository, error) {
	if len(fullRepoPath) == 0 {
		return nil, fmt.Errorf("Repository path can not be empty")
	}

	repo, err := git.OpenRepository(fullRepoPath)
	if err != nil {
		return nil, err
	}

	RepositoryClient := &RawRepository{
		Repository: repo,
		RepoPath: fullRepoPath,
		RepoName: strings.Split(fullRepoPath, "/")[1],
	}

	return RepositoryClient,nil
}

func (r *RawRepository) rootRef() {

}

// getTags return all tags from repository
func (r *RawRepository) getTags() ([]string, error){
	iter, err := r.Repository.NewReferenceNameIterator()
	if err != nil {
		return nil, err
	}
	var tags []string

	ref, err := iter.Next()
	tagRegex := regexp.MustCompile("^refs/tags/.*")

	for err == nil {
		if tagRegex.FindStringSubmatch(ref) != nil {
			tags = append(tags, strings.Split(ref, "refs/tags/")[1])
		}
		ref, err = iter.Next()
	}
	return tags,err
}

// getBranches return all branches from repository
func (r *RawRepository) getBranches() ([]string, error){
	iter, err := r.Repository.NewBranchIterator(git.BranchAll)
	if err != nil {
		return nil, err
	}

	var branches []string

	branch, _, err := iter.Next()

	for err == nil {
		name, _ := branch.Name()
		branches = append(branches, name)
		branch, _, err = iter.Next()
	}

	return branches, nil
}

// Discovers the default branch based on the repository's available branches
// - If no branches are present, returns nil
// - If one branch is present, returns its name
// - If two or more branches are present, returns current HEAD or master or first branch
// HEAD > master > first branch
func (r *RawRepository) findDefaultBranch() (branchName string, err error){
	branchNames, err := r.getBranches()
	if err != nil {
		return emptyBranchName, err
	}

	if len(branchNames) == 0 {
		return emptyBranchName, nil
	}

	if len(branchNames) == 1 {
		return branchNames[0], nil
	}

	headRef, err := r.getHeadBranch()
	if err != nil {
		return emptyBranchName, err
	}

	if len(headRef) != 0 {
		return headRef, nil
	}

	for _,tmpBranchName := range branchNames {
		if tmpBranchName == "master" {
			return "master", nil
		}
	}

	return branchNames[0], nil
}

func (r *RawRepository) getHeadBranch() (branchName string,err error){
	head,err := r.Repository.Head()
	if err != nil || head == nil {
		return emptyBranchName, err
	}

	originRef := strings.Split(head.Name(), "/")
	if len(originRef) == 3 {
		return originRef[2], nil
	}
	return originRef[0],nil
}
