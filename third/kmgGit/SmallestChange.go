package kmgGit
import (
	"github.com/bronze1man/kmg/kmgCmd"
	"bytes"
)
/*
// 返回最小变化的提交的名字
// 找到当前文件和commitId里面变化最小的版本
// localCommit 是当前文件的提交
// targetCommit 是需要寻找最小版本的commit
 无.git文件恢复
 git clone xxx ./tmp1
 mv ./tmp1/.git ./
 git checkout --orphan current
 kmg GitSmallestChange -local=current -target=master
 git checkout xx //返回的那个分支地址
*/
func (repo *Repository) MustSmallestChange(localCommit string,targetCommit string) string{
	commitList:=repo.MustGetAllParentCommitId(targetCommit)
	minNum:=2<<31
	minCommit:=""
	lineBreak:=[]byte("\n")
	for _,commitName:=range commitList{
		output:=kmgCmd.MustCombinedOutput("git diff "+localCommit +" "+commitName)
		num:=bytes.Count(output,lineBreak)
		if minNum>num{
			minNum = num
			minCommit = commitName
			if num==0{
				break
			}
		}
	}
	return minCommit
}