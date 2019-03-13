package api_root

func fixtureRootResource(host string) *RootResource {
	return &RootResource{
		ID:		1,
		Title: "RootResource",
		Endpoints: []*Endpoint{
			{
				ID:	1,
				Path: host + "/",
			},
			{
				ID: 2,
				Path: host + "/tasks",
			},
			{
				ID: 3,
				Path: host + "/images",
			},
		},
	}
}
