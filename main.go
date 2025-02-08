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
package main

import (
	"os"

	"github.com/mlofjard/contrack/cmd"
)

func main() {
	cmd.Execute()
	os.Exit(0)

	// `COLUMNSPEC:
	//
	//	A comma separated line of column names
	//
	//	container            The container name
	//	status               Short processing status (OK/ERR)
	//	detail               Long processing status error explaination
	//	repository           Repository (<domain>/<path>)
	//	image                Image (<domain>/<path>:<tag>)
	//	domain               Image domain
	//	path                 Image path
	//	tag                  Image tag
	//	update               Newer tag found
	//
	// `
}
