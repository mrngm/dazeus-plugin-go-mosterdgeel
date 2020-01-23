package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
	"time"

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

	_, err = dz.SubscribeCommand("mosterdgeel", dazeus.NewUniversalScope(), HandleRecept)
	_, err = dz.SubscribeCommand("mosterd", dazeus.NewUniversalScope(), HandleRecept)
	_, err = dz.SubscribeCommand("recept", dazeus.NewUniversalScope(), HandleRecept)
	_, err = dz.SubscribeCommand("yolorecept", dazeus.NewUniversalScope(), HandleRandomRecept)
	if err != nil {
		panic(err)
	}

	listenerr := dz.Listen()
	panic(listenerr)
}

func HandleRecept(ev dazeus.Event) {
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

	if len(recipe.Content.String()) > 200 {
		ev.Reply(fmt.Sprintf("Hier is je recept: %s", recipe.Link), true)
		return
	}

	ev.Reply(recipe.Content.String(), true)
}

func HandleRandomRecept(ev dazeus.Event) {
	rand.Seed(time.Now().UnixNano())

	res, err := GetAllRecipes(ev)
	if err != nil {
		ev.Reply(fmt.Sprintf("E_MOSTERD: %v", err), true)
		return
	}

	choice := rand.Intn(len(res))

	ev.Reply(fmt.Sprintf("Hier is je recept: %s", res[choice].Link), true)
}
