package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"github.com/bonsai/aw.tui/graph"
	"github.com/bonsai/aw.tui/repos"
)

func main() {
	owner := flag.String("owner", "bonsai", "GitHub owner")
	out := flag.String("out", "graph.json", "graph output")
	flag.Parse()

	rs, err := repos.FetchRecent(*owner, 30)
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	g := graph.Build(rs)
	b, err := json.MarshalIndent(g, "", "  "); if err != nil { panic(err) }
	if err := os.WriteFile(*out, append(b, '\n'), 0644); err != nil { panic(err) }
	fmt.Printf("aw.tui: %d repos -> %d nodes, %d edges -> %s\n", len(rs), len(g.Nodes), len(g.Edges), *out)
}
