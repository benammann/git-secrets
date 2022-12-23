package gcp

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func ListProjects(ctx context.Context) (*ProjectsList, error) {
	crm, errCrm := cloudresourcemanager.NewService(ctx)
	if errCrm != nil {
		return nil, fmt.Errorf("could not init cloudresourcemanager: %s", errCrm.Error())
	}
	projects, errListProjects := crm.Projects.List().Do()
	if errListProjects != nil {
		return nil, fmt.Errorf("could not list gcp projects via cloudresourcemanager: %s", errCrm.Error())
	}

	return NewProjectsList(projects.Projects), nil
}

type ProjectsList struct {
	projects []*cloudresourcemanager.Project
}

func NewProjectsList(projects []*cloudresourcemanager.Project) *ProjectsList {
	return &ProjectsList{
		projects: projects,
	}
}

func (l *ProjectsList) GetProjects() []*cloudresourcemanager.Project {
	return l.projects
}

func (l *ProjectsList) GetProjectById(id string) *cloudresourcemanager.Project {
	for _, project := range l.projects {
		if project.ProjectId == id {
			return project
		}
	}
	return nil
}

func (l *ProjectsList) ProjectIds() (ids []string) {
	for _, project := range l.projects {
		ids = append(ids, project.ProjectId)
	}
	return ids
}

func (l *ProjectsList) AskProject() (*cloudresourcemanager.Project, error) {
	projectId := ""
	prompt := &survey.Select{
		Message: "GCP Project:",
		Options: l.ProjectIds(),
		Description: func(value string, index int) string {
			project := l.GetProjectById(value)
			if project != nil {
				return project.Name
			}
			return ""
		},
	}
	errAsk := survey.AskOne(prompt, &projectId)
	return l.GetProjectById(projectId), errAsk
}