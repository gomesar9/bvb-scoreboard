package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gomesar9/bvb-core/public/notify"
)

type Scoreboard struct {
	Title                string
	MatchNo              string
	Round, Phase, Gender string
	LocCity, LocCountry  string
	PlayerA1, PlayerA2   string
	PlayerB1, PlayerB2   string
	ScoreA, ScoreB       int
	FlagAURL, FlagBURL   string
	CountryA, CountryB   string
	SetScores            []string
}

func MatchToScoreboard(m notify.MatchInfo) *Scoreboard {
	var s Scoreboard
	const (
		teamA   int = 0
		teamB   int = 1
		player1 int = 0
		player2 int = 1
	)

	var players [2][]string
	players[teamA] = teamToPlayerNames(m.Teams[teamA].Name)
	players[teamB] = teamToPlayerNames(m.Teams[teamB].Name)
	round, _ := translateRound(m.RoundName)
	phase, _ := translatePhase(m.Phase)
	gender := translateGender(m.Gender)

	s = Scoreboard{
		MatchNo:    strconv.FormatUint(uint64(m.MatchNo), 10),
		Title:      genTitle(m, ""), // Main Draw - Pool F - Central court - Men #24
		Round:      round,
		Phase:      phase,
		Gender:     gender,
		LocCity:    m.LocCity,
		LocCountry: m.LocCountry,
		PlayerA1:   players[teamA][player1],
		PlayerA2:   players[teamA][player2],
		PlayerB1:   players[teamB][player1],
		PlayerB2:   players[teamB][player2],
		ScoreA:     m.Teams[teamA].Score,
		ScoreB:     m.Teams[teamB].Score,
		FlagAURL:   imgUrlOfCountry(m.Teams[teamA].Country), // Bandeira do time A
		FlagBURL:   imgUrlOfCountry(m.Teams[teamB].Country), // Bandeira do time B
		CountryA:   m.Teams[teamA].Country.Name,
		CountryB:   m.Teams[teamB].Country.Name,
		SetScores:  formatSets(m.Sets), // Parciais dos sets
	}

	return &s
}

func genTitle(m notify.MatchInfo, kind string) string {
	var title string
	phase, err := translatePhase(m.Phase)
	if err != nil {
		// TODO: log
		phase = m.Phase
	}

	round, err := translateRound(m.RoundName)
	if err != nil {
		// TODO: log
		round = m.RoundName
	}

	gender := translateGender(m.Gender)

	switch kind {
	case "bvbsite":
		// Match.Phase.Name - Match.RoundName - Match.CourtText - Match.Gender #Match.MatchNoInTournament
		itens := []string{phase, round, m.CourtName, gender}
		title = fmt.Sprintf("%s #%d", strings.Join(itens, " - "), m.MatchNoInTournament)

	default:
		title = fmt.Sprintf("%s, %s (%s)", phase, round, gender)
	}

	return title
}

func teamToPlayerNames(team string) []string {
	return strings.Split(team, "/")
}

func imgUrlOfCountry(country notify.CountryInfo) string {
	return fmt.Sprintf("https://flagcdn.com/80x60/%s.png", strings.ToLower(country.Alpha2))
}

func formatSets(setsInfo [][]int) []string {
	var s []string
	for _, s := range setsInfo {
		strings.Join([]string{strconv.Itoa(s[0]), strconv.Itoa(s[1])}, "-")
	}
	return s
}

func translatePhase(p string) (string, error) {
	switch p {
	case "Qualification":
		return "Qualificatórias", nil
	case "Main Draw":
		return "Principal", nil
	default:
		return "", fmt.Errorf("Invalid Phase %s", p)

	}
}

func translateRound(r string) (string, error) {
	switch r {
	case "Preliminary Phase":
		return "Fase preliminar", nil
	case "Second round":
		return "Segunda rodada", nil
	case "Pool A":
		return "Grupo A", nil
	case "Pool B":
		return "Grupo B", nil
	case "Pool C":
		return "Grupo C", nil
	case "Pool D":
		return "Grupo D", nil
	case "Pool E":
		return "Grupo E", nil
	case "Pool F":
		return "Grupo F", nil
	case "Pool G":
		return "Grupo G", nil
	case "Pool H":
		return "Grupo H", nil
	case "Round of 24": // Para torneios com 8 grupos
		return "Rodada de 24", nil
	case "Round of 16": // Para torneio com 8 grupos
		return "Oitavas", nil
	case "Round of 18": // Para torneios com 6 grupos
		return "Rodada de 18", nil
	case "Round of 12": // Para torneios com 6 grupos
		return "Oitavas modificadas", nil
	case "Quarter-finals":
		return "Quartas de final", nil
	case "Semi-finals":
		return "Semi-finais", nil
	case "3rd place match":
		return "Disputa de terceiro", nil
	case "Final":
		return "Final", nil
	default:
		return "", fmt.Errorf("Invalid round %s", r)
	}
}

func translateGender(g string) string {
	if g == "Men" {
		return "Masc"
	}
	return "Fem"
}
