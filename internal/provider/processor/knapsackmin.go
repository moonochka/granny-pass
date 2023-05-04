package processor

import (
	"math"
)

// KnapsackMinTable
// maxLen - knapsack capacity = count of symbols in password = 24
// pathLen  - 1/pathLen because we should find a min distance
func (v *vocabulary) KnapsackMinTable(items []*wordMetric) [][]knapsack {
	var (
		b, bMin, bMinNew, bNew                                              knapsack
		newPathLen, gapPathLen, lenLeftover, bMinPathLen, bNewPathLen, iNew int
		invCur, invNew                                                      float32
		needContinue                                                        bool
	)

	// n - count words in file = count of items
	n := len(items)

	bp := make([][]knapsack, n+1)
	for i := range bp {
		bp[i] = make([]knapsack, v.maxLen+1)
	}

	for i := 0; i < n+1; i++ {
		for j := 0; j < v.maxLen+1; j++ {

			if i == 0 || j == 0 {
				//нулевую строку и столбец заполняем нулями
				bp[i][j] = knapsack{
					items:   nil,
					pathLen: 0,
					count:   0,
				}
			} else if i == 1 {
				//первая строка заполняется просто: первый предмет кладём или не кладём в зависимости от веса
				if items[0].len <= j {
					b = knapsack{
						items:   []*wordMetric{items[0]},
						pathLen: items[0].pathLen,
						count:   1,
					}
				} else {
					b = knapsack{
						items:   nil,
						pathLen: 0,
						count:   0,
					}
				}

				bp[i][j] = b
			} else {
				//если очередной предмет не влезает в рюкзак,
				if items[i-1].len > j {
					//записываем предыдущий максимум
					bp[i][j] = bp[i-1][j]
				} else {
					//рассчитаем длину пути очередного слова +
					//+ максимальную длину пути для слова, подходящего по длине оставшимся символам
					//(максимально возможный для текущего рюкзака вес − вес предмета) +
					//+ длина пути между последней буквой имеющегося и первой нового

					//длина остатка = (максимально возможное для текущего рюкзака число букв − число букв текущего слова)
					lenLeftover = j - items[i-1].len

					if lenLeftover > 0 {
						//выберем лучшее слово для добивки оставшихся символов, учитывая расстояние между словами
						// инициализация параметров:
						iNew = i - 1
						for l := 0; (i - 1 - l) > 0; l++ {
							bMin = bp[i-1-l][lenLeftover]
							//элемент подходит, если в нем не максимум слов
							if bMin.count < v.wordCnt {
								iNew = i - 1 - l
								break
							}
						}

						gapPathLen, _ = v.GapPathLen(bMin.lastWord(), items[i-1].word)
						bMinPathLen = bMin.pathLen + gapPathLen

						//решаем нужно ли продолжать поиск
						needContinue = false

						//если его длина меньше, а это лучший кандидат, то дальше искать бессмысленно
						//и если расстояние до слова недостаточно хорошее
						if !(len(bMin.GetDescription()) < lenLeftover) && !isGood(gapPathLen) {
							needContinue = true
						}

						if needContinue {
							for k := 1; (iNew - k) > 1; k++ {
								// берем элемент из предыдущей строки
								bNew = bp[iNew-k][lenLeftover]
								if len(bNew.GetDescription()) < lenLeftover {
									//если его длина меньше, то дальше искать бессмысленно
									break
								} else {
									gapPathLen, _ = v.GapPathLen(bNew.lastWord(), items[i-1].word)
									if isGood(gapPathLen) {
										//если расстояние до слова хорошее
										bNewPathLen = bNew.pathLen + gapPathLen
										if bNewPathLen < bMinPathLen {
											//и если суммарная длина пути меньше, то останавливаем поиск
											bMinPathLen = bNewPathLen
											bMin = bNew
											break
										}
									}
								}
							}
						}

						newPathLen = items[i-1].pathLen + bMinPathLen

						newItems := append([]*wordMetric{}, bMin.items...)
						bMinNew = knapsack{
							items:   append(newItems, items[i-1]),
							pathLen: newPathLen,
							count:   bMin.count + 1,
						}

					} else {
						//если добивать не надо, минимальный рюкзак bMin пустой, кладем туда слово, а длина всего пути = пути слова
						bMinNew = knapsack{
							items:   []*wordMetric{items[i-1]},
							pathLen: items[i-1].pathLen,
							count:   1,
						}
					}

					//смотрим длины элементов сверху и слева, тк по столбцу и строке длина не должна уменьшаться
					bMinNew = v.MaxLenKnapsack(bMinNew, bp[i-1][j], bp[i][j-1])

					// если длина bMinNew больше длины элемента сверху, пишем bMinNew
					if len(bMinNew.GetDescription()) > len(bp[i-1][j].GetDescription()) {
						bp[i][j] = bMinNew
					} else {
						//иначе длины равны, считаем обратные (inversion) величины для нахождения минимальной длины пути
						if bp[i-1][j].pathLen > 0 {
							invCur = 1 / float32(bp[i-1][j].pathLen)
						} else {
							invCur = 0
						}

						if bMinNew.pathLen > 0 {
							invNew = 1 / float32(bMinNew.pathLen)
						} else {
							invNew = 0
						}

						//сравниваем обратные значения
						//если предыдущий максимум больше
						if invCur > invNew {
							//запишем его
							bp[i][j] = bp[i-1][j]
						} else {
							//иначе фиксируем новый максимум
							bp[i][j] = bMinNew
						}
					}
				}
			}

		}
	}

	return bp
}

func (v *vocabulary) MinChoice(bc [][]knapsack) (knapsack, int) {
	var minKnapsack knapsack

	n := len(bc)
	k := len(bc[0])

	minPathLen := math.MaxInt

	for i := 0; i < n; i++ {
		for j := v.minLen; j < k; j++ {
			if bc[i][j].pathLen < minPathLen && bc[i][j].count == v.wordCnt && len(bc[i][j].GetDescription()) >= v.minLen {
				//fmt.Printf("i=%v j=%v \n", i, j)
				minPathLen = bc[i][j].pathLen
				minKnapsack = bc[i][j]
			}
		}
	}

	return minKnapsack, minPathLen
}

func isGood(n int) bool {
	if n <= 2 {
		return true
	}
	return false
}

func (v *vocabulary) MaxLenKnapsack(b1, b2, b3 knapsack) knapsack {
	var b, bMax knapsack
	bMax = b1

	for _, b = range []knapsack{b2, b3} {
		if len(bMax.GetDescription()) < len(b.GetDescription()) {
			bMax = b
		} else if len(bMax.GetDescription()) == len(b.GetDescription()) && bMax.pathLen > b.pathLen {
			bMax = b
		}
	}
	return bMax
}
