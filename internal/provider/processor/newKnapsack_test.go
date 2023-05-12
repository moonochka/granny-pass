package processor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"granny-pass/internal/provider/graph"
)

type testParam struct {
	fileName string
	minLen   int
	maxLen   int
	wordCnt  int
}

func TestNewKnapsack(t *testing.T) {
	t.Run("test NewKnapsack functions", func(t *testing.T) {

		dist := getDistanceMapForTests()
		v := NewVocab(dist, 0, 0, 0)

		t.Run("FindBestCombination", func(t *testing.T) {

			var (
				err         error
				k, kNew     knapsack
				wm1, wm2    *wordMetric
				wordMetrics []*wordMetric
				needStop    bool
			)

			t.Run("Empty knapsack & empty word", func(t *testing.T) {
				needStop, kNew, err = v.FindBestCombination(k, wm1)
				assert.Error(t, err)
			})

			t.Run("Empty knapsack & filled word", func(t *testing.T) {
				wm1 = &wordMetric{
					word:    "zas",
					pathLen: 2,
				}

				needStop, kNew, err = v.FindBestCombination(k, wm1)
				assert.NoError(t, err)
				assert.Equal(t, "zas", kNew.GetDescription())
				assert.Equal(t, true, needStop)

			})

			t.Run("Knapsack from test1.txt & empty word", func(t *testing.T) {
				wordMetrics, err = v.ReadFile("tests/test1.txt", false)
				assert.NoError(t, err)

				k = knapsack{
					items: wordMetrics,
				}

				needStop, kNew, err = v.FindBestCombination(k, wm2)
				assert.Error(t, err)
			})

			t.Run("Filled knapsack & filled word", func(t *testing.T) {
				needStop, kNew, err = v.FindBestCombination(k, wm1)
				assert.NoError(t, err)
				assert.Equal(t, "zasaofthebike", kNew.GetDescription())
				assert.Equal(t, false, needStop)

				wm2 = &wordMetric{
					word:    "eew",
					pathLen: 2,
				}

				needStop, kNew, err = v.FindBestCombination(k, wm2)
				assert.NoError(t, err)
				assert.Equal(t, "aofthebikeeew", kNew.GetDescription())
				assert.Equal(t, true, needStop)
			})

		})

		t.Run("ChooseCandidate", func(t *testing.T) {

			k11 := knapsack{
				items:   []*wordMetric{{word: "a", pathLen: 0}},
				pathLen: 0,
				count:   1,
			}
			k21 := knapsack{
				items:   []*wordMetric{{word: "of", pathLen: 5}},
				pathLen: 5,
				count:   1,
			}
			k31 := knapsack{
				items:   []*wordMetric{{word: "is", pathLen: 6}},
				pathLen: 6,
				count:   1,
			}

			k12 := knapsack{
				items:   []*wordMetric{{word: "a", pathLen: 0}, {word: "of", pathLen: 5}},
				pathLen: 13,
				count:   2,
			}
			k22 := knapsack{
				items:   []*wordMetric{{word: "of", pathLen: 5}, {word: "to", pathLen: 4}},
				pathLen: 13,
				count:   2,
			}
			k32 := knapsack{
				items:   []*wordMetric{{word: "to", pathLen: 4}, {word: "of", pathLen: 5}},
				pathLen: 9,
				count:   2,
			}

			ks1 := make(map[uint8]knapsack)
			ks2 := make(map[uint8]knapsack)
			ks3 := make(map[uint8]knapsack)

			ks1[1] = k11
			ks1[2] = k12

			ks2[1] = k21
			ks3[2] = k22

			ks3[1] = k31
			ks3[2] = k32

			t.Run("1,2,3", func(t *testing.T) {
				ksRes := v.ChooseCandidate(ks1, ks2, ks3)
				assert.Equal(t, k21, ksRes[1])
				assert.Equal(t, k32, ksRes[2])
			})

			k41 := knapsack{
				items:   []*wordMetric{{word: "the", pathLen: 6}},
				pathLen: 6,
				count:   1,
			}
			ks4 := make(map[uint8]knapsack)
			ks4[1] = k41

			t.Run("4,2,3", func(t *testing.T) {
				ksRes := v.ChooseCandidate(ks4, ks2, ks3)
				assert.Equal(t, k41, ksRes[1])
				_, ok := ksRes[2]
				assert.Equal(t, false, ok)
			})

		})

		t.Run("KnapsackTable and MinChoice", func(t *testing.T) {
			var (
				maxPathLen, p int
				err           error
				wordMetrics   []*wordMetric
				k             knapsack
				kt            *[][]map[uint8]knapsack
			)

			tests := []testParam{
				{
					fileName: "tests/test3.txt",
					minLen:   4,
					maxLen:   6,
					wordCnt:  2,
				},
				{
					fileName: "tests/out5a.txt",
					minLen:   20,
					maxLen:   24,
					wordCnt:  4,
				},
				//{
				//	fileName: "tests/out5b.txt",
				//	minLen:   20,
				//	maxLen:   24,
				//	wordCnt:  4,
				//},
			}

			for i, param := range tests {
				dist = getDistanceMapForTests()
				t.Run(fmt.Sprintf("Test %d, from file %s", i, param.fileName), func(t *testing.T) {

					v = NewVocab(dist, param.minLen, param.maxLen, uint8(param.wordCnt))
					wordMetrics, err = v.ReadFile(param.fileName, true)
					assert.NoError(t, err)

					n := len(wordMetrics)
					kt = v.NewKnapsackTable(wordMetrics)

					if i == 0 {
						printMap(kt, n, param.maxLen, uint8(param.wordCnt))
					}

					for n1, x := range *kt {
						for n2, y := range x {
							for cnt, k1 := range y {
								// кол-во слов = индексу map
								assert.Equal(t, uint8(len(k1.items)), cnt)

								p, err = v.PathLen(k1.GetDescription())
								assert.NoError(t, err)
								if p != k1.pathLen {
									fmt.Printf("AAAA! [%v %v] %v \n %v", n1, n2, k1.GetDescriptionWithSpace(), k1)
								}
								//проверяем правильность pathLen
								assert.Equal(t, p, k1.pathLen)

								//длина слова не больше номера столбца
								assert.Equal(t, true, len(k1.GetDescription()) <= n2)

								// по столбцу значение pathLen не может увеличиваться, только если слово удлиннилось
								if n1 != 0 && n2 != 0 {
									k2 := (*kt)[n1-1][n2][cnt]
									assert.Equal(t, true, len(k1.GetDescription()) >= len(k2.GetDescription()))
									if len(k1.GetDescription()) == len(k2.GetDescription()) {
										assert.Equal(t, true, k1.pathLen <= k2.pathLen)
									}
								}
							}
						}
					}

					k, maxPathLen = v.MinChoice(kt)
					fmt.Printf("pathLen: %d\npassword: %s\nwords: %s\n", maxPathLen, k.GetDescription(), k.GetDescriptionWithSpace())

					if i == 0 {
						assert.Equal(t, 4, maxPathLen)
						assert.Equal(t, "ploki", k.GetDescription())
					}
				})
			}
		})

		/////

		t.Run("", func(t *testing.T) {

		})
	})
}

