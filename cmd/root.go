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
	cfgFile   string
	glHost    string
	tokenVar  string
	glToken   string
	semVer    string
	gitCommit string
	gitRef    string
	buildDate string

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
		glToken = os.Getenv(tokenVar)
		if len(glToken) == 0 {
			logrus.Fatal("GL_TOKEN is not set")
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

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".gitlab-tool" (without extension).
		viper.AddConfigPath(fmt.Sprintf("%s/.config/gitlab-tool", home))
		viper.SetConfigType("yaml")
		viper.SetConfigName("gitlab-tool")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	// }
}
