package list

/*-------------------------------------------------------------------------------------------
 * qblocks - fast, easily-accessible, fully-decentralized data from blockchains
 * copyright (c) 2016, 2021 TrueBlocks, LLC (http://trueblocks.io)
 *
 * This program is free software: you may redistribute it and/or modify it under the terms
 * of the GNU General Public License as published by the Free Software Foundation, either
 * version 3 of the License, or (at your option) any later version. This program is
 * distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even
 * the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details. You should have received a copy of the GNU General
 * Public License along with this program. If not, see http://www.gnu.org/licenses/.
 *-------------------------------------------------------------------------------------------*/

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/cmd/root"
	"github.com/spf13/cobra"
)

func Run(cmd *cobra.Command, args []string) {
	options := " --appearances"
	if Options.Count {
		options += " --count"
	}
	if Options.First_Block > 0 {
		options += " --first_block " + fmt.Sprintf("%d", Options.First_Block)
	}
	if Options.Last_Block > 0 {
		options += " --last_block " + fmt.Sprintf("%d", Options.Last_Block)
	}
	arguments := ""
	for _, arg := range args {
		arguments += " " + arg
	}
	root.PassItOn("acctExport", options, arguments)
}