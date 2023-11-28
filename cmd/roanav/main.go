package main

import (
	"fmt"
	"time"

	"github.com/ebiiim/conbukun/pkg/ao/data"
	"github.com/ebiiim/conbukun/pkg/ao/roanav"
)

func main() {

	p := roanav.NewKrokiPlantUMLPainter(roanav.DefaultKrokiEndpoint, roanav.DefaultKrokiTimeout, data.Maps)

	mdSQV, ok := data.GetMapDataFromName("Qiitun-Si-Vynsom")
	if !ok {
		panic("map data not found")
	}
	mdQQV, ok := data.GetMapDataFromName("Qiient-Qi-Vynsis")
	if !ok {
		panic("map data not found")
	}
	mdQV, ok := data.GetMapDataFromName("Qiitun-Vietis")
	if !ok {
		panic("map data not found")
	}
	mdDE, ok := data.GetMapDataFromName("Dryvein End")
	if !ok {
		panic("map data not found")
	}

	n := &roanav.Navigation{
		Name: "MyGuild#ROA (conbukun@v1.3.0)",
		Portals: []*roanav.Portal{
			roanav.NewPortal(mdSQV.ID, mdQQV.ID, roanav.PortalTypeBlue, time.Now().Add(3*time.Hour), map[string]string{roanav.PortalDataKeyUser: "user1"}),
			roanav.NewPortal(mdSQV.ID, mdQV.ID, roanav.PortalTypeYellow, time.Now().Add(9*time.Hour), map[string]string{roanav.PortalDataKeyUser: "user2"}),
			roanav.NewPortal(mdQQV.ID, mdQV.ID, roanav.PortalTypeBlue, time.Now().Add(-3*time.Hour), map[string]string{roanav.PortalDataKeyUser: "user3"}),
			roanav.NewPortal(mdDE.ID, mdQV.ID, roanav.PortalTypeBlue, time.Now().Add(25*time.Minute), map[string]string{roanav.PortalDataKeyUser: "user4"}),
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
