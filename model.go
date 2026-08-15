package zoo

import (
	"context"
	"errors"
)

type NodeKind string

const (
	NodeEntrance NodeKind = "entrance"
	NodeAnimal   NodeKind = "animal"
	NodeLunch    NodeKind = "lunch"
	NodeExit     NodeKind = "exit"
	NodeRest     NodeKind = "rest"
)

type Node struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Kind             NodeKind `json:"kind"`
	Animal           string   `json:"animal,omitempty"`
	Slope            int      `json:"slope"`
	StrollerFriendly bool     `json:"stroller_friendly"`
}

type Connection struct {
	From             string `json:"from"`
	To               string `json:"to"`
	DistanceMeters   int    `json:"distance_meters"`
	Slope            int    `json:"slope"`
	StrollerFriendly bool   `json:"stroller_friendly"`
	Crowd            int    `json:"crowd"`
}

type Show struct {
	NodeID      string `json:"node_id"`
	StartMinute int    `json:"start_minute"`
	Duration    int    `json:"duration_minutes"`
}

type Catalog struct {
	nodes       map[string]Node
	connections []Connection
	shows       map[string][]Show
}

type RouteRequest struct {
	EntranceID    string   `json:"entrance_id"`
	AnimalAreaIDs []string `json:"animal_area_ids"`
	LunchID       string   `json:"lunch_id"`
	ExitID        string   `json:"exit_id"`
	StartMinute   int      `json:"start_minute"`
	Stroller      bool     `json:"stroller"`
}

type RouteStop struct {
	Node            Node   `json:"node"`
	ArrivalMinute   int    `json:"arrival_minute"`
	DepartureMinute int    `json:"departure_minute"`
	Purpose         string `json:"purpose"`
	ShowWaitMinutes int    `json:"show_wait_minutes,omitempty"`
}

type RouteLeg struct {
	From           string       `json:"from"`
	To             string       `json:"to"`
	Path           []string     `json:"path"`
	Connections    []Connection `json:"connections"`
	DistanceMeters int          `json:"distance_meters"`
	TravelMinutes  int          `json:"travel_minutes"`
	Score          int          `json:"score"`
}

type RoutePlan struct {
	Stops                []RouteStop `json:"stops"`
	Legs                 []RouteLeg  `json:"legs"`
	TotalDistanceMeters  int         `json:"total_distance_meters"`
	WalkingMinutes       int         `json:"walking_minutes"`
	WaitingMinutes       int         `json:"waiting_minutes"`
	VisitMinutes         int         `json:"visit_minutes"`
	TotalDurationMinutes int         `json:"total_duration_minutes"`
	Score                int         `json:"score"`
}

type RouteStore interface {
	Save(context.Context, RoutePlan) error
}

var (
	ErrInvalidRequest = errors.New("invalid route request")
	ErrUnknownNode    = errors.New("unknown zoo node")
	ErrNoRoute        = errors.New("no route between requested nodes")
	ErrOutOfHours     = errors.New("route exceeds half-day limit")
)
