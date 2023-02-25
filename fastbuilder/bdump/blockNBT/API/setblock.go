package blockNBT_API

import (
	"fmt"
	"phoenixbuilder/fastbuilder/commands_generator"
	"phoenixbuilder/fastbuilder/types"
)

func (g *GlobalAPI) SetBlock(pos [3]int32, name string, states string) error {
	request := commands_generator.SetBlockRequest(&types.Module{
		Block: &types.Block{
			Name:        &name,
			BlockStates: states,
		},
		Point: types.Position{
			X: int(pos[0]),
			Y: int(pos[1]),
			Z: int(pos[2]),
		},
	}, &types.MainConfig{})
	_, err := g.SendWSCommandWithResponce(request)
	if err != nil {
		return fmt.Errorf("SetBlock: %v", err)
	}
	return nil
}

func (g *GlobalAPI) SetBlockFastly(pos [3]int32, name string, states string) error {
	request := commands_generator.SetBlockRequest(&types.Module{
		Block: &types.Block{
			Name:        &name,
			BlockStates: states,
		},
		Point: types.Position{
			X: int(pos[0]),
			Y: int(pos[1]),
			Z: int(pos[2]),
		},
	}, &types.MainConfig{})
	err := g.SendSettingsCommand(request, true)
	if err != nil {
		return fmt.Errorf("SetBlockFastly: %v", err)
	}
	return nil
}
