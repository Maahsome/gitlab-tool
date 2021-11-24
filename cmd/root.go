package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/maahsome/gitlab-tool/cmd/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	cfgFile      string
	updateConfig bool
	workDir      string
	glHost       string
	tokenVar     string
	glToken      string
	semVer       string
	gitCommit    string
	gitRef       string
	buildDate    string

	semVerReg = regexp.MustCompile(`(v[0-9]+\.[0-9]+\.[0-9]+).*`)

	c = &config.Config{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitlab-tool",
	Short: "A collection of REST API Calls",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		c.VersionDetail.SemVer = semVer
		c.VersionDetail.BuildDate = buildDate
		c.VersionDetail.GitCommit = gitCommit
		c.VersionDetail.GitRef = gitRef
		c.VersionJSON = fmt.Sprintf("{\"SemVer\": \"%s\", \"BuildDate\": \"%s\", \"GitCommit\": \"%s\", \"GitRef\": \"%s\"}", semVer, buildDate, gitCommit, gitRef)
		if updateConfig {
			if len(glHost) > 0 {
				viper.Set("gitlabhost", glHost)
				verr := viper.WriteConfig()
				if verr != nil {
					logrus.WithError(verr).Info("Failed to write config")
				} else {
					logrus.Info("Successfully saved gitlab-host (%s) to config.yaml\n", glHost)
				}
			}
			if len(tokenVar) > 0 {
				viper.Set("tokenvar", tokenVar)
				verr := viper.WriteConfig()
				if verr != nil {
					logrus.WithError(verr).Info("Failed to write config")
				} else {
					logrus.Info("Successfully saved token-var (%s) to config.yaml\n", tokenVar)
				}
			}
		}
		glHostFromConfig := viper.GetString("gitlabhost")
		if len(glHostFromConfig) > 0 {
			glHost = glHostFromConfig
		}
		tokenVarFromConfig := viper.GetString("tokenvar")
		if len(tokenVarFromConfig) > 0 {
			tokenVar = tokenVarFromConfig
		}
		glToken = os.Getenv(tokenVar)
		if len(glToken) == 0 {
			logrus.Fatal(fmt.Sprintf("%s ENV VAR does not have a value", tokenVar))
		}

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
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gitlab-tool.yaml)")
	rootCmd.PersistentFlags().StringVar(&tokenVar, "token-variable", "GL_TOKEN", "Specify the ENV variable containing the gitlab PAT")
	rootCmd.PersistentFlags().StringVar(&glHost, "gitlab-host", "gitlab.com", "Base gitlab host")
	rootCmd.PersistentFlags().StringVarP(&c.OutputFormat, "output", "o", "", "Set an output format: json, text, yaml, gron")
	rootCmd.PersistentFlags().BoolVar(&updateConfig, "update-config", false, "Update the config file with --gitlab-host and/or --token-variable")
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
