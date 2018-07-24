package main

import (
	"testing"
)

func TestShovelLockCheck(t *testing.T) {
	r, _ := initRoute()
	g := r.Graph

	// make sure that needing a shovel in advance passes
	// this route is via the swamp portal
	g.Map["shovel"].AddParents(g.Map["d0 sword chest"])
	g.Map["bracelet"].AddParents(g.Map["maku key fall"])
	g.Map["flippers"].AddParents(g.Map["blaino gift"])
	g.Map["feather L-1"].AddParents(g.Map["star ore spot"])
	if canShovelSoftlock(g) {
		t.Error("false positive shovel softlock w/ shovel prereq")
	}
	// and make sure the shovel's parents are unchanged
	if len(g.Map["shovel"].Parents()) != 1 {
		t.Fatal("shovel parents altered by safety check")
	}

	// make sure that getting there with no shovel fails
	g.Map["shovel"].ClearParents()
	g.Map["bracelet"].ClearParents()
	g.Map["bracelet"].AddParents(g.Map["d0 sword chest"])
	g.Map["feather L-1"].ClearParents()
	g.Map["feather L-1"].AddParents(g.Map["maku key fall"])
	if !canShovelSoftlock(g) {
		t.Error("false negative shovel softlock w/ no shovel")
	}

	// make sure that getting a shovel as the gift passes
	g.Map["shovel"].ClearParents()
	g.Map["shovel"].AddParents(g.Map["shovel gift"])
	if canShovelSoftlock(g) {
		t.Error("false positive shovel softlock w/ shovel as gift")
	}

	// and make sure that getting there with an optional shovel fails
	g.Map["shovel"].ClearParents()
	g.Map["shovel"].AddParents(g.Map["boomerang gift"])
	if !canShovelSoftlock(g) {
		t.Error("false negative shovel softlock w/ optional shovel")
	}
}

func TestPortalLockCheck(t *testing.T) {
	r, _ := initRoute()
	g := r.Graph

	// make sure that an obviously safe graph passes
	g.Map["sword L-1"].AddParents(g.Map["d0 sword chest"])
	g.Map["satchel"].AddParents(g.Map["maku key fall"])
	if canRosaPortalSoftlock(g) {
		t.Error("false positive rosa portal softlock")
	}
	// and make sure the nodes' parents are unchanged
	if len(g.Map["sword L-1"].Parents()) != 1 {
		t.Fatal("sword L-1 parents altered by safety check")
	}
	if len(g.Map["rosa portal in wrapper"].Parents()) == 0 {
		t.Fatal("rosa portal in wrapper parents altered by safety check")
	}

	// make sure than an obviously unsafe graph fails.
	// this vulnerability is via the lake portal.
	g.Map["sword L-1"].ClearParents()
	g.Map["fool's ore"].AddParents(g.Map["d0 sword chest"])
	g.Map["flippers"].AddParents(g.Map["d0 rupee chest"]) // contrived, yes
	g.Map["feather L-2"].AddParents(g.Map["blaino gift"])
	if !canRosaPortalSoftlock(g) {
		t.Error("false negative rosa portal softlock")
	}

	// make sure that check works even if there's an optional sword
	g.Map["sword L-1"].AddParents(g.Map["floodgate gift"])
	if !canRosaPortalSoftlock(g) {
		t.Error("false negative rosa portal softlock")
	}
}
