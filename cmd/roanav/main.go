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

func main() {

	p := roanav.NewKrokiPlantUMLPainter(roanav.DefaultKrokiEndpoint, roanav.DefaultKrokiTimeout, data.Maps)

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

	fmt.Println(roanav.BriefNavigation(n, data.Maps))

	fmt.Print("\n\n\n")

	s, err := p.ToPlantUML(n)
	if err != nil {
		panic(err)
	}
	fmt.Println(s)

	fmt.Print("\n\n\n")

	path, err := p.Paint(n)
	if err != nil {
		panic(err)
	}
	fmt.Println(path)
}
