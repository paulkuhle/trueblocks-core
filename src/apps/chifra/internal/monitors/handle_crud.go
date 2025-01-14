// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package monitorsPkg

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/monitor"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"
)

// HandleCrudCommands
//
// [State]     | Delete | Undelete | Remove | Decache                    |
// ------------|--------|------------------------------------------------|
// Not Deleted | Delete	| Error    | Error  | Delete, Remove and Decache |
// Deleted     | Error  | Undelete | Remove | Remove and Decache         |
// ------------|--------|------------------------------------------------|
func (opts *MonitorsOptions) HandleCrudCommands() error {
	testMode := opts.Globals.TestMode
	for _, addr := range opts.Addrs {
		m := monitor.NewMonitor(opts.Globals.Chain, addr, false)
		if !file.FileExists(m.Path()) {
			return validate.Usage("No monitor was found for address " + addr + ".")

		} else {
			if opts.Decache {
				logger.Info("Decaching", addr)
				if opts.Globals.TestMode {
					logger.Info("Decaching monitor for address", addr, "not tested.")
					return nil
				}

				itemsSeen := int64(0)
				itemsRemoved := int64(0)
				bytesRemoved := int64(0)
				processorFunc := func(fileName string) bool {
					itemsSeen++
					if !file.FileExists(fileName) {
						logger.Progress(!testMode && itemsSeen%203 == 0, "Already removed ", fileName)
						return true // continue processing
					}

					itemsRemoved++
					bytesRemoved += file.FileSize(fileName)
					logger.Progress(!testMode && itemsRemoved%20 == 0, "Removed ", itemsRemoved, " items and ", bytesRemoved, " bytes.", fileName)

					os.Remove(fileName)
					if opts.Globals.Verbose {
						logger.Info(fileName, "was removed.")
					}
					path, _ := filepath.Split(fileName)
					if empty, _ := file.IsFolderEmpty(path); empty {
						os.RemoveAll(path)
						if opts.Globals.Verbose {
							logger.Info("Empty folder", path, "was removed.")
						}
					}

					return true
				}

				// Visits every item in the cache related to this monitor and calls into `processorFunc`
				err := m.Decache(opts.Globals.Chain, processorFunc)
				if err != nil {
					return err
				}
				logger.Info(itemsRemoved, "items totaling", bytesRemoved, "bytes were removed from the cache.", strings.Repeat(" ", 60))

				// We've visited them all, so delete the monitor itself
				m.Delete()
				logger.Info(("Monitor " + addr + " was deleted but not removed."))
				wasRemoved, err := m.Remove()
				if !wasRemoved || err != nil {
					logger.Info(("Monitor for " + addr + " was not removed (" + err.Error() + ")"))
				} else {
					logger.Info(("Monitor for " + addr + " was permanently removed."))
				}

			} else if opts.Undelete && !m.IsDeleted() {
				return validate.Usage("Monitor for {0} must be deleted before being undeleted.", addr)

			} else {
				if opts.Delete && opts.Remove {
					// do nothing, it will be resolved below...

				} else {
					if opts.Delete && m.IsDeleted() {
						return validate.Usage("Monitor for {0} is already deleted.", addr)

					} else if opts.Remove && !m.IsDeleted() {
						return validate.Usage("Cannot remove a file that has not previously been deleted.")
					}
				}
			}
		}
	}

	if opts.Decache {
		return nil
	}

	for _, addr := range opts.Addrs {
		m := monitor.NewMonitor(opts.Globals.Chain, addr, false)
		if opts.Undelete {
			m.UnDelete()
			logger.Info(("Monitor " + addr + " was undeleted."))

		} else {
			if opts.Delete {
				m.Delete()
				logger.Info(("Monitor " + addr + " was deleted but not removed."))
			}

			if opts.Remove {
				wasRemoved, err := m.Remove()
				if !wasRemoved || err != nil {
					logger.Info(("Monitor for " + addr + " was not removed (" + err.Error() + ")"))
				} else {
					logger.Info(("Monitor for " + addr + " was permanently removed."))
				}
			}
		}
	}

	return nil
}
