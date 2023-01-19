package Happy2018new

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"phoenixbuilder/minecraft/protocol/packet"
	"phoenixbuilder/omega/defines"
	"strings"
	"time"
)

type RecordBlockChanges struct {
	*defines.BasicComponent
	MaxPlayerRecord        int     `json:"每次至多追踪的玩家数"`
	IsOutputJsonDatas      bool    `json:"启动本组件时统计数据并输出 JSON 日志"`
	DiscardUnknwonOperator bool    `json:"丢弃未知操作来源的方块"`
	TrackingRadius         float64 `json:"追踪半径"`
	FileName               string  `json:"文件名称"`
	DataReceived           []struct {
		Time             string
		BlockPos         [3]int32
		BlockName_Result string
		Situation        uint32
		Operator         []string
	}
}

func (o *RecordBlockChanges) Init(settings *defines.ComponentConfig) {
	marshal, _ := json.Marshal(settings.Configs)
	if err := json.Unmarshal(marshal, o); err != nil {
		panic(err)
	}
	if o.MaxPlayerRecord <= 0 {
		o.MaxPlayerRecord = 1
	}
}

func (o *RecordBlockChanges) Inject(frame defines.MainFrame) {
	o.Frame = frame
	o.FileName = "RecordBlockChanges.Happy2018new"
}

func (o *RecordBlockChanges) RequestBlockChangesInfo(BlockInfo packet.UpdateBlock) {
	defer func() {
		err := recover()
		if err != nil {
			o.RequestBlockChangesInfo(BlockInfo)
		}
	}()
	var blockName_Result string = "air"
	var resp packet.CommandOutput
	var operator []string = []string{}
	o.Frame.GetGameControl().SendCmdAndInvokeOnResponse(
		fmt.Sprintf("testforblock %v %v %v air", BlockInfo.Position.X(), BlockInfo.Position.Y(), BlockInfo.Position.Z()),
		func(output *packet.CommandOutput) {
			resp = *output
			if resp.SuccessCount <= 0 {
				blockName_Result = resp.OutputMessages[0].Parameters[3]
				blockName_Result = strings.Replace(blockName_Result, fmt.Sprintf("%vtile.", "%"), "", 1)
				blockName_Result = strings.Replace(blockName_Result, ".name", "", 1)
			}
			resp = packet.CommandOutput{}
			o.Frame.GetGameControl().SendCmdAndInvokeOnResponse(
				fmt.Sprintf("testfor @a[c=%v,x=%v,y=%v,z=%v,r=%v]", o.MaxPlayerRecord, BlockInfo.Position.X(), BlockInfo.Position.Y(), BlockInfo.Position.Z(), o.TrackingRadius),
				func(output *packet.CommandOutput) {
					resp = *output
					if resp.SuccessCount > 0 {
						operator = strings.Split(resp.OutputMessages[0].Parameters[0], ", ")
					} else {
						operator = []string{"unknown"}
					}
					if o.DiscardUnknwonOperator && resp.SuccessCount > 0 {
						o.DataReceived = append(o.DataReceived, struct {
							Time             string
							BlockPos         [3]int32
							BlockName_Result string
							Situation        uint32
							Operator         []string
						}{
							Time:             time.Now().Format("2006-01-02 15:04:05"),
							BlockPos:         BlockInfo.Position,
							BlockName_Result: blockName_Result,
							Situation:        BlockInfo.Flags,
							Operator:         operator,
						})
					}
					if !o.DiscardUnknwonOperator {
						o.DataReceived = append(o.DataReceived, struct {
							Time             string
							BlockPos         [3]int32
							BlockName_Result string
							Situation        uint32
							Operator         []string
						}{
							Time:             time.Now().Format("2006-01-02 15:04:05"),
							BlockPos:         BlockInfo.Position,
							BlockName_Result: blockName_Result,
							Situation:        BlockInfo.Flags,
							Operator:         operator,
						})
					}
				},
			)
		},
	)
}

