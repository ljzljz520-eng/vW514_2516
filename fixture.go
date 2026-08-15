package zoo

func FixedFixture() *Catalog {
	nodes := []Node{
		{ID: "entrance-north", Name: "北门", Kind: NodeEntrance, Slope: 1, StrollerFriendly: true},
		{ID: "entrance-south", Name: "南门", Kind: NodeEntrance, Slope: 2, StrollerFriendly: true},
		{ID: "tiger", Name: "虎山", Kind: NodeAnimal, Animal: "老虎", Slope: 5, StrollerFriendly: true},
		{ID: "panda", Name: "熊猫馆", Kind: NodeAnimal, Animal: "大熊猫", Slope: 7, StrollerFriendly: true},
		{ID: "penguin", Name: "企鹅馆", Kind: NodeAnimal, Animal: "企鹅", Slope: 2, StrollerFriendly: true},
		{ID: "savanna", Name: "非洲草原", Kind: NodeAnimal, Animal: "长颈鹿", Slope: 3, StrollerFriendly: false},
		{ID: "picnic-lawn", Name: "湖畔野餐点", Kind: NodeLunch, Slope: 1, StrollerFriendly: true},
		{ID: "family-cafe", Name: "家庭餐厅", Kind: NodeLunch, Slope: 0, StrollerFriendly: true},
		{ID: "exit-east", Name: "东门", Kind: NodeExit, Slope: 2, StrollerFriendly: true},
		{ID: "exit-west", Name: "西门", Kind: NodeExit, Slope: 2, StrollerFriendly: true},
	}
	connections := []Connection{
		{From: "entrance-north", To: "tiger", DistanceMeters: 420, Slope: 3, StrollerFriendly: true, Crowd: 2},
		{From: "entrance-north", To: "panda", DistanceMeters: 260, Slope: 8, StrollerFriendly: false, Crowd: 4},
		{From: "entrance-north", To: "penguin", DistanceMeters: 500, Slope: 2, StrollerFriendly: true, Crowd: 3},
		{From: "entrance-south", To: "tiger", DistanceMeters: 360, Slope: 4, StrollerFriendly: true, Crowd: 5},
		{From: "entrance-south", To: "savanna", DistanceMeters: 300, Slope: 2, StrollerFriendly: false, Crowd: 3},
		{From: "tiger", To: "panda", DistanceMeters: 180, Slope: 4, StrollerFriendly: true, Crowd: 5},
		{From: "tiger", To: "savanna", DistanceMeters: 350, Slope: 2, StrollerFriendly: true, Crowd: 3},
		{From: "panda", To: "penguin", DistanceMeters: 240, Slope: 1, StrollerFriendly: true, Crowd: 4},
		{From: "panda", To: "picnic-lawn", DistanceMeters: 300, Slope: 6, StrollerFriendly: false, Crowd: 2},
		{From: "penguin", To: "picnic-lawn", DistanceMeters: 220, Slope: 3, StrollerFriendly: true, Crowd: 2},
		{From: "savanna", To: "picnic-lawn", DistanceMeters: 260, Slope: 1, StrollerFriendly: true, Crowd: 6},
		{From: "picnic-lawn", To: "family-cafe", DistanceMeters: 100, Slope: 0, StrollerFriendly: true, Crowd: 5},
		{From: "picnic-lawn", To: "exit-east", DistanceMeters: 280, Slope: 4, StrollerFriendly: true, Crowd: 3},
		{From: "family-cafe", To: "exit-east", DistanceMeters: 220, Slope: 2, StrollerFriendly: true, Crowd: 2},
		{From: "savanna", To: "exit-west", DistanceMeters: 300, Slope: 2, StrollerFriendly: true, Crowd: 3},
		{From: "entrance-south", To: "exit-west", DistanceMeters: 620, Slope: 5, StrollerFriendly: true, Crowd: 4},
		{From: "exit-west", To: "picnic-lawn", DistanceMeters: 480, Slope: 3, StrollerFriendly: true, Crowd: 4},
	}
	shows := []Show{
		{NodeID: "tiger", StartMinute: 630, Duration: 15},
		{NodeID: "panda", StartMinute: 690, Duration: 20},
		{NodeID: "penguin", StartMinute: 660, Duration: 15},
	}
	return NewCatalog(nodes, connections, shows)
}
