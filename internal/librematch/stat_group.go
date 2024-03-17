package librematch

type StatGroup struct {
	Id      uint               `json:"id"`
	Members []*StatGroupMember `json:"members"`
}
