package core

import (
	"context"
	"regexp"
	"testing"

	"github.com/tj/assert"
)

func TestFilterer(t *testing.T) {
	testdata := []*Article{
		{Title: "[徵/台北宜蘭] Tamron 70-180 a056 Sony E-mount", Link: "/bbs/DC_SALE/M.1691245105.A.729.html"},
		{Title: "[售/竹苗] SONY FE 50mm F2.8 MACRO", Link: "/bbs/DC_SALE/M.1691245916.A.87B.html"},
		{Title: "[售/ 高屏] 售/高屏 騰龍 Tamron  17-70 f2.8 for S", Link: "/bbs/DC_SALE/M.1691246346.A.3B6.html"},
		{Title: "[售/全國] Nikon D750 + 24-70mm F2.8G", Link: "/bbs/DC_SALE/M.1691246622.A.84D.html"},
		{Title: "[售/台北] DJI Mini 3、Zeiss 12mm F2.8", Link: "/bbs/DC_SALE/M.1691247694.A.47D.html"},
		{Title: "[售/雙北] Sony CEA-G80T CFA 記憶卡 ", Link: "/bbs/DC_SALE/M.1691247785.A.027.html"},
		{Title: "[售/ 高屏] GoPro 11 black", Link: "/bbs/DC_SALE/M.1691249060.A.A35.html"},
		{Title: "[售/台北] Fujifilm XF 14mm F2.8 R 二手 盒裝 ", Link: "/bbs/DC_SALE/M.1691249293.A.E85.html"},
		{Title: "[售/全國]  Sony RX0與相關配件", Link: "/bbs/DC_SALE/M.1691249464.A.96A.html"},
		{Title: "[售/高雄］Gopro11 忠新公司貨保固中", Link: "/bbs/DC_SALE/M.1691254408.A.597.html"},
		{Title: "[售/雙北] Sony A7M3", Link: "/bbs/DC_SALE/M.1691254720.A.861.html"},
		{Title: "[公告] DC_SALE 板規 20160606", Link: "/bbs/DC_SALE/M.1465193037.A.FB9.html"},
		{Title: "[公告] 注意詐騙行為", Link: "/bbs/DC_SALE/M.1581597790.A.F13.html"},
		{Title: "[公告] 禁止於推文內徵求、販售物品", Link: "/bbs/DC_SALE/M.1602469128.A.834.html"},
		{Title: "[黑名] clareCC", Link: "/bbs/DC_SALE/M.1610284022.A.B1E.html"},
		{Title: ">>>Canon R5、控制環轉接環!遭竊!請勿購買<<<", Link: "/bbs/DC_SALE/M.1621574839.A.49C.html"},
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	f := NewFilterer(ctx)
	f.SetInput(make(chan []*Article, 1))
	f.AddInterestTopic(InterestTopic{
		Name: "Sony A7M3",
		CompiledPatterns: []*regexp.Regexp{
			regexp.MustCompile("(?i)sony"),
			regexp.MustCompile("(?i)a7m3"),
		},
	})
	f.AddInterestTopic(InterestTopic{
		Name: "50mm f2.8",
		CompiledPatterns: []*regexp.Regexp{
			regexp.MustCompile("(?i)(sony|sigma)"),
			regexp.MustCompile("(?i)50(mm)?"),
			regexp.MustCompile("(?i)f2.8"),
		},
	})
	f.Run()

	f.Input <- testdata
	matchedArticles := <-f.Output
	assert.Equal(t, 2, len(matchedArticles))
}
