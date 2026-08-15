package zoo

import (
	"fmt"
	"sort"
)

func NewCatalog(nodes []Node, connections []Connection, shows []Show) *Catalog {
	c := &Catalog{
		nodes:       make(map[string]Node, len(nodes)),
		connections: append([]Connection(nil), connections...),
		shows:       make(map[string][]Show),
	}
	for _, node := range nodes {
		c.nodes[node.ID] = node
	}
	for _, show := range shows {
		c.shows[show.NodeID] = append(c.shows[show.NodeID], show)
	}
	for id := range c.shows {
		sort.Slice(c.shows[id], func(i, j int) bool {
			return c.shows[id][i].StartMinute < c.shows[id][j].StartMinute
		})
	}
	return c
}

func (c *Catalog) Node(id string) (Node, bool) {
	node, ok := c.nodes[id]
	return node, ok
}

func (c *Catalog) Nodes() []Node {
	ids := make([]string, 0, len(c.nodes))
	for id := range c.nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	result := make([]Node, 0, len(ids))
	for _, id := range ids {
		result = append(result, c.nodes[id])
	}
	return result
}

func (c *Catalog) Connections() []Connection {
	return append([]Connection(nil), c.connections...)
}

func (c *Catalog) Shows(nodeID string) []Show {
	return append([]Show(nil), c.shows[nodeID]...)
}

func (c *Catalog) UpsertNode(node Node) error {
	if node.ID == "" || node.Name == "" || node.Kind == "" {
		return fmt.Errorf("%w: node fields are required", ErrInvalidRequest)
	}
	if node.Slope < 0 || node.Slope > 100 {
		return fmt.Errorf("%w: node slope must be 0..100", ErrInvalidRequest)
	}
	c.nodes[node.ID] = node
	return nil
}

func (c *Catalog) AddConnection(connection Connection) error {
	if _, ok := c.nodes[connection.From]; !ok {
		return fmt.Errorf("%w: %s", ErrUnknownNode, connection.From)
	}
	if _, ok := c.nodes[connection.To]; !ok {
		return fmt.Errorf("%w: %s", ErrUnknownNode, connection.To)
	}
	if connection.From == connection.To || connection.DistanceMeters <= 0 || connection.Crowd < 0 || connection.Crowd > 100 {
		return fmt.Errorf("%w: invalid connection", ErrInvalidRequest)
	}
	for i, existing := range c.connections {
		if existing.From == connection.From && existing.To == connection.To {
			c.connections[i] = connection
			return nil
		}
		if existing.From == connection.To && existing.To == connection.From {
			c.connections[i] = connection
			return nil
		}
	}
	c.connections = append(c.connections, connection)
	return nil
}

func (c *Catalog) SetCrowd(from, to string, crowd int) error {
	if crowd < 0 || crowd > 100 {
		return fmt.Errorf("%w: crowd must be 0..100", ErrInvalidRequest)
	}
	for i := range c.connections {
		if (c.connections[i].From == from && c.connections[i].To == to) || (c.connections[i].From == to && c.connections[i].To == from) {
			c.connections[i].Crowd = crowd
			return nil
		}
	}
	return fmt.Errorf("%w: connection %s-%s", ErrNoRoute, from, to)
}
