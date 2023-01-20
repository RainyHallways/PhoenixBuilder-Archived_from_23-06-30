package Happy2018new_depends

import (
	_ "embed"
	"encoding/json"
)

type NEMCBlock struct {
	Name    string
	Val     int
	NEMCRID int
}

//go:embed runtimeIds_2_5_15_proc.json
var blockRunTimeIdTable []byte

func readNemcData() []NEMCBlock {
	NewNEMCBlock := func(p [2]interface{}, nemcRID int) NEMCBlock {
		s, ok := p[0].(string)
		if !ok {
			panic("fail")
		}
		i, ok := p[1].(float64)
		if !ok {
			panic("fail")
		}
		return NEMCBlock{s, int(i), nemcRID}
	}

	runtimeIDData := make([][2]interface{}, 0)
	err := json.Unmarshal(blockRunTimeIdTable, &runtimeIDData)
	if err != nil {
		panic(err)
	}
	nemcBlocks := make([]NEMCBlock, 0)
	for rid, jd := range runtimeIDData {
		nemcBlocks = append(nemcBlocks, NewNEMCBlock(jd, rid))
	}
	return nemcBlocks
}

func InitRunTimeIdTable() map[int]string {
	ans := map[int]string{}
	for _, value := range readNemcData() {
		ans[value.NEMCRID] = value.Name
	}
	return ans
}
