package zoo

import (
	"testing"
)

func TestPlanIncludesRequestedStopsAndExit(t *testing.T) {
	planner := NewPlanner(FixedFixture())
	plan, err := planner.Plan(RouteRequest{EntranceID: "entrance-north", AnimalAreaIDs: []string{"tiger", "panda"}, LunchID: "picnic-lawn", ExitID: "exit-east", StartMinute: 540, Stroller: true})
	if err != nil {
		t.Fatalf("plan failed: %v", err)
	}
	if len(plan.Stops) != 5 {
		t.Fatalf("got %d stops", len(plan.Stops))
	}
	if plan.Stops[0].Node.ID != "entrance-north" || plan.Stops[len(plan.Stops)-1].Node.ID != "exit-east" {
		t.Fatalf("unexpected endpoints: %s to %s", plan.Stops[0].Node.ID, plan.Stops[len(plan.Stops)-1].Node.ID)
	}
	seen := map[string]bool{}
	for _, stop := range plan.Stops {
		seen[stop.Node.ID] = true
	}
	for _, id := range []string{"tiger", "panda", "picnic-lawn"} {
		if !seen[id] {
			t.Fatalf("missing stop %s", id)
		}
	}
}

func TestStrollerPlanAvoidsUnfriendlyDirectConnection(t *testing.T) {
	planner := NewPlanner(FixedFixture())
	leg, err := planner.shortestPath("entrance-north", "panda", true)
	if err != nil {
		t.Fatalf("path failed: %v", err)
	}
	for _, edge := range leg.Connections {
		if !edge.StrollerFriendly {
			t.Fatalf("selected non stroller-friendly edge %+v", edge)
		}
	}
}

func TestAdminCanUpdateCrowd(t *testing.T) {
	catalog := FixedFixture()
	admin := NewAdminService(catalog)
	if err := admin.SetCrowd("tiger", "panda", 1); err != nil {
		t.Fatalf("set crowd failed: %v", err)
	}
	found := false
	for _, edge := range catalog.Connections() {
		if (edge.From == "tiger" && edge.To == "panda") || (edge.From == "panda" && edge.To == "tiger") {
			found = edge.Crowd == 1
		}
	}
	if !found {
		t.Fatal("crowd was not updated")
	}
}
