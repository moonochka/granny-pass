package processor

import (
	"strings"
)

type wordMetric struct {
	word    string
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
	firstWord() string
	Length() int
}

func (b *knapsack) GetDescription() string {
	if len(b.items) == 0 {
		return ""
	}

	var strBuilder strings.Builder

	for _, item := range b.items {
		strBuilder.WriteString(item.word)
	}
	return strBuilder.String()

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

func (b *knapsack) firstWord() string {
	l := len(b.items)
	if l == 0 {
		return ""
	}

	return b.items[0].word
}

func (b *knapsack) lastWord() string {
	l := len(b.items)
	if l == 0 {
		return ""
	}

	return b.items[l-1].word
}

func (b *knapsack) Length() int {
	sum := 0
	for _, i := range b.items {
		sum += len((*i).word)
	}
	return sum
}
