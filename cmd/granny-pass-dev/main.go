package main

import (
	"flag"
	"fmt"
	"granny-pass/internal/provider/graph"
	"granny-pass/internal/provider/processor"
	"log"
)

const (
	maxKeyboardPathLen = 20
	vocabularyDir      = "vocabularies/"
	distMapDir         = "distanceMaps/"
	distMapFilePrefix  = "dm"

	defaultMinPasswordLen = 20
	defaultMaxPasswordLen = 24
	defaultWordCnt        = 4
	defaultVocabularyFile = "short.txt"
)

func main() {
	var (
		minLen, maxLen, wordCnt     int
		useNormalizedKeyboard, help bool
		vocFile                     string
	)

	flag.BoolVar(&help, "help", false, "Help")
	flag.IntVar(&minLen, "min", defaultMinPasswordLen, "Provide minimum length of password")
	flag.IntVar(&maxLen, "max", defaultMaxPasswordLen, "Provide maximum length of password")
	flag.IntVar(&wordCnt, "cnt", defaultWordCnt, "Count of words")
	flag.BoolVar(&useNormalizedKeyboard, "k", false, "Use normalized keyboard - natural movement of one-finger typing method. By default will use keyboard from the task: only horizontal and vertical connections of buttons")
	flag.StringVar(&vocFile, "file", defaultVocabularyFile, "Vocabulary file name. Should consist of low-case words, no numbers, no special symbols. New line separator")

	flag.Parse()

	//TODO: cnt>=2
	if help {
		flag.PrintDefaults()
	} else {
		fmt.Println("Generating password for a grandmother. Parameters:")
		fmt.Printf(" min lenth: %d \n max lenth: %d \n count of words: %d \n", minLen, maxLen, wordCnt)
		fmt.Printf(" vocabulary file: %s \n", vocabularyDir+vocFile)
		if useNormalizedKeyboard {
			fmt.Println(" with normalized keyboard")
		} else {
			fmt.Println(" with keyboard from task")
		}

		m, err := GetBigramDistanceMap(useNormalizedKeyboard)
		if err != nil {
			log.Fatal(err)
		}

		p := processor.NewVocab(m, minLen, maxLen, uint8(wordCnt))

		wm, err := p.ReadFile(vocabularyDir+vocFile, true)
		if err != nil {
			log.Fatal(err)
		}
		kt := p.NewKnapsackTable(wm)
		k, pathLen := p.MinChoice(kt)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\nRESULT:\n%s \n used words: %s, lenth: %d, path lenth: %d\n", k.GetDescription(), k.GetDescriptionWithSpace(), len(k.GetDescription()), pathLen)

		/*
			m := PrepareDistMap(useNormalizedKeyboard)

			p := processor.New(m, minLen, maxLen, wordCnt)
			wm, err := p.ReadFile(vocabularyDir+vocFile, true)
			if err != nil {
				log.Fatal(err)
			}
			bt := p.KnapsackMinTable(wm)
			b, pathLen := p.MinChoice(bt)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("\nRESULT:\n%s \n used words: %s, lenth: %d, path lenth: %d\n", b.GetDescription(), b.GetDescriptionWithSpace(), len(b.GetDescription()), pathLen)
		*/
	}
}

func GetBigramDistanceMap(useNormalizedKeyboard bool) ([]int, error) {
	var (
		//err      error
		m []int
		//filename string
	)

	//filename = distMapDir + distMapFilePrefix + ".json"
	//if useNormalizedKeyboard {
	//	filename = distMapDir + distMapFilePrefix + "_norm.json"
	//}

	//if _, err = os.Stat(filename); err == nil {
	//	m, err = graph.ReadFromJson(filename)
	//	if err != nil {
	//		return nil, err
	//	}
	//} else {
	dist := PrepareDistMap(useNormalizedKeyboard)
	m = graph.BigramDistanceArray(dist)
	//err = graph.SaveToJson(m, filename)
	//if err != nil {
	//	fmt.Printf("%v", err)
	//}
	//}
	return m, nil
}

