package gitbackend

import (
	"testing"
	"fmt"
)

const TestRepoPath = "/tmp/test.git"

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