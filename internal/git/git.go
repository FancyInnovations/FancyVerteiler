package git

type Service struct {
	githubRepoURL string
	cachedCommit  string
	cachedMessage string
}

func New(githubRepoURL, sha, message string) *Service {
	return &Service{
		githubRepoURL: githubRepoURL,
		cachedCommit:  sha,
		cachedMessage: message,
	}
}

func (s *Service) GitHubRepoURL() string {
	return s.githubRepoURL
}

func (s *Service) CommitSHA() string {
	return s.cachedCommit
}

func (s *Service) CommitURL() string {
	return s.githubRepoURL + "/commit/" + s.cachedCommit
}

func (s *Service) CommitMessage() string {
	return s.cachedMessage
}