func PrepareDistMap(useNormalizedKeyboard bool) map[string]map[string]int {
	hash := func(v graph.Vertex) string {
		return v.Name
	}
	g := graph.New(hash)

	//add all key buttons
	for r := 'a'; r <= 'z'; r++ {
		_ = g.AddVertex(graph.Vertex{Name: string(r)})
	}

	//add all connections, weight=1
	//row1
	_ = g.AddEdge("q", "w")
	_ = g.AddEdge("w", "e")
	_ = g.AddEdge("e", "r")
	_ = g.AddEdge("r", "t")
	_ = g.AddEdge("t", "y")
	_ = g.AddEdge("y", "u")
	_ = g.AddEdge("u", "i")
	_ = g.AddEdge("i", "o")
	_ = g.AddEdge("o", "p")

	//row2
	_ = g.AddEdge("a", "s")
	_ = g.AddEdge("s", "d")
	_ = g.AddEdge("d", "f")
	_ = g.AddEdge("f", "g")
	_ = g.AddEdge("g", "h")
	_ = g.AddEdge("h", "j")
	_ = g.AddEdge("j", "k")
	_ = g.AddEdge("k", "l")

	//row3
	_ = g.AddEdge("z", "x")
	_ = g.AddEdge("x", "c")
	_ = g.AddEdge("c", "v")
	_ = g.AddEdge("v", "b")
	_ = g.AddEdge("b", "n")
	_ = g.AddEdge("n", "m")

	if useNormalizedKeyboard {
		//natural movement of one-finger typing method
		_ = g.AddEdge("q", "a")
		_ = g.AddEdge("w", "a")
		_ = g.AddEdge("w", "s")
		_ = g.AddEdge("e", "s")
		_ = g.AddEdge("e", "d")
		_ = g.AddEdge("r", "d")
		_ = g.AddEdge("r", "f")
		_ = g.AddEdge("t", "f")
		_ = g.AddEdge("t", "g")
		_ = g.AddEdge("y", "g")
		_ = g.AddEdge("y", "h")
		_ = g.AddEdge("u", "h")
		_ = g.AddEdge("u", "j")
		_ = g.AddEdge("i", "j")
		_ = g.AddEdge("i", "k")
		_ = g.AddEdge("o", "k")
		_ = g.AddEdge("o", "l")
		_ = g.AddEdge("p", "l")

		_ = g.AddEdge("a", "z")
		_ = g.AddEdge("s", "z")
		_ = g.AddEdge("s", "x")
		_ = g.AddEdge("d", "x")
		_ = g.AddEdge("d", "c")
		_ = g.AddEdge("f", "c")
		_ = g.AddEdge("f", "v")
		_ = g.AddEdge("g", "v")
		_ = g.AddEdge("g", "b")
		_ = g.AddEdge("h", "b")
		_ = g.AddEdge("h", "n")
		_ = g.AddEdge("j", "n")
		_ = g.AddEdge("j", "m")
		_ = g.AddEdge("k", "m")
	} else {
		//vertical only
		//left part
		_ = g.AddEdge("q", "a")
		_ = g.AddEdge("a", "z")

		_ = g.AddEdge("w", "s")
		_ = g.AddEdge("s", "x")

		_ = g.AddEdge("e", "d")
		_ = g.AddEdge("d", "c")

		_ = g.AddEdge("r", "f")
		_ = g.AddEdge("f", "v")

		_ = g.AddEdge("t", "g")
		_ = g.AddEdge("g", "b")

		//right part
		_ = g.AddEdge("y", "h")
		_ = g.AddEdge("h", "n")

		_ = g.AddEdge("u", "j")
		_ = g.AddEdge("j", "m")

		_ = g.AddEdge("i", "k")

		_ = g.AddEdge("o", "l")
	}
	m, _ := g.WFI(maxKeyboardPathLen)
	return m
}
