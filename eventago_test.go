package eventago

import (
	"fmt"
	"testing"
)

func TestOn(t *testing.T) {
	// Adds two listeners for "some_event"
	On("some_event", boolRes)
	On("some_event", intRes)

	// Runs GoFire to run in full concurrency mode
	GoFire("some_event", 123)

	// Checks if "some_event" has listeners Fire it and log result of them
	if HaveListener("some_event") {
		fireRes := Fire("some_event", 124)
		fmt.Println("fireRes:", fireRes)
	}

	fmt.Println("Count of listeners  :", CountListeners("some_event"))
	fmt.Println("Count of listeners 2:", CountListeners("some_event2"))

}

func boolRes(number int) bool {
	fmt.Println("isMosbat", number)
	return number > 0
}

func intRes(number int) int {
	fmt.Println("zarbdarer", number)
	return number * 2
}
