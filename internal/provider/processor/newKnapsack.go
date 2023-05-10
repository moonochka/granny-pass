package processor

import "fmt"

// i - row of the kt table = word number
// j - column of the kt table = length of symbols in knapsacks
func (v *vocab) parallelCalc(i, j int, wm *wordMetric, kt *[][]map[uint8]knapsack) error {
	var (
		setKnapsacks, candidateKnapsacks map[uint8]knapsack
		kBest, kNew, kLeftover           knapsack
		ok                               bool
		err                              error
	)

	if len(wm.word) < j {
		//если очередной предмет не влезает в рюкзак, записываем предыдущий максимум
		setKnapsacks = (*kt)[i-1][j]
	} else {
		candidateKnapsacks = make(map[uint8]knapsack)

		lenLeftover := j - len(wm.word)
		if lenLeftover > 0 {
			//выберем лучшее слово/слова для добивки оставшихся символов, учитывая расстояние между словами
			for cnt := 1; cnt <= 3; cnt++ {
				//добивка с количеством слов cnt
				kLeftover, ok = (*kt)[i-1][lenLeftover][uint8(cnt)]
				if !ok {
					//такой нет, значит кандидата не составить
					continue
				}

				kBest, err = v.FindBestCombination(kLeftover, wm)
				if err != nil {
					return err
				}

				//двигаемся вверх по столбцу
				for i1 := i - 2; i1 > 0; i1-- {
					kLeftover, ok = (*kt)[i1][lenLeftover][uint8(cnt)]
					if !ok {
						break
					}

					kNew, err = v.FindBestCombination(kLeftover, wm)
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

				//фиксируем кандидата для наполнения рюкзака из cnt слов
				candidateKnapsacks[uint8(cnt)] = kBest
			}
		} else {
			//слово подходит впритык, фиксируем его как кандидата для наполнения рюкзака из 1 слова
			candidateKnapsacks[1] = knapsack{
				items:   []*wordMetric{wm},
				pathLen: wm.pathLen,
				count:   1,
			}
		}

	}

	(*kt)[i][j] = setKnapsacks
	return nil
}

// FindBestCombination insert word wm in different positions in knapsack k
// return new knapsack with the shortest pathLen
func (v *vocab) FindBestCombination(k knapsack, wm *wordMetric) (knapsack, error) {
	var (
		g1, g2 int
		err    error
	)

	if wm == nil {
		return k, fmt.Errorf("wm is not set, probably you got an empty line(word) in file")
	}

	//add in the front
	g1, err = v.GapPathLen(wm.word, k.firstWord())
	if err != nil {
		return k, err
	}

	//add in the end
	g2, err = v.GapPathLen(k.lastWord(), wm.word)
	if err != nil {
		return k, err
	}

	if g1 < g2 {
		//add in the front
		newItems := append([]*wordMetric{}, wm)
		return knapsack{
			items:   append(newItems, k.items...),
			pathLen: wm.pathLen + g1 + k.pathLen,
			count:   k.count + 1,
		}, nil

	}

	//add in the end
	newItems := append([]*wordMetric{}, k.items...)
	return knapsack{
		items:   append(newItems, wm),
		pathLen: k.pathLen + g1 + wm.pathLen,
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
	return bestKs
}
