package objects

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/maahsome/gron"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// ~/.config/gitlab-tool/config.yaml
// OLD
// currenthost: alteryx-private
// gitlabhost: git.alteryx.com
// hosts:
// - envvar: GLA_TOKEN
//   host: git.alteryx.com
//   name: alteryx-private
// - envvar: GL_TOKEN
//   host: gitlab.com
//   name: alteryx-public
// tokenvar: GLA_TOKEN
// NEW
// configs:
// - directory: /Users/christopher.maahs/dev/futurama
//   envvar: GLA_TOKEN
//   group: futurama
//   host: git.alteryx.com

type ConfigList []DirectoryConfig
type DirectoryConfig struct {
	Directory string `mapstructure:"directory"`
	EnvVar    string `mapstructure:"envvar"`
	Group     string `mapstructure:"group"`
	Host      string `mapstructure:"host"`
}

// ToJSON - Write the output as JSON
func (cl *ConfigList) ToJSON() string {
	clJSON, err := json.MarshalIndent(cl, "", "  ")
	if err != nil {
		logrus.WithError(err).Error("Error extracting JSON")
		return ""
	}
	return string(clJSON[:])
}

func (cl *ConfigList) ToGRON() string {
	clJSON, err := json.MarshalIndent(cl, "", "  ")
	if err != nil {
		logrus.WithError(err).Error("Error extracting JSON for GRON")
	}
	subReader := strings.NewReader(string(clJSON[:]))
	subValues := &bytes.Buffer{}
	ges := gron.NewGron(subReader, subValues)
	ges.SetMonochrome(false)
	if serr := ges.ToGron(); serr != nil {
		logrus.WithError(serr).Error("Problem generating GRON syntax")
		return ""
	}
	return string(subValues.Bytes())
}

func (cl *ConfigList) ToYAML() string {
	clYAML, err := yaml.Marshal(cl)
	if err != nil {
		logrus.WithError(err).Error("Error extracting YAML")
		return ""
	}
	return string(clYAML[:])
}

func (cl *ConfigList) ToTEXT(noHeaders bool) string {
	buf, row := new(bytes.Buffer), make([]string, 0)

	// ************************** TableWriter ******************************
	table := tablewriter.NewWriter(buf)
	if !noHeaders {
		table.SetHeader([]string{"DIRECTORY", "ENVVAR", "GROUP", "HOST"})
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	}

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	// for i=0; i<=len(hl); i++ {
	for _, v := range *cl {
		row = []string{
			v.Directory,
			v.EnvVar,
			v.Group,
			v.Host,
		}
		table.Append(row)
	}

	table.Render()

	return buf.String()

}
