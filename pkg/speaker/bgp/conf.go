package bgp

import (
	bgpapi "github.com/openelb/openelb/api/v1alpha2"
	api "github.com/osrg/gobgp/api"
	"golang.org/x/net/context"
)

func (b *Bgp) HandleBgpGlobalConfig(global *bgpapi.BgpConf, rack string, delete bool) error {
	b.rack = rack

	if delete {
		return b.bgpServer.StopBgp(context.Background(), nil)
	}

	request, err := global.Spec.ToGoBgpGlobalConf()
	if err != nil {
		return err
	}

	b.bgpServer.StopBgp(context.Background(), nil)
	//cur:=config.ReadConfigFile(goBGP.toml)
	//config.InitialConfig(cur)
	return b.bgpServer.StartBgp(context.Background(), &api.StartBgpRequest{
		Global: request,
	})
}
