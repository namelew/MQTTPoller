package messages

type Status struct{
	Type string `json:"type"`
	Status string `json:"status"`
	Attr Command `json:"attr"`
}

type Session struct{
	Id int
	Finish bool
	Status Status
}

type Command struct{
	Name string `json:"name"`
	Type string `json:"type"`
	Arguments map[string]interface{} `json:"arguments"`
}