func getDistanceMapForTests() []int {
	hash := func(v graph.Vertex) string {
		return v.Name
	}
	g := graph.New(hash)

	//add all key buttons
	for r := 'a'; r <= 'z'; r++ {
		_ = g.AddVertex(graph.Vertex{Name: string(r)})
	}

	//add all connections weight=1
	_ = g.AddEdge("q", "w")
	_ = g.AddEdge("w", "e")
	_ = g.AddEdge("e", "r")
	_ = g.AddEdge("r", "t")
	_ = g.AddEdge("t", "y")
	_ = g.AddEdge("y", "u")
	_ = g.AddEdge("u", "i")
	_ = g.AddEdge("i", "o")
	_ = g.AddEdge("o", "p")

	_ = g.AddEdge("a", "s")
	_ = g.AddEdge("s", "d")
	_ = g.AddEdge("d", "f")
	_ = g.AddEdge("f", "g")
	_ = g.AddEdge("g", "h")
	_ = g.AddEdge("h", "j")
	_ = g.AddEdge("j", "k")
	_ = g.AddEdge("k", "l")

	_ = g.AddEdge("z", "x")
	_ = g.AddEdge("x", "c")
	_ = g.AddEdge("c", "v")
	_ = g.AddEdge("v", "b")
	_ = g.AddEdge("b", "n")
	_ = g.AddEdge("n", "m")

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

	m, _ := g.WFI(20)
	dist := graph.BigramDistanceArray(m)
	return dist
}

func printMap(kt *[][]map[uint8]knapsack, n, k int, wordCnt uint8) {
	for i := 0; i < n+1; i++ {
		for j := 0; j < k+1; j++ {
			s := ""
			for cnt := uint8(1); cnt <= wordCnt; cnt++ {
				if b, ok := (*kt)[i][j][uint8(cnt)]; ok {
					s = fmt.Sprintf("%s[%d:%d=%v]", s, cnt, b.pathLen, b.GetDescription())
				} else {
					s = fmt.Sprintf("%s[%d:_=_]", s, cnt)
				}
			}
			fmt.Printf("%s \t", s)
		}
		fmt.Println()
	}
}
