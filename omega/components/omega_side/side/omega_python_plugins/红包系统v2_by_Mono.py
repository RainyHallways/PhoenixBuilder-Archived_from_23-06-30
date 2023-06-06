# 插件: 关
# 需要使用的请把这个"关"改为"开"

# 此文件请勿随意更改

import os
import random
import uuid

from omega_side.python3_omega_sync import API
from omega_side.python3_omega_sync import frame as omega
from omega_side.python3_omega_sync.protocol import *
from tools_by_mono import api


class redBag(object):
    def version(self):
        return "2.1"

    def algorithm(self, total_amount: int, total_people: int):
        """
        红包金额计算
        total_amount -> 红包金额
        total_people: -> 红包数量

        返回dict
        """
        if total_amount < 0 or total_amount == 0:
            return {"count": 0, "error": "红包金额小于或者等于0"}
        if total_people < 0 or total_people == 0:
            return {"count": 0, "error": "红包数量小于或者等于0"}
        if total_amount < total_people:
            return {"count": 0, "error": "红包金额不能小于红包数量"}
        amount_list = []
        rest_amount = total_amount
        rest_people = total_people
        for i in range(0, total_people - 1):
            amount = random.randint(1, int(rest_amount / rest_people * 2) - 1)
            rest_amount = rest_amount - amount
            rest_people -= 1
            amount_list.append(amount)
        amount_list.append(rest_amount)
        return {"count": 1, "reduce": amount_list}

    def send(self, player: str, coin_num: int, num: int, leavingAMessage: str = None):
        if coin_num <= 0:
            return {"count": 0, "error": "红包金额小于或者等于0"}
        if num <= 0:
            return {"count": 0, "error": "红包数量小于或者等于0"}
        if not os.path.isdir(f"data/player/{player}"):
            os.mkdir(f"data/player/{player}")
        oneData = self.api.read_player_json_file(player, "红包数据")
        First = False
        if oneData == {}:
            First = True
        calculationReturnResult = self.algorithm(coin_num, num)
        if calculationReturnResult == 0:
            self.api.sayTo(player, "错误:" + calculationReturnResult["error"])
            return {"count": 0, "error": calculationReturnResult["error"]}
        collectionOrder = calculationReturnResult['reduce']
        writeData = {"玩家名": player, "红包金额": coin_num, "红包数量": num, "附加的消息": leavingAMessage,
                     "领取顺序": collectionOrder, "已领玩家": []}
        if First is True:
            num_ = 1
            oneData[f"{num_}"] = writeData
            self.api.write_player_json_file(player, "红包数据", oneData)
        else:
            num_ = len(oneData) + 1
            oneData[f"{num_}"] = writeData
            self.api.write_player_json_file(player, "红包数据", oneData)
        RedBagData = self.api.read_json_file("红包系统数据")
        UUID = str(uuid.uuid4()).split("-")[0]
        RedBagData["信息"][UUID] = [f"{player}", num_]
        self.api.write_json_file("红包系统数据", RedBagData)
        self.api.sayTo("@a",
                       f"§6{player} §b| §f发送了§c红包§e{coin_num}r 0/{num} \n祝福:§e{leavingAMessage}§f\n§f红包代号:{UUID}")
        print(player, "发送了总金额", coin_num, "r,总数为", num, "的红包,祝福:", leavingAMessage)
        return {"count": 1, "code": UUID}

    def receive(self, UUID: str, player: str):
        RedBagData = self.api.read_json_file("红包系统数据")
        if UUID not in RedBagData['信息']:
            return {"count": 0, "error": "未找到红包"}
        Sender, Num = RedBagData['信息'][UUID][0], RedBagData['信息'][UUID][1]
        oneData = self.api.read_player_json_file(Sender, "红包数据")
        if oneData[f"{Num}"]["红包数量"] == 0:  # 红包金额 红包数量 顺序
            return {"count": 0, "error": f"你的手速太慢了,§c红包§r已经领完了"}
        elif player in oneData[f'{Num}']["已领玩家"]:
            return {"count": 0, "error": f"你的已经领取过了"}
        RemCoin = oneData[f'{Num}']["领取顺序"][0]
        oneData[f'{Num}']["红包金额"] -= RemCoin
        oneData[f'{Num}']["红包数量"] -= 1
        del (oneData[f'{Num}']["领取顺序"][0])
        oneData[f'{Num}']["已领玩家"].append(player)
        if oneData[f"{Num}"]["红包数量"] == 0:
            del RedBagData['信息'][UUID]
            del oneData[f"{Num}"]
            self.api.write_json_file("红包系统数据", RedBagData)
        self.api.write_player_json_file(Sender, "红包数据", oneData)
        self.api.sayTo(player, f"§6{Sender} §b| §f成功领取§e{RemCoin}r")
        self.api.add_player_score(player, self.scoreboardName, RemCoin)

        return {"count": 1}

    def back(self, player: str):
        try:
            RedBagData = self.api.read_json_file("红包系统数据")
            oneData = self.api.read_player_json_file(player, "红包数据")
            backData_UUID_Num = []
            will_del = []
            for i in RedBagData["信息"]:
                if player in RedBagData["信息"][i]:
                    backData_UUID_Num.append(RedBagData["信息"][i][1])
                    will_del.append(i)
            for i in will_del:
                del RedBagData["信息"][i]
            back_num = 0
            for i in backData_UUID_Num:
                back_num += oneData[f"{i}"]["红包金额"]
                del oneData[f"{i}"]
            self.api.write_json_file("红包系统数据", RedBagData)
            self.api.write_player_json_file(player, "红包数据", oneData)
            self.api.add_player_score(player, self.scoreboardName, back_num)
            self.api.sayTo(player, f"§b| §a成功退回§e{back_num}§cr§f.")
            print(player, "退回了", back_num, "积分")
            return {"count": 1}
        except Exception as e:
            return {"count": 0, "error": e}

    def menu(self, player_input: PlayerInput):
        player = player_input.Name
        msg = player_input.Msg
        if len(msg) == 0:
            self.api.sayTo(player,
                           f"§b| §a红包 Menu§6v{self.version()} §r§a页数§e1|1§r\n§b| §f发 <金额> <数量> <祝福语> §7发送红包\n§b| §f领/抢 <代号> §7领取红包\n§b| §f退 §7退回还未领完的红包.")
        elif msg[0] in ["s", "send", "发", "发红包"]:
            if len(msg) in [1, 2]:
                self.api.sayTo(player, "§b| §c参数不全")
                return
            coin_num, num = int(msg[1]), int(msg[2])
            leavingAMessage = self.leavingAMessage
            if len(msg) == 4:
                leavingAMessage = msg[3]
            response = self.api.do_get_scoreboard()
            if player in response[self.scoreboardName]:
                if response[self.scoreboardName][player] < int(msg[1]):  # 如果玩家输入的红包金额大于自己的金额的话
                    self.api.sayTo(player, "§b| §c金额不足§r")
                    return
            else:
                self.api.sayTo(player, "§b| §c金额不足§r")
                return
            result = self.send(player, coin_num, num, leavingAMessage)
            if result["count"] == 0:
                return
        elif msg[0] in ["l", "领", "抢", "q"]:
            if len(msg) == 1:
                self.api.sayTo(player, "§b| §c参数不全")
                return
            UUID = msg[1]
            self.receive(UUID, player)
        elif msg[0] in ["back", "b", "退", "退款"]:
            result = self.back(player)
            if result["count"] == 0:
                self.api.sayTo(player, f'§b| §cERROR:{result["error"]}')

    def __call__(self, API: API):
        # Start
        self.api = api(API)
        if not os.path.exists(os.path.join('data', '红包系统数据.json')):
            self.api.write_json_file("红包系统数据", {"名称": "红包系统数据存储",
                                                      "描述": "储存谁发了红包和发了多少红包(这不是配置文件哦)",
                                                      "信息": {}})
        if not os.path.isdir(os.path.join('data', 'player')):
            os.mkdir(os.path.join('data', 'player'))
        if not os.path.exists('组件_红包系统.json'):
            self.api.write_json_file('组件_红包系统.json',
                                     {"名称": "红包系统", "描述": "红包系统配置文件,跨服红包暂不支持",
                                      "配置": {"计分板名称": "存储", "默认祝福语": "恭喜发财", "禁用跨服红包": True,
                                               "禁用跨服红包发来的指令": True}})
        data = self.api.read_json_file('组件_红包系统.json')
        self.scoreboardName = data["配置"]["计分板名称"]
        self.leavingAMessage = data["配置"]["默认祝福语"]
        self.api.listen_omega_menu(triggers=["红包", "hb"], argument_hint="", usage="获取红包-菜单", cb=None,
                                   on_menu_invoked=self.menu)
        print("红包系统启动成功.")


if __name__ == "__main__":
    omega.add_plugin(plugin=redBag())
    omega.run(addr=None)
else:
    omega.add_plugin(plugin=redBag())