func (o *RecordBlockChanges) OutputDatas() []byte {
	ans := []byte{}
	// prepare
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, int32(len(o.DataReceived)))
	ans = append(ans, buf.Bytes()...)
	// data length
	for _, value := range o.DataReceived {
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, int32(len([]byte(value.Time))))
		ans = append(ans, buf.Bytes()...)
		ans = append(ans, []byte(value.Time)...)
		// time
		for _, val := range value.BlockPos {
			buf := bytes.NewBuffer([]byte{})
			binary.Write(buf, binary.BigEndian, val)
			ans = append(ans, buf.Bytes()...)
		}
		// pos
		buf = bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, int32(len([]byte(value.BlockName_Result))))
		ans = append(ans, buf.Bytes()...)
		ans = append(ans, []byte(value.BlockName_Result)...)
		// blockName_Result
		buf = bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, value.Situation)
		ans = append(ans, buf.Bytes()...)
		// situation
		buf = bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, int32(len(value.Operator)))
		ans = append(ans, buf.Bytes()...)
		for _, val := range value.Operator {
			buf = bytes.NewBuffer([]byte{})
			binary.Write(buf, binary.BigEndian, int32(len([]byte(val))))
			ans = append(ans, buf.Bytes()...)
			ans = append(ans, []byte(val)...)
		}
		// operator
	}
	return ans
}

func (o *RecordBlockChanges) GetDatas() {
	ans := []struct {
		Time             string
		BlockPos         [3]int32
		BlockName_Result string
		Situation        uint32
		Operator         []string
	}{}
	got, err := o.Frame.GetFileData(o.FileName)
	if len(got) <= 0 || err != nil {
		o.DataReceived = ans
		return
	}
	// prepare
	reader := bytes.NewReader(got)
	p := make([]byte, 4)
	n, err := reader.Read(p)
	if n < 4 || err != nil {
		panic("无法读取保存的文件，请检查您的文件是否已经损坏！")
	}
	// get length
	buf := bytes.NewBuffer(p)
	var length int32
	binary.Read(buf, binary.BigEndian, &length)
	// decode length
	for i := 0; i < int(length); i++ {
		p = make([]byte, 4)
		n, err = reader.Read(p)
		if n < 4 || err != nil {
			panic("无法读取保存的文件，请检查您的文件是否已经损坏！")
		}
		// get length of time
		buf = bytes.NewBuffer(p)
		var timeLength int32
		binary.Read(buf, binary.BigEndian, &timeLength)
		// decode length of time
		p = make([]byte, timeLength)
		n, err = reader.Read(p)
		if n < int(timeLength) || err != nil {
			panic("无法读取保存的文件，请检查您的文件是否已经损坏！")
		}
		time := string(p)
		// time
		pos := [3]int32{}
		for j := 0; j < 3; j++ {
			p = make([]byte, 4)
			n, err = reader.Read(p)
			if n < 4 || err != nil {
				panic("无法读取保存的文件，请检查您的文件是否已经损坏！")
			}
			// get pos[j]
			buf = bytes.NewBuffer(p)
			var posSingle int32
			binary.Read(buf, binary.BigEndian, &posSingle)
			// decode pos[j]
			pos[j] = posSingle
		}
		// blockPos
		p = make([]byte, 4)
		n, err = reader.Read(p)
		if n < 4 || err != nil {
			panic("无法读取保存的文件，请检查您的文件是否已经损坏！")
		}
		// get length of blockName_Result
		buf = bytes.NewBuffer(p)
		var blockName_Result_length int32
		binary.Read(buf, binary.BigEndian, &blockName_Result_length)
		// decode length of blockName_Result
		p = make([]byte, blockName_Result_length)
		n, err = reader.Read(p)
		if n < int(blockName_Result_length) || err != nil {
			panic("无法读取保存的文件，请检查您的文件是否已经损坏！")
		}
		blockName_Result := string(p)
		// blockName_Result
		p = make([]byte, 4)
		n, err = reader.Read(p)
		if n < 4 || err != nil {
			panic("无法读取保存的文件，请检查您的文件是否已经损坏！")
		}
		// get situation
		buf = bytes.NewBuffer(p)
		var situation uint32
		binary.Read(buf, binary.BigEndian, &situation)
		// decode situation
		p = make([]byte, 4)
		n, err = reader.Read(p)
		if n < 4 || err != nil {
			panic("无法读取保存的文件，请检查您的文件是否已经损坏！")
		}
		// get length of operator
		buf = bytes.NewBuffer(p)
		var operatorLength int32
		binary.Read(buf, binary.BigEndian, &operatorLength)
		// decode length of operator
		operator := []string{}
		for j := 0; j < int(operatorLength); j++ {
			p = make([]byte, 4)
			n, err = reader.Read(p)
			if n < 4 || err != nil {
				panic("无法读取保存的文件，请检查您的文件是否已经损坏！")
			}
			// get length of operator(single)
			buf = bytes.NewBuffer(p)
			var operatorSingleLength int32
			binary.Read(buf, binary.BigEndian, &operatorSingleLength)
			// decode length of operator(single)
			p = make([]byte, operatorSingleLength)
			n, err = reader.Read(p)
			if n < int(operatorSingleLength) || err != nil {
				panic("无法读取保存的文件，请检查您的文件是否已经损坏！")
			}
			operator = append(operator, string(p))
			// operator(single)
		}
		ans = append(ans, struct {
			Time             string
			BlockPos         [3]int32
			BlockName_Result string
			Situation        uint32
			Operator         []string
		}{
			Time:             time,
			BlockPos:         pos,
			BlockName_Result: blockName_Result,
			Situation:        situation,
			Operator:         operator,
		})
	}
	o.DataReceived = ans
}

