# Name：tools
# Version：0.0.1
# Author：Mono
# 请不要随意更改
import json
import threading

from omega_side.python3_omega_sync import API


class api(object):
    def __init__(self, API: API) -> None:
        self.api = API
        self.json_lock = threading.Lock()

    def execute_after(self, func, *args, delay_time):
        return self.api.execute_after(func, *args, delay_time=delay_time)

    def execute_with_repeat(self, func, *args, repeat_time):
        return self.api.execute_with_repeat(func, *args, repeat_time=repeat_time)

    def execute_in_individual_thread(self, func, *args):
        return self.api.execute_in_individual_thread(func, *args)

    def do_get_get_player_next_param_input(self, player, hint="请随便说点什么", cb=None):
        return self.api.do_get_get_player_next_param_input(player=player, hint=hint, cb=cb)

    def do_send_ws_cmd(self, cmd, cb=None):
        return self.api.do_send_ws_cmd(cmd=cmd, cb=cb)

    def do_send_player_cmd(self, cmd, cb=None):
        return self.api.do_send_player_cmd(cmd=cmd, cb=cb)

    def do_send_wo_cmd(self, cmd, cb=None):
        return self.api.do_send_wo_cmd(cmd=cmd, cb=cb)

    def do_get_players_list(self, cb=None):
        return self.api.do_get_players_list(cb=cb)

    def sayTo(self, player: str, msg: str, cb=None):
        """
        给玩家发生消息
        ---
        player : 玩家名
        msg : 消息
        """
        self.api.do_send_player_msg(player=player, msg=msg, cb=cb)

    def echo(self, msg, cb=None):
        return self.api.do_echo(msg=msg, cb=cb)

    def sendwscmd(self, cmd, cb=None):
        return self.api.do_send_ws_cmd(cmd=cmd, cb=cb)

    def sendplayercmd(self, cmd, cb=None):
        return self.api.do_send_player_cmd(cmd=cmd, cb=cb)

    def sendwocmd(self, cmd, cb=None):
        return self.api.do_send_wo_cmd(cmd=cmd, cb=cb)

    def get_uqholder(self, cb=None):
        return self.api.do_get_uqholder(cb=cb)

    def get_players_list(self, cb=None):
        return self.api.do_get_players_list(cb=cb)

    def set_title(self, player, msg, cb=None):
        return self.api.do_set_player_title(player=player, msg=msg, cb=cb)

    def set_subtitle(self, player, msg, cb=None):
        return self.api.do_set_player_subtitle(player=player, msg=msg, cb=cb)

    def set_actionbar(self, player, msg, cb=None):
        return self.api.do_set_player_actionbar(player=player, msg=msg, cb=cb)

    def getPos(self, player, limit, cb=None):
        """
        limit实例:@p[name=[player]]
        返回:dict
        response.pos -> [1, 117, -2]
        """
        return self.api.do_get_player_pos(player=player, limit=limit, cb=cb)

    def set_player_data(self, player, entry, data, cb=None):
        return self.api.do_set_player_data(player=player, entry=entry, data=data, cb=cb)

    def get_player_data(self, player, entry, cb=None):
        return self.api.do_get_player_data(player=player, entry=entry, cb=cb)

    def do_get_item_mapping(self, cb=None):
        return self.api.do_get_item_mapping(cb=cb)

    def do_get_block_mapping(self, cb=None):
        return self.api.do_get_block_mapping(cb=cb)

    def do_get_scoreboard(self, cb=None):
        return self.api.do_get_scoreboard(cb=cb)

    def do_send_fb_cmd(self, cmd, cb=None):
        return self.api.do_send_fb_cmd(cmd=cmd, cb=cb)

    def do_send_qq_msg(self, msg, cb=None):
        return self.api.do_send_qq_msg(msg=msg, cb=cb)

    def listen_omega_menu(self, triggers=[], argument_hint="", usage="", on_menu_invoked=None, cb=None):
        return self.api.listen_omega_menu(triggers=triggers, argument_hint=argument_hint, usage=usage,
                                          on_menu_invoked=on_menu_invoked, cb=cb)

    def listen_mc_packet(self, pkt_type, on_new_packet_cb, cb=None):
        return self.api.listen_mc_packet(pkt_type=pkt_type, cb=cb, on_new_packet_cb=on_new_packet_cb)

    def listen_any_mc_packet(self, on_new_packet_cb, cb=None):
        return self.api.listen_any_mc_packet(cb=cb, on_new_packet_cb=on_new_packet_cb)

    def listen_player_login(self, on_player_login_cb, cb=None):
        return self.api.listen_player_login(cb=cb, on_player_login_cb=on_player_login_cb)

    def listen_player_logout(self, on_player_logout_cb, cb=None):
        return self.api.listen_player_logout(cb=cb, on_player_logout_cb=on_player_logout_cb)

    def listen_block_update(self, on_block_update, cb=None):
        return self.api.listen_block_update(cb=cb, on_block_update=on_block_update)

    # 写入数据到JSON文件
    def write_json_file(self, filename, dict):
        try:
            self.json_lock.acquire()
            with open(f"./data/{filename}.json", 'w', encoding='utf-8') as file:
                json.dump(dict, file, indent=4, ensure_ascii=False)
        except Exception:
            return False
        finally:
            self.json_lock.release()
        return True

    # 从JSON文件读取数据
    def read_json_file(self, filename):
        try:
            self.json_lock.acquire()
            with open(f"./data/{filename}.json", 'r', encoding='utf-8') as file:
                result = json.load(file)
        except Exception:
            result = {}
        finally:
            self.json_lock.release()
        return result

    # 写入数据到JSON文件

    def write_player_json_file(self, player, filename, dict):
        try:
            self.json_lock.acquire()
            with open(f"./data/player/{player}/{filename}.json", 'w', encoding='utf-8') as file:
                json.dump(dict, file, indent=4, ensure_ascii=False)
        except Exception:
            return False
        finally:
            self.json_lock.release()
        return True

    # 从JSON文件读取数据
    def read_player_json_file(self, player, filename):
        try:
            self.json_lock.acquire()
            with open(f"./data/player/{player}/{filename}.json", 'r', encoding='utf-8') as file:
                result = json.load(file)
        except Exception:
            result = {}
        finally:
            self.json_lock.release()
        return result

    # 通过玩家UUID查询名字
    def get_player_name(self, UUID):
        for player in self.do_get_players_list():
            if player.uuid == UUID:
                return player.name
        return None

    # 通过玩家名查询UUID
    def get_player_uuid(self, name):
        for player in self.do_get_players_list():
            if player.name == name:
                return player.uuid
        return None

    # 根据玩家名查询权限等级
    def get_player_permission(self, name):
        # 获取uqholder
        uqholder = self.api.do_get_uqholder(cb=None)
        # 权限列表
        OPPermissionLevelList = ["访客", "成员", "操作员"]
        # 解析
        for data in uqholder.PlayersByEntityID.values():
            if data.Username == name:
                try:
                    return OPPermissionLevelList[data.OPPermissionLevel]
                except Exception:
                    return f'unknow<{data.OPPermissionLevel}>'
        return None

    # 根据玩家名查询维度与位置信息
    def get_player_pos(self, name):
        # 发送指令
        response = self.do_send_ws_cmd(f"querytarget @a[name=\"{name}\"]")
        # 没有目标则不处理
        if not response.result.OutputMessages[0].Success:
            return None
        # 解析
        for data in json.loads(response.result.OutputMessages[0].Parameters[0]):
            return {"d": data['dimension'], "x": int(data['position']['x']), "y": int(data['position']['y']),
                    "z": int(data['position']['z'])}
        return None

    # 根据玩家名查询yRot
    def get_player_yRot(self, name):
        # 发送指令
        response = self.do_send_ws_cmd(f"querytarget @a[name=\"{name}\"]")
        # 没有目标则不处理
        if not response.result.OutputMessages[0].Success:
            return None
        # 解析
        for data in json.loads(response.result.OutputMessages[0].Parameters[0]):
            return data['yRot']
        return None

    # 获取玩家分数
    def get_player_score(self, name, scoreboard):
        # 发送指令
        response = self.do_send_ws_cmd(f"scoreboard players add @a[name=\"{name}\"] {scoreboard} 0")
        # 解析
        if response.result.OutputMessages[0].Success:
            return response.result.OutputMessages[0].Parameters[3]
        return None

    # 设置玩家分数
    def set_player_score(self, name, scoreboard, score):
        response = self.do_send_ws_cmd(f"scoreboard players set @a[name=\"{name}\"] {scoreboard} {score}")
        # 解析
        if response.result.OutputMessages[0].Success:
            return True
        return False

    # 增加玩家分数
    def add_player_score(self, name, scoreboard, score):
        response = self.do_send_ws_cmd(f"scoreboard players add @a[name=\"{name}\"] {scoreboard} {score}")
        # 解析
        if response.result.OutputMessages[0].Success:
            return True
        return False

    # 扣除玩家分数 - 扣除后小于0则失败
    def remove_player_score(self, name, scoreboard, score):
        response = self.do_send_ws_cmd(
            f"scoreboard players remove @a[name=\"{name}\",scores={{{scoreboard}={score}..}}] {scoreboard} {score}")
        # 解析
        if response.result.OutputMessages[0].Success:
            return True
        return False

    # 向所有玩家发送一条消息
    def send_all_player_msg(self, msg):
        self.execute_after(func=lambda: self.do_send_wo_cmd(f"tellraw @a {{\"rawtext\":[{{\"text\":\"{msg}\"}}]}}"),
                           delay_time=0.1)
        return True

    # 获取两个坐标之间的距离
    def get_distance(self, posx1, posy1, posz1, posx2, posy2, posz2):
        return pow(pow(posx2 - posx1, 2) + pow(posy2 - posy1, 2) + pow(posz2 - posz1, 2), 0.5)

    def searchName(self, key, data):
        """
        模糊查找器
        key: 关键字(str)
        data: 数据(list)
        :return: list
        """
        return [i for i in data if key in i]

    def getWord(self) -> dict:
        import requests
        url = "https://v1.hitokoto.cn"
        response = requests.get(url)
        data = response.json()
        list_values = [i for i in data.values()]
        if "hitokoto" in data:
            return {"main": list_values[2], "from": list_values[4], "author": list_values[5], "succ": True}
        return {"succ": False}
