package bgp

import (
	bgpapi "github.com/openelb/openelb/api/v1alpha2"
	api "github.com/osrg/gobgp/api"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"os"
	"text/template"
)

func (b *Bgp) HandleBgpGlobalConfig(global *bgpapi.BgpConf, rack string, delete bool) error {

	/*
		// file /comfigmap/data/goBGPConfig is mounted from configMap
		file, err := os.Open("/comfigmap/data/goBGPConfig")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		goBGPConfig, err := ioutil.ReadAll(file)

	*/

	// The contents of the goBGPConfig file are as follows:
	goBGPConfig := `
[global.config]

#   The golang template is set here to be compatible with BgpConf
    as = {{ .As }}
    router-id = {{ .RouterId }}
    port = {{ .ListenPort }}
    local-address-list = [  {{  range .ListenAddresses }}  {{ . }}, {{ end }}]
   

# User configuration policy

    [global.apply-policy.config]
        import-policy-list = ["policy1"]
        default-import-policy = "reject-route"
        export-policy-list = ["policy2"]
        default-export-policy = "accept-route"

[[defined-sets.prefix-sets]]
    prefix-set-name = "ps0"
    [[defined-sets.prefix-sets.prefix-list]]
        ip-prefix = "10.0.0.0/8"
        masklength-range = "24..32"

.....
....

`
	f, _ := os.Create("goBGP.toml")
	defer f.Close()

	template.Must(template.New("goBGP Configuration").Parse(goBGPConfig)).Execute(f, global.Spec)

	v := viper.New()
	v.SetConfigFile("goBGP.toml")
	v.SetConfigType("toml")

	v.ReadInConfig()

	var request api.Global
	v.UnmarshalExact(&request)

	var policyRequest *api.Policy
	v.UnmarshalExact(&policyRequest)

	b.bgpServer.StartBgp(context.Background(), &api.StartBgpRequest{
		Global: &request,
	})

	b.bgpServer.AddPolicy(context.Background(), &api.AddPolicyRequest{
		Policy: policyRequest,
	})
	return nil
}