func (o *RecordBlockChanges) StatisticsDatas() {
	type blockCube struct {
		Posx int32
		Posy int32
		Posz int32
	}
	type single struct {
		Time             string
		BlockName_Result string
		Situation        uint32
		Operator         []string
	}
	type set []single
	blockCubeMap := map[blockCube]set{}
	for _, value := range o.DataReceived {
		got, ok := blockCubeMap[blockCube{value.BlockPos[0], value.BlockPos[1], value.BlockPos[2]}]
		if !ok {
			blockCubeMap[blockCube{value.BlockPos[0], value.BlockPos[1], value.BlockPos[2]}] = set{
				single{
					Time:             value.Time,
					BlockName_Result: value.BlockName_Result,
					Situation:        value.Situation,
					Operator:         value.Operator,
				},
			}
		} else {
			got = append(got, single{
				Time:             value.Time,
				BlockName_Result: value.BlockName_Result,
				Situation:        value.Situation,
				Operator:         value.Operator,
			})
			blockCubeMap[blockCube{value.BlockPos[0], value.BlockPos[1], value.BlockPos[2]}] = got
		}
	}
	new := map[string]interface{}{}
	for key, value := range blockCubeMap {
		singleNew := []interface{}{}
		for _, val := range value {
			operatorNew := []interface{}{}
			for _, v := range val.Operator {
				operatorNew = append(operatorNew, v)
			}
			singleNew = append(singleNew, map[string]interface{}{
				"操作时间":   val.Time,
				"关联的方块名": val.BlockName_Result,
				"附加数据":   float64(val.Situation),
				"可能的操作者": operatorNew,
			})
		}
		new[fmt.Sprintf("方块 (%v,%v,%v)", key.Posx, key.Posy, key.Posz)] = singleNew
	}
	o.Frame.WriteJsonData(fmt.Sprintf("%v.json", o.FileName), new)
}

func (o *RecordBlockChanges) Activate() {
	o.GetDatas()
	if o.IsOutputJsonDatas {
		o.StatisticsDatas()
	}
	o.Frame.GetGameListener().SetOnTypedPacketCallBack(packet.IDUpdateBlock, func(p packet.Packet) {
		go func() {
			o.RequestBlockChangesInfo(*p.(*packet.UpdateBlock))
		}()
	})
}

func (o *RecordBlockChanges) Signal(signal int) error {
	switch signal {
	case defines.SIGNAL_DATA_CHECKPOINT:
		return o.Frame.WriteFileData(o.FileName, o.OutputDatas())
	}
	return nil
}

func (o *RecordBlockChanges) Stop() error {
	fmt.Println("正在保存 " + o.FileName)
	return o.Frame.WriteFileData(o.FileName, o.OutputDatas())
}
