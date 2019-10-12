package main

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func Pursue(actor *Actor, stage *Stage, target *Actor) *Actor {

	actor.Behavior = func() *Action {
		dx := target.Position.X - actor.Position.X
		dy := target.Position.Y - actor.Position.Y
		dxAbs := abs(dx)
		dyAbs := abs(dy)

		// do not follow at 5 distance
		if dxAbs > 5 || dyAbs > 5 {
			return nil
		}

		if dxAbs > dyAbs {
			if dx < 0 {
				return ActionMove(stage, actor, DirectionLeft)
			} else if dy > 0 {
				return ActionMove(stage, actor, DirectionRight)
			}
		} else {
			if dy < 0 {
				return ActionMove(stage, actor, DirectionTop)
			} else if dy > 0 {
				return ActionMove(stage, actor, DirectionDown)
			}
		}

		return nil
	}

	return actor
}
