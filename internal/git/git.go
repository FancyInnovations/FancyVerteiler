package git

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

type Service struct {
	cachedCommit  string
	cachedMessage string
}

func New() *Service {
	return &Service{
		cachedCommit:  "unknown",
		cachedMessage: "unknown",
	}
}

func (s *Service) Setup() error {
	//sha, err := fetchLatestCommitSHA()
	//if err != nil {
	//	return err
	//}
	//
	//msg, err := fetchLatestCommitMessage()
	//if err != nil {
	//	return err
	//}
	//
	//s.cachedCommit = strings.TrimSpace(sha)
	//s.cachedMessage = strings.TrimSpace(msg)

	s.cachedCommit = os.Getenv("COMMIT_SHA")
	s.cachedMessage = os.Getenv("COMMIT_MESSAGE")

	return nil
}

func (s *Service) CommitSHA() string {
	return s.cachedCommit
}

func (s *Service) CommitMessage() string {
	return s.cachedMessage
}

func fetchLatestCommitSHA() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Dir = os.Getenv("GITHUB_WORKSPACE")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return out.String(), nil
}

func fetchLatestCommitMessage() (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=%B")
	cmd.Dir = os.Getenv("GITHUB_WORKSPACE")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}
