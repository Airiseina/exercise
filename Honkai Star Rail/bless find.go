package main

import "fmt"

func main() {
	fmt.Println("请输入你想得到的黄金裔的祝福")
	var input string
	xilian := &PhLia093{}
	fmt.Scan(&input)
	other := &huangjingyi{name: input}
	if input == "昔涟" {
		print(xilian)
	} else {
		print(other)
	}
}

type huangjingyi struct {
	code    string
	name    string
	talking string
}
type PhLia093 struct {
	name string
	talk string
}
type allhuangjingyi interface {
	bless()
}

func (name *PhLia093) bless() {
	name.name = "昔涟"
	name.talk = "愿[开拓]的结局如我们所书"
	fmt.Println(*name)
}
func (name *huangjingyi) bless() {
	var huangjingyiList = []huangjingyi{
		{code: "黄金之茧", name: "阿格莱雅", talking: "愿[浪漫]相伴你的前程"},
		{code: "万径之门", name: "缇里西庇俄丝", talking: "[门径]为你指明前路"},
		{code: "天谴之矛", name: "迈德漠斯", talking: "让[纷争]给予你鼓舞"},
		{code: "灰黯之手", name: "遐蝶", talking: "[死亡]呵护你的灵魂"},
		{code: "裂分之枝", name: "阿那克萨戈拉斯", talking: "由[理性]赐你以启豪 "},
		{code: "晨昏之眼", name: "雅辛忒丝", talking: "[天空]为你洒落晨曦"},
		{code: "翻飞之币", name: "赛法利娅", talking: "[诡计]保你万无一失"},
		{code: "▇▇▇▇▇▇", name: "卡厄斯兰那", talking: "[负世]定永志不忘"},
		{code: "满溢之杯", name: "海列屈拉", talking: "令[海洋]为你起舞"},
		{code: "公正之秤", name: "刻律德莉", talking: "由[律法]踏碎不公"},
		{code: "永夜之帷", name: "三月七", talking: "[岁月]会铭记旅途"},
		{code: "磐若之脊", name: "丹恒", talking: "[大地]护卫你我前行"},
	}
	for _, people := range huangjingyiList {
		if (*name).name == people.name {
			fmt.Println(people)
			return
		}
	}
	fmt.Println("请输入黄金裔的名字")
	return

}
func print(name allhuangjingyi) {
	name.bless()
}
