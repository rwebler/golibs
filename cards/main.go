package main

func main() {
	cards := newDeck()

	//hand, pile := deal(cards, 2)

	//hand.saveToFile("hand.bin")

	//pile.print()

	//cards := newDeckFromFile("hand.bin")

	cards.shuffle()

	cards.print()
}
