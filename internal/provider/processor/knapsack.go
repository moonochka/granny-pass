package processor

type wordMetric struct {
	word    string
	len     int
	pathLen int
}

type knapsack struct {
	items   []*wordMetric
	pathLen int // в классической интерпретации price - ценность предметов
	count   int // items count
}

type Features interface {
	GetDescription() string
	GetDescriptionWithSpace() string
	lastWord() string
}

func (b *knapsack) GetDescription() string {
	if len(b.items) == 0 {
		return ""
	}

	s := ""
	for _, item := range b.items {
		s += item.word
	}
	return s

}

func (b *knapsack) GetDescriptionWithSpace() string {
	if len(b.items) == 0 {
		return ""
	}

	s := ""
	for _, item := range b.items {
		s += item.word + " "
	}
	return s

}

func (b *knapsack) lastWord() string {
	l := len(b.items)
	if l == 0 {
		return ""
	}

	return b.items[l-1].word
}

// KnapsackTable
// n - count words in file = count of items
// k - knapsack capacity = count of symbols in password = 24
func (v *vocabulary) KnapsackTable(k int, items []*wordMetric) [][]knapsack {
	var (
		b          knapsack
		newPathLen int
	)

	n := len(items)

	bp := make([][]knapsack, n+1)
	for i := range bp {
		bp[i] = make([]knapsack, k+1)
	}

	for i := 0; i < n+1; i++ {
		for j := 0; j < k+1; j++ {

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
					//рассчитаем цену очередного предмета + максимальную цену для
					//(максимально возможный для рюкзака вес − вес предмета)
					newPathLen = items[i-1].pathLen + bp[i-1][j-items[i-1].len].pathLen

					//если предыдущий максимум больше
					if bp[i-1][j].pathLen > newPathLen {
						//запишем его
						bp[i][j] = bp[i-1][j]
					} else {
						//иначе фиксируем новый максимум: текущий предмет + стоимость свободного пространства
						bp[i][j] = knapsack{
							items:   append(bp[i-1][j-items[i-1].len].items, items[i-1]),
							pathLen: newPathLen,
							count:   bp[i-1][j-items[i-1].len].count + 1,
						}
					}
				}
			}

		}
	}

	return bp
}

func (v *vocabulary) MaxChoice(bc [][]knapsack) (knapsack, int) {
	n := len(bc)
	k := len(bc[0])

	maxKnapsnack := bc[0][0]
	maxPathLen := 0

	for i := 0; i < n; i++ {
		if bc[i][k-1].pathLen > maxPathLen {
			//fmt.Printf("i=%v \n", i)
			maxPathLen = bc[i][k-1].pathLen
			maxKnapsnack = bc[i][k-1]

		}
	}

	//fmt.Printf("%v \n", maxKnapsnack.pathLen)
	//fmt.Printf("%v \n", maxKnapsnack.GetDescription())

	return maxKnapsnack, maxPathLen
}
