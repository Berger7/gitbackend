package gitbackend

import (
	"testing"
	"fmt"
)

const TestRepoPath = "./git-test.git"

func TestGetBranches(t *testing.T){
	repo,err := Init(TestRepoPath)
	if err != nil{
		fmt.Println(err)
	}
	branchNames,err := repo.getBranches()
	if err != nil {
		t.Errorf("get head branch error")
	}
	if len(branchNames) == 0 {
		t.Errorf("get branch name failed")
	}
	return
}

func TestFindDefaultBranch(t *testing.T) {
	repo,err := Init(TestRepoPath)
	if err != nil{
		t.Errorf("init repo error")
	}

	defaultBranch, err := repo.findDefaultBranch()
	if err != nil {
		t.Errorf("get findDefaultBranch error")
	}

	if defaultBranch != "master" {
		t.Errorf(fmt.Sprintf(""))
	}

	return
}

func TestSize(t *testing.T) {
	repo, err := Init(TestRepoPath)
	if err != nil {
		t.Errorf("init repo error")
	}

	size, err := repo.size()
	if err != nil {
		t.Errorf("get repo size error : %v", err)
	}

	if len(size) == 0 {
		t.Errorf("get repo size string err : %v", err)
	}

	return
}

func TestLog(t *testing.T)  {
	repo, err := Init(TestRepoPath)
	if err != nil {
		t.Errorf("init repo error")
	}

	commits, err := repo.log(10, 0, "master", false, true);
	if err != nil {
		t.Errorf("get commits error : %v", err)
	}
	if len(commits) == 0 {
		t.Errorf("get commits failed")
	}
}

func TestRootRef(t *testing.T) {
	repo, err := Init(TestRepoPath)
	if err != nil {
		t.Errorf("init repo error")
	}

	refName, err := repo.findDefaultBranch()
	if err != nil {
		t.Errorf("get refName error : %v", err)
	}
	if len(refName) == 0 {
		t.Errorf("get refName is empty")
	}
}

func TestCommitsBetween(t *testing.T) {
	repo, err := Init(TestRepoPath)
	if err != nil {
		t.Errorf("init repo error")
	}

	fromRefName := "testA"
	toRefName := "testB"
	commits, err := repo.commitsBetween(fromRefName, toRefName)
	if err != nil {
		t.Errorf("get commits between %s..%s", fromRefName, toRefName)
	}
	if len(commits) != 2 {
		t.Errorf("commits between %s..%s should be 2", fromRefName, toRefName)
	}
}