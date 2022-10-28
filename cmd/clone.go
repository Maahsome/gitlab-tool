package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	git "github.com/go-git/go-git/v5"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type (
	projectAnswer struct {
		ProjectName string `survey:"projectname"` // or you can tag fields to match a specific name
	}
)

var (
	projectAnswers = &projectAnswer{}
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a gitlab project that is related to the config and pwd",
	Long: `EXAMPLE:
> gitlab-tool clone <project-name>

* project cloned

`,
	Run: func(cmd *cobra.Command, args []string) {
		if inProject {
			fmt.Println("You are currently in a PROJECT directory, cannot clone here")
		} else {
			if len(args) > 0 {
				cloneProject(args[0])
			} else {
				// Let's prompt for available projects to clone
				projects, gerr := gitClient.GetGroupProjects(cwdGroupID)
				if gerr != nil {
					logrus.WithError(gerr).Error("Failed to fetch sub-groups")
				}
				var projectArray []string
				for _, v := range projects {
					if v.Namespace.ID == cwdGroupID {
						if !exists(v.Path) {
							projectArray = append(projectArray, v.Path)
						}
					}
					// fmt.Printf("%s -> %s\n", v.Path, v.FullPath)
				}

				// the questions to ask
				var projectSurvey = []*survey.Question{
					{
						Name: "projectname",
						Prompt: &survey.Select{
							Message: "Choose a project to clone:",
							Options: projectArray,
						},
					},
				}

				opts := survey.WithStdio(os.Stdin, os.Stderr, os.Stderr)

				if len(projectArray) > 0 {
					// perform the question
					if err := survey.Ask(projectSurvey, projectAnswers, opts); err != nil {
						logrus.Fatal("No projects selected")
					}
					// fmt.Printf("Selected Project: %s\n", projectAnswers.ProjectName)
					cloneProject(projectAnswers.ProjectName)
				} else {
					logrus.Info("All of the projects in this group have been cloned locally")
				}
			}
		}
	},
}

func cloneProject(name string) {

	topGroup := filepath.Base(configDir)
	subSlug := strings.TrimPrefix(currentWorkDir, configDir)
	prPath := fmt.Sprintf("%s%s/%s", topGroup, subSlug, name)
	prSlug := url.PathEscape(fmt.Sprintf("%s%s/%s", topGroup, subSlug, name))
	projectID, err := gitClient.GetProjectID(prSlug)
	if err != nil {
		logrus.WithError(err).Fatal("Error getting Project ID")
	}
	if projectID != 0 {
		// Does the directory already exist?  If so, no clone
		if exists(name) {
			logrus.Fatal("The project directory already exists, cannot clone")
		} else {
			if err := os.Mkdir(name, os.ModePerm); err != nil {
				logrus.WithError(err).Fatal("Failed to create the project directory")
			}
			directory := name
			url := fmt.Sprintf("https://%s/%s.git", cwdGitlabHost, prPath)
			token := gitClient.GetProperty("Token")
			// r, err := git.PlainClone(directory, false, &git.CloneOptions{
			_, err := git.PlainClone(directory, false, &git.CloneOptions{
				URL: url,
				Auth: &gitHttp.BasicAuth{
					Username: "access_token",
					Password: token,
				},
				// Depth:    1,
				Progress: os.Stderr,
			})
			if err != nil {
				logrus.WithError(err).Fatal("Failed to clone the project")
			}
		}
	} else {
		logrus.Fatal(fmt.Sprintf("Failed to lookup the ProjectID for %s", prSlug))
	}
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}
