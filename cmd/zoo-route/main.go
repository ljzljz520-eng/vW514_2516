package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"zoo-route"
)

func main() {
	entrance := flag.String("entrance", "entrance-north", "entrance node id")
	animals := flag.String("animals", "tiger,panda,penguin", "comma-separated animal area ids")
	lunch := flag.String("lunch", "picnic-lawn", "lunch node id")
	exit := flag.String("exit", "exit-east", "exit node id")
	start := flag.Int("start", 540, "start minute after midnight")
	stroller := flag.Bool("stroller", true, "prefer stroller-friendly paths")
	list := flag.Bool("list", false, "list fixture nodes")
	flag.Parse()

	catalog := zoo.FixedFixture()
	if *list {
		for _, node := range catalog.Nodes() {
			fmt.Printf("%s\t%s\t%s\n", node.ID, node.Kind, node.Name)
		}
		return
	}
	store := zoo.NewMemoryStore()
	service := zoo.NewService(zoo.NewPlanner(catalog), store)
	request := zoo.RouteRequest{
		EntranceID:    *entrance,
		AnimalAreaIDs: splitIDs(*animals),
		LunchID:       *lunch,
		ExitID:        *exit,
		StartMinute:   *start,
		Stroller:      *stroller,
	}
	plan, err := service.Execute(context.Background(), request)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	output, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}

func splitIDs(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	sort.Strings(result)
	return result
}
