package utils

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GenerateSlug(t *testing.T) {
	testCases := []struct {
		in   string
		want string
	}{
		{"to-my-grandson-\u00a0never-forget-that-i-love-you-fleece-blanket", "to-my-grandson-never-forget-that-i-love-you-fleece-blanket"},
		{"DOBROSLAWZYBORT", "dobroslawzybort"},
		{"Dobroslaw Zybort", "dobroslaw-zybort"},
		{"  Dobroslaw     Zybort  ?", "dobroslaw-zybort"},
		{"  Dobroslaw     Zybort?e", "dobroslaw-zybort-e"},
		{"Dobrosław Żybort", "dobroslaw-zybort"},
		{"Ala ma 6 kotów.", "ala-ma-6-kotow"},

		{"áÁàÀãÃâÂäÄąĄą̊Ą̊", "aaaaaaaaaaaaåå"},
		{"ćĆĉĈçÇ", "cccccc"},
		{"éÉèÈẽẼêÊëËęĘ", "eeeeeeeeeeee"},
		{"íÍìÌĩĨîÎïÏįĮ", "iiiiiiiiiiii"},
		{"łŁ", "ll"},
		{"ńŃ", "nn"},
		{"óÓòÒõÕôÔöÖǫǪǭǬø", "ooooooooooooooo"},
		{"śŚ", "ss"},
		{"úÚùÙũŨûÛüÜųŲ", "uuuuuuuuuuuu"},
		{"y̨Y̨", "y̨y̨"},
		{"źŹżŹ", "zzzz"},

		{"·/,:;`˜'\"", ""},
		{"a·/,:;`˜'\"a", "a-a"},

		{"2000-2013", "2000-2013"},
		{"style—not", "style-not"},
		{"test_slug", "test_slug"},
		{"_test_slug_", "_test_slug_"},
		{"-test-slug-", "test-slug"},
		{"Æ", "ae"},
		{"Ich heiße", "ich-heisse"},

		{"This & that", "this-that"},
		{"fácil €", "facil-eu"},
		{"smile ☺", "smile-☺"},
		{"Hellö Wörld хелло ворлд", "hello-world-хелло-ворлд"},
		{"\"C'est déjà l’été.\"", "cest-deja-l-ete"},
		{"jaja---lol-méméméoo--a", "jaja-lol-mememeoo-a"},
		{"我在", "我在"},
		{"影師", "影師"},
		{"我❤️😁", "我-😁"},
		{"\"·`~！!`@#$¥%^*（）()-=_+【】[]{}\\|、|:;：；，《<,。》.>/?/？～……「」\"", "_"},
		{"1\"2·3`4~5！6!7`8@9#0$1¥2%3^4*5（6）7(8)9-0=1_2+3【4】5[6]7{8}\\9|0、1|2:3;4：5；6，7《8<9,0。1》2.3>4/5?6/7？8～9…0…1「2」3\"4", "12-3-4-5-6-7-8-9-0-1-2-3-4-5-6-789-0-1_2-3-4-567-8-9-0-1-2-3-4-5-6-7-8-9-0-1-2-3-4-5-6-7-8-9-0-1-2-34"},
		{"emoj👿test", "emoj👿test"},
		{"This \n that", "this-that"},
		{"!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~®©_", "0123456789-abcdefghijklmnopqrstuvwxyz-_-abcdefghijklmnopqrstuvwxyz-®©_"},
	}

	fmt.Printf("%s\n%s", GenerateSlug("test"), GenerateSlug("test"))

	Convey("make correct slug", t, func() {
		for _, st := range testCases {
			got := GenerateSlug(st.in)
			So(got, ShouldEqual, st.want)
		}
	})
}
