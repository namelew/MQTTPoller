package messages

type Status struct{
	Type string `json:"type"`
	Status string `json:"status"`
	Attr Command `json:"attr"`
}

type Command struct{
	Name string `json:"name"`
	Type string `json:"type"`
	Arguments map[string]interface{} `json:"arguments"`
}