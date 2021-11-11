package export

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
/*
 * The file was auto generated with makeClass --gocmds. DO NOT EDIT.
 */

type ExportOptionsType struct {
	Appearances  bool
	Receipts     bool
	Statements   bool
	Logs         bool
	Traces       bool
	Accounting   bool
	Articulate   bool
	Cache        bool
	Cache_Traces bool
	Factory      bool
	Count        bool
	First_Record uint64
	Max_Records  uint64
	Relevant     bool
	Emitter      []string
	Topic        []string
	Clean        bool
	Freshen      bool
	Staging      bool
	Unripe       bool
	Load         string
	Reversed     bool
	By_Date      bool
	Summarize_By string
	Skip_Ddos    bool
	Max_Traces   uint64
	First_Block  uint64
	Last_Block   uint64
}

var Options ExportOptionsType
