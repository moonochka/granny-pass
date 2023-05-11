//go:build processorTest
// +build processorTest

package processor

import (
	"fmt"
	"math"
)

func (v *vocab) NewKnapsackTable(items []*wordMetric) *[][]map[uint8]knapsack {

	// n - count words in file = count of items
	n := len(items)

	kt := make([][]map[uint8]knapsack, n+1)
	for i := range kt {
		kt[i] = make([]map[uint8]knapsack, v.maxLen+1)
	}

	for i := 0; i < n+1; i++ {
		for j := 0; j < v.maxLen+1; j++ {
			if i == 0 || j == 0 {
				//нулевую строку и столбец заполняем нулями
				kt[i][j] = make(map[uint8]knapsack)
			} else {
				//TODO: handle error
				_ = v.calcSet(i, j, items[i-1], &kt)
			}
		}
	}

	return &kt
}

// i - row of the kt table = word number
// j - column of the kt table = length of symbols in knapsacks
func (v *vocab) calcSet(i, j int, wm *wordMetric, kt *[][]map[uint8]knapsack) error {
	var (
		setKnapsacks, candidateKnapsacks map[uint8]knapsack
		kBest, kNew, kLeftover           knapsack
		ok, needStop                     bool
		err                              error
		cnt                              uint8
	)

	if len(wm.word) > j {
		//если очередной предмет не влезает в рюкзак, записываем предыдущий максимум
		setKnapsacks = (*kt)[i-1][j]
	} else {
		candidateKnapsacks = make(map[uint8]knapsack)

		//слово подходит впритык или с запасом, фиксируем его как кандидата для наполнения рюкзака из 1 слова
		candidateKnapsacks[1] = knapsack{
			items:   []*wordMetric{wm},
			pathLen: wm.pathLen,
			count:   1,
		}

		lenLeftover := j - len(wm.word)
		if lenLeftover > 0 {
			//выберем лучшее слово/слова для добивки оставшихся символов, учитывая расстояние между словами
			for cnt = 1; cnt < v.wordCnt; cnt++ {
				//добивка с количеством слов cnt
				kLeftover, ok = (*kt)[i-1][lenLeftover][cnt]
				if !ok {
					//такой нет, значит записываем предыдущий максимум
					if _, ok = (*kt)[i-1][j][cnt+1]; ok {
						candidateKnapsacks[cnt+1] = (*kt)[i-1][j][cnt+1]
					}
					continue
				}

				needStop, kBest, err = v.FindBestCombination(kLeftover, wm)
				if err != nil {
					return err
				}

				//двигаемся вверх по столбцу
				for i1 := i - 2; i1 > 0 && !needStop; i1-- {
					kLeftover, ok = (*kt)[i1][lenLeftover][cnt]
					if !ok {
						break
					}

					needStop, kNew, err = v.FindBestCombination(kLeftover, wm)
					if err != nil {
						return err
					}

					//TODO: тут можно сравнивать рюкзаки до добавления wn, если запоминать еще и kLeftover
					//если его длина меньше, то дальше искать бессмысленно
					if len(kBest.GetDescription()) > len(kNew.GetDescription()) {
						break
					}

					if kNew.pathLen < kBest.pathLen {
						kBest = kNew
					}
				}

				//фиксируем кандидата для наполнения рюкзака из cnt+1 слов
				candidateKnapsacks[cnt+1] = kBest
			}
		}

		setKnapsacks = v.ChooseCandidate(candidateKnapsacks, (*kt)[i-1][j], (*kt)[i][j-1])
	}

	(*kt)[i][j] = setKnapsacks
	return nil
}

// FindBestCombination insert word wm in different positions in knapsack k
// return bool flag nedStop when gap==0
// return new knapsack with the shortest pathLen
func (v *vocab) FindBestCombination(k knapsack, wm *wordMetric) (bool, knapsack, error) {
	var (
		g1, g2 int
		err    error
	)

	if wm == nil {
		return false, k, fmt.Errorf("wm is not set, probably you got an empty line(word) in file")
	}

	//add in the front
	g1, err = v.GapPathLen(wm.word, k.firstWord())
	if err != nil {
		return false, k, err
	}

	//add in the end
	g2, err = v.GapPathLen(k.lastWord(), wm.word)
	if err != nil {
		return false, k, err
	}

	if g1 < g2 {
		//add in the front
		newItems := append([]*wordMetric{}, wm)
		return g1 == 0, knapsack{
			items:   append(newItems, k.items...),
			pathLen: wm.pathLen + g1 + k.pathLen,
			count:   k.count + 1,
		}, nil

	}

	//add in the end
	newItems := append([]*wordMetric{}, k.items...)
	return g2 == 0, knapsack{
		items:   append(newItems, wm),
		pathLen: k.pathLen + g2 + wm.pathLen,
		count:   k.count + 1,
	}, nil

}

// ChooseCandidate compare length, if equal than compare pathLen
func (v *vocab) ChooseCandidate(candidateKs, upKs, leftKs map[uint8]knapsack) map[uint8]knapsack {
	bestKs := candidateKs

	for cnt := range candidateKs {
		//create list of knapsacks with cnt words
		var list []knapsack

		k1, ok := upKs[cnt]
		if ok {
			list = append(list, k1)
		}

		k2, ok := leftKs[cnt]
		if ok {
			list = append(list, k2)
		}

		for _, kl := range list {
			b := bestKs[cnt]
			if len(b.GetDescription()) < len(kl.GetDescription()) {
				bestKs[cnt] = kl
			} else if len(b.GetDescription()) == len(kl.GetDescription()) && bestKs[cnt].pathLen > kl.pathLen {
				bestKs[cnt] = kl
			}
		}
	}

	//если рюкзак из cnt слов есть в upKs, но нет в bestKs, добавляем
	for cnt := range upKs {
		_, ok := bestKs[cnt]
		if !ok {
			bestKs[cnt] = upKs[cnt]
		}
	}
	return bestKs
}

func (v *vocab) MinChoice(kt *[][]map[uint8]knapsack) (knapsack, int) {
	var minKnapsack knapsack

	n := len(*kt)
	k := len((*kt)[0])

	minPathLen := math.MaxInt

	cnt := v.wordCnt

	for i := 0; i < n; i++ {
		for j := v.minLen; j < k; j++ {
			kn, ok := (*kt)[i][j][cnt]
			if ok {
				if kn.pathLen < minPathLen && len(kn.GetDescription()) >= v.minLen {
					//fmt.Printf("i=%v j=%v \n", i, j)
					minPathLen = kn.pathLen
					minKnapsack = kn
				}
			}
		}
	}

	return minKnapsack, minPathLen
}
