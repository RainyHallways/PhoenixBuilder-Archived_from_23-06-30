# 插件 开

import os,time,json
import jieba
import gensim
import argparse
from threading import Thread
from omega_side.python3_omega_sync.bootstrap import install_lib
install_lib(lib_name="jieba",lib_install_name="jieba")
install_lib(lib_name="gensim",lib_install_name="gensim")
install_lib(lib_name="gensim",lib_install_name="gensim")
from omega_side.python3_omega_sync import API
from omega_side.python3_omega_sync import frame as omega
from omega_side.python3_omega_sync.protocol import *

def main(msg):    
    try:
        def split_word(sentence, stoplist=[]):
            '''分词+删除停用词，返回列表'''
            words = jieba.cut(sentence)
            result = [i for i in words if i not in stoplist]
            return result
        # 解析参数
        parser = argparse.ArgumentParser(description='问答机器人参数')
        parser.add_argument('--data_filepath', default='./omega_python_plugins/data/msg.json')  # 语料路径
        parser.add_argument('--stopwords_filepath', default='./omega_python_plugins/data/stopwords.txt')  # 停用词路径
        parser.add_argument('--splitdata_filepath', default='./omega_python_plugins/data/splitdata.json')  # 分词结果路径
        parser.add_argument('--dictionary_filepath', default='./omega_python_plugins/data/dictionary')  # gensim字典路径
        parser.add_argument('--model_filepath', default='./omega_python_plugins/data/tfidf.model')  # tfidf模型路径
        parser.add_argument('--index_filepath', default='./omega_python_plugins/data/tfidf.index')  # 相似度比较序列路径
        args = parser.parse_args()
        # 语料库
        with open(args.data_filepath, encoding='utf-8') as f:
            data = json.load(f)
        # 停用词
        with open(args.stopwords_filepath, encoding='utf-8') as f:
            stoplist = f.read().splitlines()
        beg = time.time()
        print('分词中...')
        # 加载分词结果，若无则生成
        splitdata_filepath = args.splitdata_filepath
        if os.path.exists(splitdata_filepath):
            with open(splitdata_filepath, encoding='utf-8') as f:
                content = json.load(f)
        else:
            content = []  # 已分词且去停用词的问题
            for key, value in data.items():
                question = value['question']
                content.append(split_word(question, stoplist))
            with open(splitdata_filepath, 'w', encoding='utf-8') as f:
                f.write(json.dumps(content, ensure_ascii=False))
        print('分词耗时 {:.2f}s'.format(time.time() - beg))
        beg = time.time()
        # 加载gensim字典，若无则生成
        dictionary_filepath = args.dictionary_filepath
        if os.path.exists(dictionary_filepath):
            dictionary = gensim.corpora.Dictionary.load(dictionary_filepath)
        else:
            dictionary = gensim.corpora.Dictionary(content)
            dictionary.save(dictionary_filepath)
        print('gensim字典耗时 {:.2f}s'.format(time.time() - beg))
        beg = time.time()
        num_features = len(dictionary)  # 特征数
        # 加载tfidf模型，若无则生成
        model_filepath = args.model_filepath
        if os.path.exists(model_filepath):
            tfidf = gensim.models.TfidfModel.load(model_filepath)
        else:
            corpus = [dictionary.doc2bow(line) for line in content]  # 语料转词袋表示
            tfidf = gensim.models.TfidfModel(corpus)  # 构建tfidf模型
            tfidf.save(args.model_filepath)  # 保存tfidf模型
        # 加载tfidf相似度比较序列，若无则生成
        index_filepath = args.index_filepath
        if os.path.exists(index_filepath):
            index = gensim.similarities.Similarity.load(index_filepath)
        else:
            index = gensim.similarities.Similarity(args.index_filepath, tfidf[corpus], num_features)  # 文本相似度序列
            index.save(index_filepath)
        print('语料转词袋耗时 {:.2f}s'.format(time.time() - beg))
        sentences = msg[3:-1]
        sentences = split_word(sentences, stoplist)  # 分词
        vec = dictionary.doc2bow(sentences)  # 转词袋表示
        sims = index[tfidf[vec]]  # 相似度比较
        sorted_sims = sorted(enumerate(sims), key=lambda x: x[1], reverse=True) 
        print('分词结果 ->  ', sentences)
        print("相似度比较 ->  ", sorted_sims[:5])
        for similarity in sorted_sims[:1]:
            print(dict(data)[str(similarity[0])])
        return dict(data)[str(similarity[0])]
    except:
        err = {'question': sentences, 'answer': '目前无法为您解答该问题!'}
        return dict(err)

def plugin_main(api:API):
    def on_menu_main(packet):
        msg = packet["Message"]
        player = packet['SourceName']
        msg = msg+" "
        msg = main(str(msg))
        msg = "问题:"+str(dict(msg)["question"])+"\n回答:"+str(dict(msg)['answer'])
        api.do_send_player_msg(player,msg,cb=None)
    def getIDAddPlayer_IDTextr_pakcet(packet):
        if "问答" in packet["Message"]:
            Thread(target=on_menu_main,args=[packet]).start()
    api.listen_mc_packet(pkt_type="IDText",cb=None,on_new_packet_cb=getIDAddPlayer_IDTextr_pakcet)

omega.add_plugin(plugin=plugin_main)