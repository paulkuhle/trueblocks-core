package names

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
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/cmd/root"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"
	"github.com/spf13/cobra"
)

func Validate(cmd *cobra.Command, args []string) error {
	if Options.Tags && anyBase() {
		return validate.Usage("Do not use the --tags option with any other option.")
	}

	if Options.Collections && anyBase() {
		return validate.Usage("Do not use the --collection option with any other option.")
	}

	err := root.ValidateGlobals(cmd, args)
	if err != nil {
		return err
	}

	return nil
}

func anyBase() bool {
	return Options.Expand ||
		Options.Match_Case ||
		Options.All ||
		Options.Custom ||
		Options.Prefund ||
		Options.Named ||
		Options.Addr ||
		Options.To_Custom ||
		Options.Clean
}