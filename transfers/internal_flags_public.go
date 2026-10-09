//go:build !internal

package transfers

import "github.com/spf13/cobra"

func (t *Transfers) addInternalFlags(_ *cobra.Command) {}
