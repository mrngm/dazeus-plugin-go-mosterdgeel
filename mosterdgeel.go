package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/dazeus/dazeus-go"
)

var myCommand string

func main() {
	connStr := "unix:/tmp/dazeus.sock"
	if len(os.Args) > 1 {
		connStr = os.Args[1]
	}
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("Paniek! %v\n", p)
			debug.PrintStack()
		}
	}()

	dz, err := dazeus.ConnectWithLoggingToStdErr(connStr)
	if err != nil {
		panic(err)
	}

	if _, hlerr := dz.HighlightCharacter(); hlerr != nil {
		panic(hlerr)
	}

	_, err = dz.SubscribeCommand("mosterdgeel", dazeus.NewUniversalScope(), func(ev dazeus.Event) {
		res, err := GetPossibleRecipes(ev)
		if err != nil {
			ev.Reply(fmt.Sprintf("E_MOSTERD: %v", err), true)
			return
		}
		recipe, err := GetRecipe(res)
		if err != nil {
			ev.Reply(fmt.Sprintf("E_MOSTERD: %v", err), true)
			return
		}

		ev.Reply(recipe.Content.String(), true)
	})
	if err != nil {
		panic(err)
	}

	listenerr := dz.Listen()
	panic(listenerr)
}
