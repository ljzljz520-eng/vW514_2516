package zoo

import (
	"fmt"
	"sort"
)

type Planner struct {
	catalog *Catalog
}

func NewPlanner(catalog *Catalog) *Planner {
	return &Planner{catalog: catalog}
}

func (p *Planner) Validate(req RouteRequest) error {
	if p == nil || p.catalog == nil {
		return fmt.Errorf("%w: catalog is required", ErrInvalidRequest)
	}
	if req.StartMinute < 0 || req.StartMinute >= 1440 || req.EntranceID == "" || req.LunchID == "" || req.ExitID == "" || len(req.AnimalAreaIDs) == 0 {
		return fmt.Errorf("%w: entrance, animals, lunch, exit and start time are required", ErrInvalidRequest)
	}
	entrance, ok := p.catalog.Node(req.EntranceID)
	if !ok || entrance.Kind != NodeEntrance {
		return fmt.Errorf("%w: entrance %s", ErrUnknownNode, req.EntranceID)
	}
	lunch, ok := p.catalog.Node(req.LunchID)
	if !ok || lunch.Kind != NodeLunch {
		return fmt.Errorf("%w: lunch %s", ErrUnknownNode, req.LunchID)
	}
	exit, ok := p.catalog.Node(req.ExitID)
	if !ok || exit.Kind != NodeExit {
		return fmt.Errorf("%w: exit %s", ErrUnknownNode, req.ExitID)
	}
	seen := make(map[string]struct{}, len(req.AnimalAreaIDs))
	for _, id := range req.AnimalAreaIDs {
		if _, duplicate := seen[id]; duplicate {
			return fmt.Errorf("%w: duplicate animal %s", ErrInvalidRequest, id)
		}
		seen[id] = struct{}{}
		node, exists := p.catalog.Node(id)
		if !exists || node.Kind != NodeAnimal {
			return fmt.Errorf("%w: animal area %s", ErrUnknownNode, id)
		}
	}
	return nil
}

type pathResult struct {
	nodes    []string
	edges    []Connection
	distance int
	minutes  int
	score    int
}

func (p *Planner) Plan(req RouteRequest) (RoutePlan, error) {
	if err := p.Validate(req); err != nil {
		return RoutePlan{}, err
	}
	best := RoutePlan{}
	var bestKey string
	order := append([]string(nil), req.AnimalAreaIDs...)
	first := true
	var visit func(int)
	visit = func(index int) {
		if index == len(order) {
			candidate, err := p.planOrder(req, order)
			if err != nil {
				return
			}
			key := orderKey(order)
			if first || candidate.Score < best.Score || (candidate.Score == best.Score && key < bestKey) {
				best, bestKey, first = candidate, key, false
			}
			return
		}
		for i := index; i < len(order); i++ {
			order[index], order[i] = order[i], order[index]
			visit(index + 1)
			order[index], order[i] = order[i], order[index]
		}
	}
	visit(0)
	if first {
		return RoutePlan{}, ErrNoRoute
	}
	return best, nil
}

func orderKey(order []string) string {
	key := ""
	for _, id := range order {
		key += id + "|"
	}
	return key
}

func (p *Planner) planOrder(req RouteRequest, order []string) (RoutePlan, error) {
	plan := RoutePlan{}
	minute := req.StartMinute
	current := req.EntranceID
	plan.Stops = append(plan.Stops, RouteStop{Node: p.catalog.nodes[current], ArrivalMinute: minute, DepartureMinute: minute, Purpose: "entrance"})
	for _, next := range order {
		leg, err := p.shortestPath(current, next, req.Stroller)
		if err != nil {
			return RoutePlan{}, err
		}
		plan.Legs = append(plan.Legs, leg)
		plan.TotalDistanceMeters += leg.DistanceMeters
		plan.WalkingMinutes += leg.TravelMinutes
		plan.Score += leg.Score
		minute += leg.TravelMinutes
		wait := p.showWait(next, minute)
		minute += wait
		departure := minute + 25
		plan.WaitingMinutes += wait
		plan.VisitMinutes += 25
		plan.Score += wait * 3
		plan.Stops = append(plan.Stops, RouteStop{Node: p.catalog.nodes[next], ArrivalMinute: minute - wait, DepartureMinute: departure, Purpose: "animal", ShowWaitMinutes: wait})
		minute = departure
		current = next
	}
	leg, err := p.shortestPath(current, req.LunchID, req.Stroller)
	if err != nil {
		return RoutePlan{}, err
	}
	plan.Legs = append(plan.Legs, leg)
	plan.TotalDistanceMeters += leg.DistanceMeters
	plan.WalkingMinutes += leg.TravelMinutes
	plan.Score += leg.Score
	minute += leg.TravelMinutes
	plan.Stops = append(plan.Stops, RouteStop{Node: p.catalog.nodes[req.LunchID], ArrivalMinute: minute, DepartureMinute: minute + 40, Purpose: "lunch"})
	plan.VisitMinutes += 40
	minute += 40
	current = req.LunchID
	leg, err = p.shortestPath(current, req.ExitID, req.Stroller)
	if err != nil {
		return RoutePlan{}, err
	}
	plan.Legs = append(plan.Legs, leg)
	plan.TotalDistanceMeters += leg.DistanceMeters
	plan.WalkingMinutes += leg.TravelMinutes
	plan.Score += leg.Score
	minute += leg.TravelMinutes
	plan.Stops = append(plan.Stops, RouteStop{Node: p.catalog.nodes[req.ExitID], ArrivalMinute: minute, DepartureMinute: minute, Purpose: "exit"})
	plan.TotalDurationMinutes = minute - req.StartMinute
	if plan.TotalDurationMinutes > 360 {
		return RoutePlan{}, ErrOutOfHours
	}
	return plan, nil
}

