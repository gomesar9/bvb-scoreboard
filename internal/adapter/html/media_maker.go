package html

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	_ "path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gomesar9/bvb-scoreboard/internal/domain/model"
	"github.com/gomesar9/bvb-scoreboard/internal/domain/service"

	"github.com/chromedp/chromedp"
	"github.com/gomesar9/bvb-core/public/notify"
)

type MediaData struct {
	Scoreboards []model.Scoreboard
	Theme       string
}

// Estrutura para os parâmetros do HTML
type MatchData struct {
}

func MakeInstagramStoryMedia(matches []notify.MatchInfo) {
	var mediaData MediaData
	var matchNoList []string
	templateFile := "templates/scoreboard_elite16.html"

	// TODO: Tratar em batches
	if len(matches) > 4 {
		matches = matches[:4]
	}

	mediaData.Theme = string(service.ThemeDark)
	for _, m := range matches {
		data := model.MatchToScoreboard(m)
		mediaData.Scoreboards = append(mediaData.Scoreboards, *data)
		matchNoList = append(matchNoList, strconv.FormatUint(uint64(m.MatchNo), 10))
	}

	// Renderiza o HTML com dados
	tpl, err := template.ParseFiles(templateFile)
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
	// absPath, err := filepath.Abs(htmlPath)
	// if err != nil {
	// 	log.Fatalf("Erro ao obter caminho absoluto: %v", err)
	// }

	fullPath := "http://localhost:8080/" + htmlPath

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
