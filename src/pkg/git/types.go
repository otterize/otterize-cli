package git

import "fmt"

type LocalGitInformation struct {
	Commit       string
	OriginUrl    string
	RelativePath string
}

func (l *LocalGitInformation) Print() {
	fmt.Printf("Git origin URL: %s\n", l.OriginUrl)
	fmt.Printf("Git commit ID: %s\n", l.Commit)
	fmt.Printf("Relative path to git root: %s\n", l.RelativePath)
}
