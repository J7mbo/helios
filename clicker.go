package helios

import (
	"math"
	"math/rand"
	"time"

	"github.com/go-vgo/robotgo"
)

type Clicker struct {
}

// Click clicks at a random X and Y coordinate within the given location.
// This doesn't use the standard Click() function of any library which is easy
// to detect as an automated click - rather it adds entropy throughout the mouse events.
func (c *Clicker) Click(match *Match) {
	c.sleepRandomly(0.2, 0.5)
	c.moveMouseRandomlyWithinBox(match.topLeft.x, match.topLeft.y, float64(match.width), float64(match.height))
	c.sleepRandomly(0.2, 0.5)
	c.performRandomisedClick()
}

func (c *Clicker) MoveMouseInRegion(region *Region) {
	c.moveMouseRandomlyWithinBox(region.topLeft.x, region.topLeft.y, float64(region.width), float64(region.height))
}

func (c *Clicker) performRandomisedClick() {
	robotgo.Toggle("left", "down")
	c.sleepRandomly(0.2, 0.5)
	robotgo.Toggle("left", "up")
}

func (c *Clicker) sleepRandomly(min, max float64) {
	time.Sleep(time.Duration(c.generateRandomNumber(min, max) * float64(time.Second)))
}

func (c *Clicker) moveMouseRandomlyWithinBox(x, y, w, h float64) {
	randomX := c.generateRandomNumber(x, x+w)
	randomY := c.generateRandomNumber(y, y+h)

	robotgo.Move(int(math.Round(randomX)), int(math.Round(randomY)))
}

func (c *Clicker) generateRandomNumber(min float64, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	randNum := (rand.Float64() * (max - min)) + min

	// Trims to two decimal places. Doesn't need to be perfect.
	return math.Floor(randNum*100) / 100
}
