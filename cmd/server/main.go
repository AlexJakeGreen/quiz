package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"math/rand"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

type GermanWord struct {
	Article string `json:"article"`
	Word string `json:"word"`
}

var words []GermanWord

func loadWords() {
	f, err := os.Open("./words.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		w := strings.Split(scanner.Text(), " ")
		word := GermanWord{
			Article: w[0],
			Word: w[1],
		}
		words = append(words, word)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rand.Seed(time.Now().Unix())
	loadWords()
}

func indexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func getWord(c *gin.Context) {
	word := words[rand.Intn(len(words))]
	c.JSON(http.StatusOK, gin.H{"word": word.Word})
}

func checkAnswer(c *gin.Context) {
	var json GermanWord
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	richtig := false
	var w GermanWord
	for _, w = range words {
		if json.Word == w.Word {
			richtig = (json.Article == w.Article)
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"richtig": richtig,
		"word": w.Word,
		"article": w.Article,
	})
}


// Run ...
func Run() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.LoadHTMLFiles("templates/index.html")

	router.Static("/static", "./static")

	router.GET("/", indexPage)
	router.GET("/quiz", getWord)
	router.POST("/quiz", checkAnswer)

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.Run(fmt.Sprintf(":%d", 3000))
}

func main() {
	Run()
}
