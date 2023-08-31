package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type App struct {
	agents  map[int]*Agent
	points  []Point
	lock    *sync.Mutex
	allStop chan struct{}
}

type Agent struct {
	Id        int
	x         float64
	y         float64
	isRunning bool
	remained  float64
}

type Point struct {
	Id int
	x  float64
	y  float64
}

func main() {
	app := new(App)
	app.lock = new(sync.Mutex)
	app.agents = make(map[int]*Agent)
	app.allStop = make(chan struct{})
	for i := 1; i <= 5; i++ {
		app.agents[i] = &Agent{Id: i}
	}

	points := []Point{
		{
			Id: 1,
			x:  2,
			y:  3,
		},
		{
			Id: 2,
			x:  4,
			y:  5,
		},
		{
			Id: 3,
			x:  -3,
			y:  -4,
		},
		{
			Id: 4,
			x:  -6,
			y:  -7,
		},
		{
			Id: 5,
			x:  6,
			y:  8,
		},
		{
			Id: 6,
			x:  -10,
			y:  -12,
		},
		{
			Id: 7,
			x:  4,
			y:  2,
		},
		{
			Id: 8,
			x:  -3,
			y:  4,
		},
		{
			Id: 9,
			x:  -8,
			y:  2,
		},
		{
			Id: 10,
			x:  5,
			y:  -6,
		},
		{
			Id: 11,
			x:  6,
			y:  -8,
		},
	}

	for i, p := range points {
		agentCh := make(chan *Agent, 1)
		app.findNearestAgentToPoint(p, agentCh)
		agent := <-agentCh
		go app.moveToPoint(agent, p)

		fmt.Printf("agent %d start going to point %d\n", agent.Id, p.Id)
		// wait for all agents to stop and after pick nearest agent to next point
		if (i+1)%5 == 0 {
			<-app.allStop
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c

}

func (app *App) findNearestAgentToPoint(point Point, ag chan *Agent) {
	var minDistance float64 = 100
	near := new(Agent)

	app.lock.Lock()
	for _, agent := range app.agents {
		d := getDistance(agent, point)
		if d <= minDistance {
			minDistance = d
			agent.remained = d
			near = agent
		}
	}

	delete(app.agents, near.Id)
	app.lock.Unlock()

	ag <- near
}

func getDistance(agent *Agent, point Point) float64 {
	distance := math.Pow(point.y-agent.y, 2) + math.Pow(point.x-agent.x, 2)
	return math.Sqrt(distance)
}

func (app *App) moveToPoint(agent *Agent, point Point) {
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-tick.C:
			agent.remained--
			if agent.remained <= 0 {
				agent.remained = 0
				agent.x = point.x
				agent.y = point.y
				fmt.Printf("agent %d is reached to point %d\n", agent.Id, point.Id)
				fmt.Printf("agent %d is now on x:%f, y:%f\n", agent.Id, agent.x, agent.y)
				app.agents[agent.Id] = agent
				if len(app.agents) == 5 {
					app.allStop <- struct{}{}
				}
				return
			}
			fmt.Printf("agent %d is moving to point %d, remained %f\n", agent.Id, point.Id, agent.remained)
		}
	}
}
