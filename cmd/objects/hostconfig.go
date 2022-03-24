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

type HostList []HostConfig
type HostConfig struct {
	Name   string `mapstructure:"name"`
	Host   string `mapstructure:"host"`
	EnvVar string `mapstructure:"envvar"`
}

// ToJSON - Write the output as JSON
func (hl *HostList) ToJSON() string {
	hlJSON, err := json.MarshalIndent(hl, "", "  ")
	if err != nil {
		logrus.WithError(err).Error("Error extracting JSON")
		return ""
	}
	return string(hlJSON[:])
}

func (hl *HostList) ToGRON() string {
	hlJSON, err := json.MarshalIndent(hl, "", "  ")
	if err != nil {
		logrus.WithError(err).Error("Error extracting JSON for GRON")
	}
	subReader := strings.NewReader(string(hlJSON[:]))
	subValues := &bytes.Buffer{}
	ges := gron.NewGron(subReader, subValues)
	ges.SetMonochrome(false)
	if serr := ges.ToGron(); serr != nil {
		logrus.WithError(serr).Error("Problem generating GRON syntax")
		return ""
	}
	return string(subValues.Bytes())
}

func (hl *HostList) ToYAML() string {
	hlYAML, err := yaml.Marshal(hl)
	if err != nil {
		logrus.WithError(err).Error("Error extracting YAML")
		return ""
	}
	return string(hlYAML[:])
}

func (hl *HostList) ToTEXT(noHeaders bool) string {
	buf, row := new(bytes.Buffer), make([]string, 0)

	// ************************** TableWriter ******************************
	table := tablewriter.NewWriter(buf)
	if !noHeaders {
		table.SetHeader([]string{"NAME", "HOST", "ENVVAR"})
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
	for _, v := range *hl {
		row = []string{
			v.Name,
			v.Host,
			v.EnvVar,
		}
		table.Append(row)
	}

	table.Render()

	return buf.String()

}
