// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were generated with makeClass --run. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package initPkg

// EXISTING_CODE
import (
	"net/http"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/internal/globals"
	outputHelpers "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output/helpers"
	"github.com/spf13/cobra"
)

// EXISTING_CODE

// RunInit handles the init command for the command line. Returns error only as per cobra.
func RunInit(cmd *cobra.Command, args []string) (err error) {
	opts := initFinishParse(args)
	outputHelpers.SetEnabledForCmds("init", opts.IsPorted())
	outputHelpers.SetWriterForCommand("init", &opts.Globals)
	// EXISTING_CODE
	// EXISTING_CODE
	err, _ = opts.InitInternal()
	return
}

// ServeInit handles the init command for the API. Returns error and a bool if handled
func ServeInit(w http.ResponseWriter, r *http.Request) (err error, handled bool) {
	opts := initFinishParseApi(w, r)
	outputHelpers.SetEnabledForCmds("init", opts.IsPorted())
	outputHelpers.InitJsonWriterApi("init", w, &opts.Globals)
	// EXISTING_CODE
	// EXISTING_CODE
	err, handled = opts.InitInternal()
	outputHelpers.CloseJsonWriterIfNeededApi("init", err, &opts.Globals)
	return
}

// InitInternal handles the internal workings of the init command.  Returns error and a bool if handled
func (opts *InitOptions) InitInternal() (err error, handled bool) {
	err = opts.validateInit()
	if err != nil {
		return err, true
	}

	// EXISTING_CODE
	handled = true

	err = opts.HandleInit()
	// EXISTING_CODE

	return
}

// GetInitOptions returns the options for this tool so other tools may use it.
func GetInitOptions(args []string, g *globals.GlobalOptions) *InitOptions {
	ret := initFinishParse(args)
	if g != nil {
		ret.Globals = *g
	}
	return ret
}

func (opts *InitOptions) IsPorted() (ported bool) {
	// EXISTING_CODE
	// EXISTING_CODE
	return
}

// EXISTING_CODE
// EXISTING_CODE
