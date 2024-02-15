package librematch

type Lobby struct {
	Id            uint      `json:"id"`
	SteamLobbyId  uint      `json:"steamlobbyid"`
	HostProfileId uint      `json:"host_profile_id"`
	Description   string    `json:"description"`
	MatchMembers  []*Member `json:"matchmembers"`
}
