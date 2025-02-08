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
package command

import (
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var Version = "dev-build"

type multiValueFlags []string

func (i multiValueFlags) Has(s string) bool {
	all := slices.Contains(i, "all")
	if all {
		return true
	}
	return slices.Contains(i, s)
}

func SetupCommandline(flagSet *pflag.FlagSet) multiValueFlags {
	var mockFlags multiValueFlags

	mockFlags, err := flagSet.GetStringSlice("mock")
	cobra.CheckErr(err)

	return mockFlags
}
