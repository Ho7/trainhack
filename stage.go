package main

import (
	"time"
)

type Stage struct {
	Name string

	Game    *Game
	Hero    *Actor
	Actors  []*Actor
	Actions *Actions
	Level   *Level

	ticker *Ticker
}

func NewStage(g *Game) *Stage {
	hero := NewHero(DefaultWeapons)
	return &Stage{
		Game: g,
		Hero: hero,

		Actors:  []*Actor{hero},
		Actions: NewActions(),
	}
}

func (s *Stage) Load(name string, location *rune) *Stage {
	s.Stop()

	level := LoadLevel(s.Game, name)
	if level == nil {
		panic("cannot load level " + name)
		// TODO: handle error
		return s
	}

	// save state
	s.Save()

	// load or create state
	if state, ok := s.Game.State.Stages[name]; ok {
		// TODO: chage / cleanup actors from levelMap
		s.Actors = state.Actors
		// TODO: use target to locate or:
		s.Hero.Position = state.HeroPosition
	} else {
		// TODO: create actors from levelMap
		s.Hero = NewHero(DefaultWeapons)
		s.Hero.Position = Position{X: 20, Y: 10}
		s.Actors = []*Actor{s.Hero}
	}

	if location != nil {
		if pos, ok := level.Doors[*location]; ok {
			s.Hero.Position = pos
		}
	}

	s.Name = name
	s.Actions.Reset()

	s.Level = level

	return s
}

func (s *Stage) Save() string {
	if s.Name == "" {
		return ""
	}
	s.Game.State.Stages[s.Name] = StateStage{
		HeroPosition: s.Hero.Position,
		Actors:       s.Actors,
	}
	return s.Name
}

func (s *Stage) Start() {
	s.ticker = NewTicker(tickTimeout, func(d time.Duration) {
		if s.Update(d) {
			s.Game.View.Draw()
		}
	})
}

func (s *Stage) Stop() {
	if s.ticker != nil {
		s.ticker.Done()
	}
}

func (s *Stage) ActorAt(pos Position) *Actor {
	for _, a := range s.Actors {
		if a.Position == pos {
			return a
		}
	}

	return nil
}

func (s *Stage) AddActor(actor *Actor) {
	s.Actors = append(s.Actors, actor)
}

func (s *Stage) Update(d time.Duration) bool {

	if d > tickTimeout {
		d = tickTimeout
	}
	timeFactor := float64(d) / tickTimeoutF

	l := len(s.Actors)
	for i := 0; i < l; i++ {
		actor := s.Actors[i]
		if actor.Behavior == nil {
			continue
		}

		if actor.Energy.CanTakeTurn() || actor.Energy.Gain(timeFactor*actor.Speed) {
			action := actor.Behavior()
			if action != nil {
				s.Actions.Add(action)
			}
		}
	}

	needToDraw := false

	for {
		action := s.Actions.Get()
		if action == nil {
			break
		}

		result := action.Perform()

		for result.Alternative != nil {
			result = result.Alternative.Perform()
		}

		if result.Success {
			needToDraw = true

			if action.Actor != nil {
				action.Actor.Energy.Spend()
			}
		}
	}

	return needToDraw
}
