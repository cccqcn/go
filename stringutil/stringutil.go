
package main

import (
	"fmt"
	"strings"
)

var s string = `高峻飞雪高朗妙芙高丽碧凡高邈思柔高旻雁桃高明丹南高爽雁菡高兴翠丝高轩幻梅高雅海莲高扬宛秋高阳问枫高义靖雁高谊蛟凤高逸大凄高懿傻姑高原金连高远梦安高韵碧曼高卓代珊光赫惜珊光华元冬光辉青梦光济书南光霁绮山光亮白桃光临从波光明访冬光启含卉光熙平蝶光耀海秋光誉沛珊光远飞兰国安凝云国兴亦竹国源梦岚冠宇寒凡冠玉傲柔晗昱凌丝晗日觅风涵畅平彤涵涤念露
涵亮翠彤涵忍秋玲涵容安蕾涵润若蕊涵涵灵萱涵煦含雁涵蓄思真涵衍盼山涵意香薇涵映碧萱涵育夏柳翰采白风翰池安双翰飞凌萱翰海盼夏翰翮幻巧翰林怜寒翰墨傲儿翰学冰枫翰音如萱瀚玥妖丽翰藻元芹瀚海涵阳瀚漠涵蕾昊苍以旋昊昊高丽昊空灭男昊乾代玉昊穹可仁昊然可兰昊天可愁昊英可燕浩波妙彤浩博易槐浩初小凝浩大妙晴浩宕冰薇浩荡涵柏浩歌语兰浩广小蕾浩涆忆翠浩瀚听云浩浩觅海浩慨静竹浩旷初蓝浩阔迎丝浩漫幻香浩淼含芙浩渺夏波浩邈冰香浩气凌香浩穰妙菱浩壤访彤浩思凡雁
浩言紫真和蔼书双和安问晴和璧惜萱和昶白萱和畅靖柔和风凡白和歌晓曼和光曼岚和平雁菱和洽雨安和惬谷菱和顺夏烟和硕问儿和颂青亦和泰夏槐和悌含蕊和通迎南和同又琴和煦冷松和雅安雁和宜飞荷和怡踏歌和玉秋莲和裕盼波和豫以蕊和悦盼兰和韵之槐和泽飞柏和正孤容和志白玉弘博傲南弘大山芙弘方夏青弘光雁山弘和曼梅弘厚如霜弘化沛芹弘济丹萱弘阔翠霜弘亮玉兰弘量汝燕弘深不乐弘盛不悔弘图可冥弘伟若男弘新素阴弘雅元彤弘扬从丹弘业曼彤弘义惋庭弘益起眸弘毅香芦弘懿绿竹`


func main() {
	z:=strings.SplitN(s, "\n", len(s))
	fmt.Println(strings.Join(z, "\n\n")) 
	fmt.Println("\n\n")
	sss := ""
	j:=0
	for i:=0; i < len(z); i++{
		ss := z[i]
		for j=0; j < len(ss); j+=6{
			sss = sss + ss[j:j+6]
			sss = sss + ","
	    }
		sss = sss + "\n"
	}
	fmt.Println(sss)
	}