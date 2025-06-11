package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gomesar9/bvb-core/public/notify"
	"github.com/gomesar9/bvb-scoreboard/pkg/adapter/html"
)

func dumpData(data []notify.MatchInfo) {
	// Serializa para JSON
	jsonData, err := json.MarshalIndent(data, "", "  ") // Indentado (mais legível)
	if err != nil {
		panic(err)
	}

	// Cria (ou sobrescreve) o arquivo
	fileName := fmt.Sprintf("%s.json", time.Now().Format("2006-01-02_15-04-05"))
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Escreve os dados no arquivo
	_, err = file.Write(jsonData)
	if err != nil {
		panic(err)
	}
}

func teamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var reader io.ReadCloser
	var err error

	// Verifica se o corpo está compactado com gzip
	if r.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "Erro ao descompactar gzip", http.StatusBadRequest)
			return
		}
		defer reader.Close()
	} else {
		reader = r.Body
	}
	defer r.Body.Close()

	var matchesInfo []notify.MatchInfo
	if err := json.NewDecoder(reader).Decode(&matchesInfo); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Recebido: %+v\n", len(matchesInfo))

	go func() {
		// dumpData(matchesInfo)
		html.MakeInstagramStoryMedia(matchesInfo)
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Dados recebidos com sucesso"))
}

func main() {
	// Serve arquivos estáticos (imagens, CSS, JS etc)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	ft := http.FileServer(http.Dir("test"))
	http.Handle("/test/", http.StripPrefix("/test/", ft))

	http.HandleFunc("/team", teamHandler)

	log.Println("Servidor ouvindo na porta 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
