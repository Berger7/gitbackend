package gitbackend

import (
	"github.com/libgit2/git2go"
	"strings"
	"fmt"
	"regexp"
	"path/filepath"
	"os/exec"
	"io/ioutil"
	"strconv"
	"bufio"
	"io"
)

type RawRepository struct {
	Repository 		*git.Repository
	RepoPath 		string
	RepoName		string
}

const emptyBranchName = ""
const tabSplit = "	"
const SIZE64 = 64
const PATHSPLIT = "/"

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

func (r *RawRepository) rootRef() (refName string, err error){
	refName,err = r.findDefaultBranch()
	return
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

	originRef := strings.Split(head.Name(), PATHSPLIT)
	if len(originRef) == 3 {
		return originRef[2], nil
	}
	return originRef[0],nil
}

func (r *RawRepository) hasCommits() (bool bool, err error) {
	repoIsEmpty,err := r.isEmpty();
	if err != nil || repoIsEmpty {
		return false, err
	}
	return true, err
}

func (r *RawRepository) isEmpty() (bool bool, err error)  {
	repoIsEmpty,err := r.Repository.IsEmpty()
	if err != nil || repoIsEmpty {
		return true, err
	}
	return false, err
}

func (r *RawRepository) isBare() (bool bool, err error)  {
	repoIsBare := r.Repository.IsBare()
	if repoIsBare {
		return true, err
	}
	return false, err
}

func (r *RawRepository) isExist() (bool bool, err error)  {
	if r.Repository == nil {
		return false, err
	}
	return true, err
}

func (r *RawRepository) archiveFilePath(name string, storagePath string, format string) (path string, err error) {
	if len(format) == 0 {
		format = "tar.gz"
	}

	extension := format

	switch format {
	case "tar.bz2", "tbz", "tbz2", "tb2", "bz2":
		extension = "tar.bz2"
	case "tar":
		extension = "tar"
	case "zip":
		extension = "zip"
	default:
		extension = "tar.gz"
	}

	fileName := fmt.Sprintf("%s.%s",name, extension)
	path = filepath.Join(storagePath, name, fileName)
	return path, err
}

func (r *RawRepository) size() (size string, err error) {
	cmd := exec.Command("du", "-sk", r.Repository.Path())
	stdout, err := cmd.StdoutPipe()

	defer stdout.Close()

	if err = cmd.Start(); err != nil {
		return
	}

	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return
	}

	if err = cmd.Wait(); err != nil {
		return
	}

	outPutString := string(opBytes)
	rawSize := strings.Split(outPutString, tabSplit)[0]
	floatSize,err := strconv.ParseFloat(rawSize, SIZE64)
	if err != nil {
		return
	}

	size = fmt.Sprintf("%.2f", floatSize / 1024)

	fmt.Printf(size)
	return
}

func (r *RawRepository) log(limit int, offset int, refName string,
	followFlag bool, skipMerges bool) (commits []git.Commit,err error) {
	oid,err := r.shaFromRev(refName)
	if err != nil {
		return
	}

	opCMD := exec.Command("git",
		fmt.Sprintf("--git-dir=%s", r.Repository.Path()),
		"log",
		"-n", fmt.Sprintf("%d", limit),
		"--format=%H",
		fmt.Sprintf("--skip=%s", offset))

	if followFlag {
		opCMD.Args = append(opCMD.Args, "--follow")
	}
	if skipMerges {
		opCMD.Args = append(opCMD.Args, "--no-merges")
	}

	opCMD.Args = append(opCMD.Args, fmt.Sprintf("%v", oid))

	stdout,err := opCMD.StdoutPipe()
	opCMD.Start()

	reader := bufio.NewReader(stdout)

	for {
		line, getCommitIDErr := reader.ReadString('\n')
		commitID := strings.TrimSpace(line)
		if getCommitIDErr != nil || io.EOF == getCommitIDErr {
			break
		}
		obj, getObjectErr := r.Repository.RevparseSingle(commitID)
		if getObjectErr != nil {
			return commits, getObjectErr
		}
		commit, getCommitErr := r.Repository.LookupCommit(obj.Id())
		if getCommitErr != nil {
			return commits, getCommitErr
		}
		commits = append(commits, *commit)
	}
	return
}

func (r *RawRepository) shaFromRev(rev string) (oid *git.Oid, err error) {
	obj,err := r.Repository.RevparseSingle(rev)
	if err != nil {
		return
	}
	return obj.Id(), nil
}