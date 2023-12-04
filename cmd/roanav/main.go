package main

import (
	"fmt"
	"time"

	"github.com/ebiiim/conbukun/pkg/ao/data"
	"github.com/ebiiim/conbukun/pkg/ao/roanav"
)

func mustGetMD(name string) data.MapData {
	md, ok := data.GetMapDataFromName(name)
	if !ok {
		panic("map data not found")
	}
	return md
}

func mustParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}

func nav1() *roanav.Navigation {
	mdQSV := mustGetMD("Qiitun-Si-Vynsom")
	mdQQV := mustGetMD("Qiient-Qi-Vynsis")
	mdQV := mustGetMD("Qiitun-Vietis")
	mdCA := mustGetMD("Ceritos-Avulsum")
	mdSO := mustGetMD("Suyitos-Ofugtum")

	mdDE := mustGetMD("Dryvein End")      // Black
	mdCM := mustGetMD("Creag Morr")       // Red
	mdSB := mustGetMD("Sleetwater Basin") // Yellow
	mdFF := mustGetMD("Fog Fen")          // Blue

	n := &roanav.Navigation{
		Name: "MyGuild#ROA",
		Portals: []*roanav.Portal{
			roanav.NewPortal(mdQSV.ID, mdQQV.ID, roanav.PortalTypeBlue, time.Now().Add(3*time.Hour), map[string]string{roanav.PortalDataKeyUser: "user1"}),
			roanav.NewPortal(mdQSV.ID, mdQV.ID, roanav.PortalTypeYellow, time.Now().Add(9*time.Hour), map[string]string{roanav.PortalDataKeyUser: "user2"}),
			roanav.NewPortal(mdQQV.ID, mdQV.ID, roanav.PortalTypeBlue, time.Now().Add(-3*time.Hour), map[string]string{roanav.PortalDataKeyUser: "user3"}),
			roanav.NewPortal(mdDE.ID, mdQV.ID, roanav.PortalTypeBlue, time.Now().Add(25*time.Minute), map[string]string{roanav.PortalDataKeyUser: "user4"}),
			roanav.NewPortal(mdSO.ID, mdFF.ID, roanav.PortalTypeYellow, time.Now().Add(3*time.Minute), map[string]string{roanav.PortalDataKeyUser: "user5"}),
			roanav.NewPortal(mdQV.ID, mdCM.ID, roanav.PortalTypeYellow, time.Now().Add(1*time.Minute), map[string]string{roanav.PortalDataKeyUser: "user6"}),
			roanav.NewPortal(mdQV.ID, mdSB.ID, roanav.PortalTypeYellow, time.Now().Add(500*time.Minute), map[string]string{roanav.PortalDataKeyUser: "user7"}),
			roanav.NewPortal(mdCA.ID, mdSO.ID, roanav.PortalTypeYellow, time.Now().Add(100*time.Minute), map[string]string{roanav.PortalDataKeyUser: "user1"}),
		},
		Data: map[string]string{
			roanav.NavigationDataHideouts: fmt.Sprintf("%s,", mdQSV.ID),
		},
	}
	n.DeleteExpiredPortals()
	return n
}

func nav2() *roanav.Navigation {
	usr := "hoge"
	n := &roanav.Navigation{
		Name: "MyGuild#ROA",
		Portals: []*roanav.Portal{
			roanav.NewPortal("TNL-258", "TNL-367", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("7h32m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-301", "TNL-367", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("3h47m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("0218", "TNL-210", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("4h37m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-175", "TNL-344", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("25m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-132", "TNL-344", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("2h26m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-132", "TNL-170", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("7h1m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-096", "TNL-301", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("16m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-007", "TNL-096", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("12h28m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("4220", "TNL-007", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("9h31m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("0333", "TNL-041", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("9h6m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-041", "TNL-374", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("5h45m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("4323", "TNL-041", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("3h48m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-034", "TNL-374", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("9h24m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-034", "TNL-372", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("13h26m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-034", "TNL-038", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("9h56m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("4309", "TNL-034", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("7h30m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-136", "TNL-372", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("6h24m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("3353", "TNL-136", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("24m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("0347", "TNL-136", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("13h43m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-136", "TNL-318", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("2h7m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("2330", "TNL-318", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("10h30m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-281", "TNL-318", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("2h11m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("0321", "TNL-281", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("25m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("4354", "TNL-281", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("2h9m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-268", "TNL-281", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("1h35m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-042", "TNL-258", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("18h7m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("2319", "TNL-042", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("5h37m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("4331", "TNL-042", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("9h14m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("1336", "TNL-108", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("11h2m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-108", "TNL-160", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("3h57m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("1335", "TNL-108", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("20h3m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-160", "TNL-330", roanav.PortalTypeYellow, time.Now().Add(mustParseDuration("5h2m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-330", "TNL-357", roanav.PortalTypeYellow, time.Now().Add(mustParseDuration("7h55m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-152", "TNL-330", roanav.PortalTypeYellow, time.Now().Add(mustParseDuration("1h1m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-280", "TNL-357", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("10h42m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-152", "TNL-248", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("1h3m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("0337", "TNL-248", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("14m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("1324", "TNL-248", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("13h2m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-123", "TNL-248", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("11h34m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("1311", "TNL-248", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("1h14m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("1343", "TNL-280", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("5h48m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-181", "TNL-280", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("37m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("1348", "TNL-280", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("2h59m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-149", "TNL-181", roanav.PortalTypeYellow, time.Now().Add(mustParseDuration("1h42m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-149", "TNL-398", roanav.PortalTypeYellow, time.Now().Add(mustParseDuration("12h27m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-149", "TNL-153", roanav.PortalTypeYellow, time.Now().Add(mustParseDuration("9h12m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-122", "TNL-398", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("3h18m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-122", "TNL-245", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("56m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-122", "TNL-200", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("4h7m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-092", "TNL-200", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("10h33m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-092", "TNL-102", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("16h22m")), map[string]string{roanav.PortalDataKeyUser: usr}),
			roanav.NewPortal("TNL-077", "TNL-092", roanav.PortalTypeBlue, time.Now().Add(mustParseDuration("5h31m")), map[string]string{roanav.PortalDataKeyUser: usr}),
		},
		Data: map[string]string{
			roanav.NavigationDataHideouts: "TNL-367",
		},
	}
	n.DeleteExpiredPortals()
	return n
}

func mustPaintPlantUML(p *roanav.KrokiPlantUMLPNGPainter, n *roanav.Navigation, maps map[string]data.MapData) {
	fmt.Println("===================================")
	fmt.Println(roanav.BriefNavigation(n, maps))
	fmt.Println("===================================")
	s, err := p.ToPlantUML(n)
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	fmt.Println("===================================")
	path, err := p.Paint(n)
	if err != nil {
		panic(err)
	}
	fmt.Println("===================================")
	fmt.Printf("saved: %s\n", path)
	fmt.Println("===================================")
}

func main() {

	p := roanav.NewKrokiPlantUMLPNGPainter(roanav.DefaultKrokiEndpoint, roanav.DefaultKrokiTimeout, data.Maps)

	n1 := nav1()
	mustPaintPlantUML(p, n1, data.Maps)

	n2 := nav2()
	mustPaintPlantUML(p, n2, data.Maps)
}
