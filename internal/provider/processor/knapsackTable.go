package processor

import (
	"fmt"
	"math"
	"sync"
)

func (v *vocab) KnapsackTable(items []*wordMetric) *[][][]knapsack {
	var wg sync.WaitGroup
	// n - count words in file = count of items
	n := len(items)

	kt := make([][][]knapsack, n+1)
	for i := range kt {
		kt[i] = make([][]knapsack, v.maxLen+1)
	}

	//нулевую строку и столбец заполняем нулями
	for i := 0; i < n+1; i++ {
		for j := 0; j < v.maxLen+1; j++ {
			if i == 0 || j == 0 {
				kt[i][j] = make([]knapsack, v.wordCnt+1)
			}
		}
	}

	for currRow := 1; currRow < n+1; currRow++ {
		num := int(math.Min(float64(currRow), float64(v.maxLen)))
		wg.Add(num)
		for i, j := currRow, 1; i >= 1 && j < v.maxLen+1; i, j = i-1, j+1 {
			i := i
			j := j
			item := items[i-1]
			go func() {
				defer wg.Done()
				//TODO: handle error
				_ = v.calcSet(i, j, item, &kt)

			}()
		}
		wg.Wait()
	}
	wg.Wait()

	//right bottom corner
	for currColm := 2; currColm < v.maxLen+1; currColm++ {
		num := v.maxLen - currColm + 1
		wg.Add(num)
		for i, j := n, currColm; j < v.maxLen+1; i, j = i-1, j+1 {
			i := i
			j := j
			item := items[i-1]
			go func() {
				defer wg.Done()
				_ = v.calcSet(i, j, item, &kt)

			}()
		}
		wg.Wait()
	}
	return &kt
}

// i - row of the kt table = word number
// j - column of the kt table = length of symbols in knapsacks
func (v *vocab) calcSet(i, j int, wm *wordMetric, kt *[][][]knapsack) error {
	var (
		setKnapsacks, candidateKnapsacks []knapsack
		kBest, kNew, kLeftover           knapsack
		needStop                         bool
		err                              error
		cnt                              uint8
		prevPass                         string
	)

	if len(wm.word) > j {
		//если очередной предмет не влезает в рюкзак, записываем предыдущий максимум
		setKnapsacks = (*kt)[i-1][j]
	} else {
		candidateKnapsacks = make([]knapsack, v.wordCnt+1)

		//слово подходит впритык или с запасом, фиксируем его как кандидата для наполнения рюкзака из 1 слова
		candidateKnapsacks[1] = knapsack{
			items:   []*wordMetric{wm},
			pathLen: wm.pathLen,
		}

		lenLeftover := j - len(wm.word)
		if lenLeftover > 0 {
			//выберем лучшее слово/слова для добивки оставшихся символов, учитывая расстояние между словами
			for cnt = 1; cnt < v.wordCnt; cnt++ {
				//добивка с количеством слов cnt
				kLeftover = (*kt)[i-1][lenLeftover][cnt]
				if kLeftover.isEmpty() {
					//такой нет, значит записываем предыдущий максимум
					if !(*kt)[i-1][j][cnt+1].isEmpty() {
						candidateKnapsacks[cnt+1] = (*kt)[i-1][j][cnt+1]
					}
					continue
				}

				needStop, kBest, err = v.FindBestCombination(kLeftover, wm)
				if err != nil {
					return err
				}

				prevPass = kLeftover.GetDescription()

				//двигаемся вверх по столбцу
				for i1 := i - 2; i1 > 0 && !needStop; i1-- {
					kLeftover = (*kt)[i1][lenLeftover][cnt]
					if kLeftover.isEmpty() {
						break
					}

					//если его длина меньше, то дальше искать бессмысленно
					if len(prevPass) > kLeftover.Length() {
						break
					}

					//если пароль тот же или длина уменьшилась - то считать незачем
					if prevPass == kLeftover.GetDescription() {
						continue
					}

					needStop, kNew, err = v.FindBestCombination(kLeftover, wm)
					if err != nil {
						return err
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
		}, nil

	}

	//add in the end
	newItems := append([]*wordMetric{}, k.items...)
	return g2 == 0, knapsack{
		items:   append(newItems, wm),
		pathLen: k.pathLen + g2 + wm.pathLen,
	}, nil

}

// ChooseCandidate compare length, if equal than compare pathLen
func (v *vocab) ChooseCandidate(candidateKs, upKs, leftKs []knapsack) []knapsack {
	bestKs := candidateKs

	for cnt, k := range candidateKs {
		// check only if it was set
		if k.isEmpty() {
			continue
		}

		//create list of knapsacks with cnt words
		var list []knapsack

		k1 := upKs[cnt]
		if !k1.isEmpty() {
			list = append(list, k1)
		}

		k2 := leftKs[cnt]
		if !k2.isEmpty() {
			list = append(list, k2)
		}

		for _, kl := range list {
			b := bestKs[cnt]
			if b.Length() < kl.Length() {
				bestKs[cnt] = kl
			} else if b.Length() == kl.Length() && bestKs[cnt].pathLen > kl.pathLen {
				bestKs[cnt] = kl
			}
		}
	}

	//если рюкзак из cnt слов есть в upKs, но нет в bestKs, добавляем
	for cnt := range upKs {
		if bestKs[cnt].isEmpty() && !upKs[cnt].isEmpty() {
			bestKs[cnt] = upKs[cnt]
		}
	}
	return bestKs
}

func (v *vocab) MinChoice(kt *[][][]knapsack) (knapsack, int) {
	var minKnapsack knapsack

	n := len(*kt)
	k := len((*kt)[0])

	minPathLen := math.MaxInt

	cnt := v.wordCnt

	for i := 0; i < n; i++ {
		for j := v.minLen; j < k; j++ {
			kn := (*kt)[i][j][cnt]
			if !kn.isEmpty() {
				if kn.pathLen < minPathLen && kn.Length() >= v.minLen {
					minPathLen = kn.pathLen
					minKnapsack = kn
				}
			}
		}
	}

	return minKnapsack, minPathLen
}
