package objects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/maahsome/gron"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type GitListing []GitListItem

type GitListItem struct {
	StatBlock string `json:"stat_block"`
	Path      string `json:"path"`
	ID        int    `json:"id"`
}

// ToJSON - Write the output as JSON
func (gl *GitListing) ToJSON() string {
	glJSON, err := json.MarshalIndent(gl, "", "  ")
	if err != nil {
		logrus.WithError(err).Error("Error extracting JSON")
		return ""
	}
	return string(glJSON[:])
}

func (gl *GitListing) ToGRON() string {
	glJSON, err := json.MarshalIndent(gl, "", "  ")
	if err != nil {
		logrus.WithError(err).Error("Error extracting JSON for GRON")
	}
	subReader := strings.NewReader(string(glJSON[:]))
	subValues := &bytes.Buffer{}
	ges := gron.NewGron(subReader, subValues)
	ges.SetMonochrome(false)
	if serr := ges.ToGron(); serr != nil {
		logrus.WithError(serr).Error("Problem generating GRON syntax")
		return ""
	}
	return string(subValues.Bytes())
}

func (gl *GitListing) ToYAML() string {
	glYAML, err := yaml.Marshal(gl)
	if err != nil {
		logrus.WithError(err).Error("Error extracting YAML")
		return ""
	}
	return string(glYAML[:])
}

func (gl *GitListing) ToTEXT(noHeaders bool, showid bool) string {
	buf, row := new(bytes.Buffer), make([]string, 0)

	// ************************** TableWriter ******************************
	table := tablewriter.NewWriter(buf)
	if !noHeaders {
		if showid {
			table.SetHeader([]string{"STAT", "PATH", "ID"})
		} else {
			table.SetHeader([]string{"STAT", "PATH"})
		}
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

	// for i=0; i<=len(mr); i++ {
	for _, v := range *gl {
		if showid {
			row = []string{
				v.StatBlock,
				v.Path,
				fmt.Sprintf("%d", v.ID),
			}
		} else {
			row = []string{
				v.StatBlock,
				v.Path,
			}
		}
		table.Append(row)
	}

	table.Render()

	return buf.String()

}
