package gitlab

import "github.com/xanzy/go-gitlab"

type rawClient interface {
	GetProjects() projectsService
	GetMergeRequests() mergeRequestsService
	GetIssues() issuesService
	//GetTags() tagsService
	SetBaseURL(baseUrl string) error
}

type projectsService interface {
	GetProject(pid interface{}, options ...gitlab.OptionFunc) (*gitlab.Project, *gitlab.Response, error)
	CreateProject(opt *gitlab.CreateProjectOptions, options ...gitlab.OptionFunc) (*gitlab.Project, *gitlab.Response, error)
}

type issuesService interface {
	ListProjectIssues(pid interface{}, opt *gitlab.ListProjectIssuesOptions, options ...gitlab.OptionFunc) ([]*gitlab.Issue, *gitlab.Response, error)
}

type mergeRequestsService interface {
	ListMergeRequests(opt *gitlab.ListMergeRequestsOptions, options ...gitlab.OptionFunc) ([]*gitlab.MergeRequest, *gitlab.Response, error)
	CreateMergeRequest(pid interface{}, opt *gitlab.CreateMergeRequestOptions, options ...gitlab.OptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error)
}

//type tagsService interface {
//	CreateTag(pid interface{}, opt *gitlab.CreateTagOptions, options ...gitlab.OptionFunc) (*gitlab.Tag, *gitlab.Response, error)
//}

type RawClient struct {
	*gitlab.Client
}

func (r *RawClient) GetProjects() projectsService {
	return projectsService(r.Projects)
}

func (r *RawClient) GetIssues() issuesService {
	return issuesService(r.Issues)
}

func (r *RawClient) GetMergeRequests() mergeRequestsService {
	return mergeRequestsService(r.MergeRequests)
}

//func (r *RawClient) GetTags() tagsService {
//	return tagsService(r.Tags)
//}
