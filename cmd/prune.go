/*
Copyright © 2025 Mikael Lofjärd <mikael@lofjard.se>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mlofjard/contrack/utils"
)

//go:embed prune.md
var pruneMarkdown string

type longDurationValue time.Duration

func newLongDurationValue(val time.Duration, p *time.Duration) *longDurationValue {
	*p = val
	return (*longDurationValue)(p)
}

func (i *longDurationValue) Set(s string) error {
	var years = int64(0)
	var months = int64(0)
	var days = int64(0)
	var hours = int64(0)
	var err error

	yearSplit := strings.SplitN(s, "y", 2)
	if len(yearSplit) == 2 {
		years, err = strconv.ParseInt(yearSplit[0], 0, 64)
		utils.HandleErr(err)
	}
	monthSplit := strings.SplitN(yearSplit[len(yearSplit)-1], "m", 2)
	if len(monthSplit) == 2 {
		months, err = strconv.ParseInt(monthSplit[0], 0, 64)
		utils.HandleErr(err)
	}
	daySplit := strings.SplitN(monthSplit[len(monthSplit)-1], "d", 2)
	if len(daySplit) == 2 {
		days, err = strconv.ParseInt(daySplit[0], 0, 64)
		utils.HandleErr(err)
	}

	hours += (years * 365 * 24)
	hours += (months * 30 * 24)
	hours += (days * 24)

	v, err := time.ParseDuration(fmt.Sprintf("%dh", hours))
	*i = longDurationValue(v)
	return err
}

func (i *longDurationValue) Type() string {
	return "duration"
}

func (d *longDurationValue) String() string { return (*time.Duration)(d).String() }

var longDuration longDurationValue

func NewPruneCommand(rootViper *viper.Viper, renderer *glamour.TermRenderer) *cobra.Command {
	help := utils.ParseHelp(renderer, pruneMarkdown)

	// pruneCmd represents the prune command
	var pruneCmd = &cobra.Command{
		Use:     "prune <manifest file>",
		Short:   "A brief description of your command",
		Long:    help.Long,
		Example: help.Example,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("prune called")
			fmt.Println("duration", longDuration)
			fmt.Println("cutoff", time.Now().Add(-time.Duration(longDuration)))
		},
	}

	// Setup flags

	pruneCmd.Flags().IntP("keep", "k", 10, "Number of digests to keep")
	pruneCmd.Flags().VarP(&longDuration, "older-than", "o", "Print help (this message) and exit")
	pruneCmd.Flags().Bool("dry-run", false, "Print out digests to remove without removing them")
	pruneCmd.Flags().Bool("help", false, "Print help (this message) and exit")

	return pruneCmd
}
