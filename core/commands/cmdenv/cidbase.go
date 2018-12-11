package cmdenv

import (
	path "gx/ipfs/QmQtg7N4XjAk2ZYpBjjv8B6gQprsRekabHBCnF6i46JYKJ/go-path"
	cmds "gx/ipfs/Qma6uuSyjkecGhMFFLfzyJDPyoDtNJSHJNweDccZhaWkgU/go-ipfs-cmds"
	cidenc "gx/ipfs/QmckgkstbdXagMTQ4e1DW2SzxGcjjudbqEvA5H2Rb7uvAT/go-cidutil/cidenc"
	cmdkit "gx/ipfs/Qmde5VP1qUkyQXKCfmEUA7bP64V2HAptbJ7phuPp7jXWwg/go-ipfs-cmdkit"
	mbase "gx/ipfs/QmekxXDhCxCJRNuzmHreuaT3BsuJcsjcXWNrtV9C8DRHtd/go-multibase"
)

var OptionCidBase = cmdkit.StringOption("cid-base", "Multibase encoding used for version 1 CIDs in output.")
var OptionOutputCidV1 = cmdkit.BoolOption("output-cidv1", "Upgrade CID version 0 to version 1 in output.")

// GetCidEncoder processes the `cid-base` and `output-cidv1` options and
// returns a encoder to use based on those parameters.
func GetCidEncoder(req *cmds.Request) (cidenc.Encoder, error) {
	return getCidBase(req, true)
}

// GetLowLevelCidEncoder is like GetCidEncoder but meant to be used by
// lower level commands.  It differs from GetCidEncoder in that CIDv0
// and not, by default, auto-upgraded to CIDv1.
func GetLowLevelCidEncoder(req *cmds.Request) (cidenc.Encoder, error) {
	return getCidBase(req, false)
}

func getCidBase(req *cmds.Request, autoUpgrade bool) (cidenc.Encoder, error) {
	base, _ := req.Options["cid-base"].(string)
	upgrade, upgradeDefined := req.Options["output-cidv1"].(bool)

	var e cidenc.Encoder = cidenc.Default

	if base != "" {
		var err error
		e.Base, err = mbase.EncoderByName(base)
		if err != nil {
			return e, err
		}
		if autoUpgrade {
			e.Upgrade = true
		}
	}

	if upgradeDefined {
		e.Upgrade = upgrade
	}

	return e, nil
}

func CidBaseDefined(req *cmds.Request) bool {
	base, _ := req.Options["cid-base"].(string)
	return base != ""
}

// CidEncoderFromPath creates a new encoder that is influenced from
// the encoded Cid in a Path.  For CidV0 the multibase from the base
// encoder is used and automatic upgrades are disabled.  For CidV1 the
// multibase from the CID is used and upgrades are eneabled.  On error
// the base encoder is returned.  If you don't care about the error
// condiation it is safe to ignore the error returned.
func CidEncoderFromPath(enc cidenc.Encoder, p string) (cidenc.Encoder, error) {
	v := extractCidString(p)
	if cidVer(v) == 0 {
		return cidenc.Encoder{Base: enc.Base, Upgrade: false}, nil
	}
	e, err := mbase.NewEncoder(mbase.Encoding(v[0]))
	if err != nil {
		return enc, err
	}
	return cidenc.Encoder{Base: e, Upgrade: true}, nil
}

func extractCidString(p string) string {
	segs := path.FromString(p).Segments()
	v := segs[0]
	if v == "ipfs" || v == "ipld" && len(segs) > 0 {
		v = segs[1]
	}
	return v
}

func cidVer(v string) int {
	if len(v) == 46 && v[:2] == "Qm" {
		return 0
	} else {
		return 1
	}
}
