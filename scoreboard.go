package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gomesar9/bvb-core/public/notify"
)

type ScoreboardData struct {
	MatchNo             string
	Stage               string
	LocCity, LocCountry string
	PlayerA1, PlayerA2  string
	PlayerB1, PlayerB2  string
	ScoreA, ScoreB      int
	FlagA, FlagB        string
	CountryA, CountryB  string
	SetScores           []string
}

type MediaData struct {
	Scoreboards []ScoreboardData
}

// Estrutura para os parâmetros do HTML
type MatchData struct {
}

func genTitle(m notify.MatchInfo, kind string) string {
	var title string

	if kind == "bvbsite" {
		// Match.Phase.Name - Match.RoundName - Match.CourtText - Match.Gender #Match.MatchNoInTournament
		itens := []string{m.Phase, m.RoundName, m.CourtName, m.Gender}
		title = fmt.Sprintf("%s #%d", strings.Join(itens, " - "), m.MatchNoInTournament)
	}

	title = fmt.Sprintf("[%s] %s", m.Phase, m.Gender)
	return title
}

func teamToPlayerNames(team string) []string {
	return strings.Split(team, "/")
}

func formatSets(setsInfo [][]int) []string {
	var s []string
	for _, s := range setsInfo {
		strings.Join([]string{strconv.Itoa(s[0]), strconv.Itoa(s[1])}, "-")
	}
	return s
}

func test() {
	_ = ScoreboardData{
		Stage:      "Chave Principal - Grupo F - Masc", // Main Draw - Pool F - Central court - Men #24
		LocCity:    "Brasília",
		LocCountry: "Brazil",
		PlayerA1:   "Arthur",
		PlayerA2:   "Adrielson",
		PlayerB1:   "Pedrosa",
		PlayerB2:   "Campos",
		ScoreA:     2,
		ScoreB:     0,
		FlagA:      "https://flagcdn.com/w40/br.png", // Bandeira do time A
		FlagB:      "https://flagcdn.com/w40/pt.png", // Bandeira do time B
		CountryA:   "Brasil",
		CountryB:   "Portugal",
		SetScores:  []string{"21–16", "21–19"}, // Parciais dos sets
	}
}

func matchToScoreboard(m notify.MatchInfo) *ScoreboardData {
	var s ScoreboardData
	const (
		teamA int = 0
		teamB int = 1
	)

	var players [2][]string
	players[teamA] = teamToPlayerNames(m.Teams[teamA].Name)
	players[teamB] = teamToPlayerNames(m.Teams[teamB].Name)

	s = ScoreboardData{
		MatchNo:    strconv.FormatUint(uint64(m.MatchNo), 10),
		Stage:      genTitle(m, ""), // Main Draw - Pool F - Central court - Men #24
		LocCity:    m.LocCity,
		LocCountry: m.LocCountry,
		PlayerA1:   players[teamA][0],
		PlayerA2:   players[teamA][1],
		PlayerB1:   players[teamB][0],
		PlayerB2:   players[teamB][1],
		ScoreA:     m.Teams[teamA].Score,
		ScoreB:     m.Teams[teamB].Score,
		FlagA:      imgUrlOfCountry(m.Teams[teamA].Country), // Bandeira do time A
		FlagB:      imgUrlOfCountry(m.Teams[teamB].Country), // Bandeira do time B
		CountryA:   m.Teams[teamA].Country.Name,
		CountryB:   m.Teams[teamB].Country.Name,
		SetScores:  formatSets(m.Sets), // Parciais dos sets
	}

	return &s
}

func imgUrlOfCountry(country notify.CountryInfo) string {
	return fmt.Sprintf("https://flagcdn.com/80x60/%s.png", strings.ToLower(country.Alpha2))
}

func makeInstagramStoryMedia(matches []notify.MatchInfo) {
	var mediaData MediaData
	var matchNoList []string

	// TODO: Tratar em batches
	if len(matches) > 4 {
		matches = matches[:4]
	}

	for _, m := range matches {
		data := matchToScoreboard(m)
		mediaData.Scoreboards = append(mediaData.Scoreboards, *data)
		matchNoList = append(matchNoList, strconv.FormatUint(uint64(m.MatchNo), 10))
	}

	// Renderiza o HTML com dados
	tpl, err := template.ParseFiles("templates/scoreboard.html")
	if err != nil {
		log.Fatalf("Erro ao carregar HTML: %v", err)
	}

	var htmlBuffer bytes.Buffer
	err = tpl.Execute(&htmlBuffer, mediaData)
	if err != nil {
		log.Fatalf("Erro ao renderizar template: %v", err)
	}

	// Salva o HTML em um arquivo temporário
	htmlPath := fmt.Sprintf("test/match-%s-rendered.html", strings.Join(matchNoList, "_"))
	err = os.WriteFile(htmlPath, htmlBuffer.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Erro ao salvar HTML: %v", err)
	}
	fmt.Printf("HTML salvo em: %v\n", htmlPath)
	time.Sleep(50 * time.Millisecond) // Aguarda arquivo ser completamente salvo ???

	// Tira screenshot com chromedp
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	absPath, err := filepath.Abs(htmlPath)
	if err != nil {
		log.Fatalf("Erro ao obter caminho absoluto: %v", err)
	}

	fullPath := "file://" + absPath

	// Executa com timeout para segurança
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Executa as ações do chromedp
	log.Println("Gerando imagem a partir do template")
	err = chromedp.Run(ctx,
		chromedp.Navigate(fullPath),
		chromedp.EmulateViewport(1080, 1920),
		chromedp.WaitVisible("#media", chromedp.ByQuery), // Aguarda o carregamento da div
		chromedp.Sleep(200*time.Millisecond),
		chromedp.Screenshot("#media", &buf), // Captura a div específica
		// chromedp.CaptureScreenshot(&buf),
	)
	if err != nil {
		log.Fatalf("Erro ao tirar screenshot: %v", err)
	}
	matchScoreImagePath := fmt.Sprintf("output/matches-%s.png", strings.Join(matchNoList, "_"))

	err = os.WriteFile(matchScoreImagePath, buf, 0644)
	if err != nil {
		log.Fatalf("Erro ao salvar imagem: %v", err)
	}

	log.Printf("Imagem do placar gerada em %s.\n", matchScoreImagePath)
}