func (p *Planner) showWait(nodeID string, arrival int) int {
	for _, show := range p.catalog.shows[nodeID] {
		if show.StartMinute >= arrival && show.StartMinute-arrival <= 30 {
			return show.StartMinute - arrival
		}
	}
	return 0
}

func edgeMinutes(edge Connection) int {
	minutes := (edge.DistanceMeters + 49) / 50
	minutes += edge.Slope / 5
	minutes += edge.Crowd / 25
	if minutes < 1 {
		return 1
	}
	return minutes
}

func edgeScore(edge Connection, stroller bool) int {
	score := edge.DistanceMeters + edge.Slope*35 + edge.Crowd*20
	if stroller && !edge.StrollerFriendly {
		score += 5000
	}
	return score
}

func (p *Planner) shortestPath(from, to string, stroller bool) (RouteLeg, error) {
	if _, ok := p.catalog.nodes[from]; !ok {
		return RouteLeg{}, ErrUnknownNode
	}
	if _, ok := p.catalog.nodes[to]; !ok {
		return RouteLeg{}, ErrUnknownNode
	}
	if from == to {
		return RouteLeg{From: from, To: to, Path: []string{from}}, nil
	}
	adj := make(map[string][]Connection)
	for _, edge := range p.catalog.connections {
		adj[edge.From] = append(adj[edge.From], edge)
		adj[edge.To] = append(adj[edge.To], Connection{From: edge.To, To: edge.From, DistanceMeters: edge.DistanceMeters, Slope: edge.Slope, StrollerFriendly: edge.StrollerFriendly, Crowd: edge.Crowd})
	}
	for id := range adj {
		sort.Slice(adj[id], func(i, j int) bool {
			if adj[id][i].To == adj[id][j].To {
				return edgeScore(adj[id][i], stroller) < edgeScore(adj[id][j], stroller)
			}
			return adj[id][i].To < adj[id][j].To
		})
	}
	const inf = int(^uint(0) >> 1)
	dist := make(map[string]int, len(p.catalog.nodes))
	previous := make(map[string]string, len(p.catalog.nodes))
	previousEdge := make(map[string]Connection, len(p.catalog.nodes))
	visited := make(map[string]bool, len(p.catalog.nodes))
	for id := range p.catalog.nodes {
		dist[id] = inf
	}
	dist[from] = 0
	for range p.catalog.nodes {
		current := ""
		best := inf
		for id, value := range dist {
			if !visited[id] && (value < best || (value == best && (current == "" || id < current))) {
				current, best = id, value
			}
		}
		if current == "" || best == inf {
			break
		}
		visited[current] = true
		if current == to {
			break
		}
		for _, edge := range adj[current] {
			candidate := best + edgeScore(edge, stroller)
			if candidate < dist[edge.To] || (candidate == dist[edge.To] && current < previous[edge.To]) {
				dist[edge.To] = candidate
				previous[edge.To] = current
				previousEdge[edge.To] = edge
			}
		}
	}
	if dist[to] == inf {
		return RouteLeg{}, ErrNoRoute
	}
	nodes := []string{to}
	edges := make([]Connection, 0)
	for current := to; current != from; {
		parent, ok := previous[current]
		if !ok {
			return RouteLeg{}, ErrNoRoute
		}
		edges = append(edges, previousEdge[current])
		nodes = append(nodes, parent)
		current = parent
	}
	for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}
	for i, j := 0, len(edges)-1; i < j; i, j = i+1, j-1 {
		edges[i], edges[j] = edges[j], edges[i]
	}
	minutes, distance, score := 0, 0, 0
	for _, edge := range edges {
		minutes += edgeMinutes(edge)
		distance += edge.DistanceMeters
		score += edgeScore(edge, stroller)
	}
	return RouteLeg{From: from, To: to, Path: nodes, Connections: edges, DistanceMeters: distance, TravelMinutes: minutes, Score: score}, nil
}
