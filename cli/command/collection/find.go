package collection

import (
	"fmt"

	"github.com/juliengk/go-utils"
	"github.com/kassisol/hbm/cli/command"
	collectionobj "github.com/kassisol/hbm/object/collection"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newFindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find [name]",
		Short: "Verify if collection exists in the whitelist",
		Long:  findDescription,
		Args:  cobra.ExactArgs(1),
		Run:   runFind,
	}

	return cmd
}

func runFind(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	c, err := collectionobj.New("sqlite", command.AppPath)
	if err != nil {
		log.Fatal(err)
	}
	defer c.End()

	fmt.Println(c.Find(args[0]))
}

var findDescription = `
Verify if collection exists in the whitelist

`
