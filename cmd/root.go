package cmd

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	git "github.com/go-git/go-git/v5"
	gl "github.com/maahsome/gitlab-go"
	"github.com/maahsome/gitlab-tool/cmd/config"
	"github.com/maahsome/gitlab-tool/cmd/objects"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	giturls "github.com/whilp/git-urls"

	"github.com/spf13/viper"
)

var (
	cfgFile       string
	updateConfig  bool
	workDir       string
	glHost        string
	tokenVar      string
	glToken       string
	semVer        string
	gitCommit     string
	gitRef        string
	buildDate     string
	cwdProjectID  int
	cwdGitlabHost string
	gitClient     gl.GitlabClient

	semVerReg = regexp.MustCompile(`(v[0-9]+\.[0-9]+\.[0-9]+).*`)

	c = &config.Config{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitlab-tool",
	Short: "A collection of REST API Calls",
	Long: `A cli tool to free you from the browser as much as possible.
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		c.VersionDetail.SemVer = semVer
		c.VersionDetail.BuildDate = buildDate
		c.VersionDetail.GitCommit = gitCommit
		c.VersionDetail.GitRef = gitRef
		c.VersionJSON = fmt.Sprintf("{\"SemVer\": \"%s\", \"BuildDate\": \"%s\", \"GitCommit\": \"%s\", \"GitRef\": \"%s\"}", semVer, buildDate, gitCommit, gitRef)
		// if updateConfig {
		// 	if len(glHost) > 0 {
		// 		viper.Set("gitlabhost", glHost)
		// 		verr := viper.WriteConfig()
		// 		if verr != nil {
		// 			logrus.WithError(verr).Info("Failed to write config")
		// 		} else {
		// 			logrus.Info("Successfully saved gitlab-host (%s) to config.yaml\n", glHost)
		// 		}
		// 	}
		// 	if len(tokenVar) > 0 {
		// 		viper.Set("tokenvar", tokenVar)
		// 		verr := viper.WriteConfig()
		// 		if verr != nil {
		// 			logrus.WithError(verr).Info("Failed to write config")
		// 		} else {
		// 			logrus.Info("Successfully saved token-var (%s) to config.yaml\n", tokenVar)
		// 		}
		// 	}
		// }

		getCurrentWorkingDirGitInfo()

		// glDefaultHost := ""
		// glDefaultEnvVar := ""
		// currentHost := viper.GetString("currentHost")
		// var hostList objects.HostList
		// err := viper.UnmarshalKey("hosts", &hostList)
		// if err != nil {
		// 	logrus.Fatal("Error unmarshalling...")
		// }
		// for _, v := range hostList {
		// 	if strings.EqualFold(v.Host, cwdGitlabHost) {
		// 		glHost = v.Host
		// 		tokenVar = v.EnvVar
		// 	}
		// 	if v.Host == currentHost {
		// 		glDefaultHost = v.Host
		// 		glDefaultEnvVar = v.EnvVar
		// 	}
		// }
		// if len(glHost) == 0 {
		// 	glHost = glDefaultHost
		// 	tokenVar = glDefaultEnvVar
		// }
		// // glHostFromConfig := viper.GetString("gitlabhost")
		// // if len(glHostFromConfig) > 0 {
		// // 	glHost = glHostFromConfig
		// // }
		// // tokenVarFromConfig := viper.GetString("tokenvar")
		// // if len(tokenVarFromConfig) > 0 {
		// // 	tokenVar = tokenVarFromConfig
		// // }
		// // glToken = os.Getenv(tokenVar)
		// // if len(glToken) == 0 {
		// // 	logrus.Fatal(fmt.Sprintf("%s ENV VAR does not have a value", tokenVar))
		// // }

		// gitClient = gl.New(glHost, "", glToken)

		if c.OutputFormat != "" {
			c.OutputFormat = strings.ToLower(c.OutputFormat)
			switch c.OutputFormat {
			case "json", "gron", "yaml", "text", "table", "raw":
				break
			default:
				fmt.Println("Valid options for -o are [json|gron|text|table|yaml|raw]")
				os.Exit(1)
			}
		}
	},
}

func getCurrentWorkingDirGitInfo() {

	cwdProjectID = 0
	cwdGitlabHost = ""

	// Are we in a GIT working directory, if so, collect the host/projectid
	workDir, werr := os.Getwd()
	if werr != nil {
		logrus.Fatal("Failed to get the current working directory?  That is odd.")
	}

	gitDir := fmt.Sprintf("%s/.git", workDir)
	if stat, err := os.Stat(gitDir); err == nil && !stat.IsDir() {
		realDir, rerr := os.ReadFile(gitDir)
		if rerr != nil {
			logrus.Fatal("Failed to read the worktree gitdir...")
		}
		workDir = strings.Split(strings.TrimSpace(strings.TrimPrefix(string(realDir[:]), "gitdir: ")), ".git")[0]
	}

	repo, rerr := git.PlainOpen(workDir)
	if rerr != nil {
		logrus.Fatal("Error retrieving git info")
	}
	repoConfig, rcerr := repo.Config()
	if rcerr != nil {
		logrus.Fatal("Error getting Config")
	}
	// fmt.Printf("%#v\n", repoConfig)
	pURLs, _ := giturls.Parse(repoConfig.Remotes["origin"].URLs[0])
	glSlug := strings.TrimPrefix(strings.TrimSuffix(pURLs.EscapedPath(), ".git"), "/")
	glSlug = url.PathEscape(glSlug)

	cwdGitlabHost = pURLs.Host

	glDefaultHost := ""
	glDefaultEnvVar := ""
	currentHost := viper.GetString("currentHost")
	var hostList objects.HostList
	err := viper.UnmarshalKey("hosts", &hostList)
	if err != nil {
		logrus.Fatal("Error unmarshalling...")
	}
	for _, v := range hostList {
		if strings.EqualFold(v.Host, cwdGitlabHost) {
			glHost = v.Host
			tokenVar = v.EnvVar
		}
		if v.Host == currentHost {
			glDefaultHost = v.Host
			glDefaultEnvVar = v.EnvVar
		}
	}
	if len(glHost) == 0 {
		glHost = glDefaultHost
		tokenVar = glDefaultEnvVar
	}
	glToken = os.Getenv(tokenVar)
	if len(glToken) == 0 {
		logrus.Fatal(fmt.Sprintf("%s ENV VAR does not have a value", tokenVar))
	}
	// fmt.Printf("Host: %s, Token: %s\n", glHost, glToken)
	gitClient = gl.New(glHost, "", glToken)

	projectID, pierr := gitClient.GetProjectID(glSlug)
	if pierr != nil {
		logrus.Fatal("Could not get ProjectID from Slug", glSlug)
	}

	cwdProjectID = projectID

}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gitlab-tool.yaml)")
	rootCmd.PersistentFlags().StringVar(&tokenVar, "token-variable", "GL_TOKEN", "Specify the ENV variable containing the gitlab PAT")
	rootCmd.PersistentFlags().StringVar(&glHost, "gitlab-host", "gitlab.com", "Base gitlab host")
	rootCmd.PersistentFlags().StringVarP(&c.OutputFormat, "output", "o", "", "Set an output format: json, text, yaml, gron")
	// rootCmd.PersistentFlags().BoolVar(&updateConfig, "update-config", false, "Update the config file with --gitlab-host and/or --token-variable")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		workDir = fmt.Sprintf("%s/.config/gitlab-tool", home)
		if _, err := os.Stat(workDir); err != nil {
			if os.IsNotExist(err) {
				mkerr := os.MkdirAll(workDir, os.ModePerm)
				if mkerr != nil {
					logrus.Fatal("Error creating ~/.config/gitlab-tool directory", mkerr)
				}
			}
		}
		if stat, err := os.Stat(workDir); err == nil && stat.IsDir() {
			configFile := fmt.Sprintf("%s/%s", workDir, "config.yaml")
			createRestrictedConfigFile(configFile)
			viper.SetConfigFile(configFile)
		} else {
			logrus.Info("The ~/.config/gitlab-tool path is a file and not a directory, please remove the 'gitlab-tool' file.")
			os.Exit(1)
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		logrus.Warn("Failed to read viper config file.")
	}
}

func createRestrictedConfigFile(fileName string) {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			file, ferr := os.Create(fileName)
			if ferr != nil {
				logrus.Info("Unable to create the configfile.")
				os.Exit(1)
			}
			mode := int(0600)
			if cherr := file.Chmod(os.FileMode(mode)); cherr != nil {
				logrus.Info("Chmod for config file failed, please set the mode to 0600.")
			}
		}
	}
}
