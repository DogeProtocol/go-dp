package hybridedsfull

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/DogeProtocol/dp/crypto/hybrideds"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"testing"
)

func testHybridedsfullSigBasic(t *testing.T) {
	var sig signaturealgorithm.SignatureAlgorithm

	sig = CreateHybridedsfullSig()
	signaturealgorithm.SignatureAlgorithmTest(t, sig)
}

func TestHybridedsfullSig_Basic(t *testing.T) {
	testHybridedsfullSigBasic(t)
}

func testBase64(t *testing.T) {
	var sig signaturealgorithm.SignatureAlgorithm
	sig = CreateHybridedsfullSig()

	serverSigningPublicKey := "gIIYJ9zt9om5ccXlkKP3T1fqco6PFGh7KcHBV5BYtaCcbiGlS1uHJolV+cXHq4MAgH8IK61niyTwLw9Dth5fZreiWXp+4AbV/1rquS77evhbCCcKfghASdtcU1xFRbYR4nFC8OHIOmvlTggXZSshp0L3XQpNYly4ZU7FsexhgW5BQyTXb+Q2AAfj2U70t1snpeG8pQRBCaNI0DuuG9rU3X/dFFvGf9LbxxZFl1+9vsgV1n1/keJyfRIa9Yf8bT4xiM5ksgGUfDHbHKsAGMmL25iynIn5SpmNJ8JckrU4M3a31bVeStydQs1NaBkbu9pFizKYDpQsm50dMcD7oLpPed+Oxq18pYx2A/A6p1HpXe7p4n9JwrWwGcltoETQ9vMAcxfuT5Pa5Vo74FT1BU2SLHH+2BtwG7tzMQYVZqwNsZMI573dA12en5tGuelAXPmpc+zqIk0WEyH9TZm0/aSFp0+Bem4eAvLT3RM8Jx2M31aIX1fP1DjSO/Gb5R3uDrPBc8sdEuA2u3nJFjthWaxQRzCO98xmQ0SNMnaE52W2amBxBNRaw5wOzCvoZZSeGJGMMpyEZ5R/8btTAH5WHDRkdG7taAo3h6aQCZAiv9ygh9BnIZWLIJBlbiJ3EycLsSqybJG4AItHES+QDsq+/YBig0QhAbuIp0sEovb8IlSV+CZnRiT+Qr0nLZ9uGgFGRChP/XwomtvAC8ld1cfOagL6lz13/AIahWXhlOzBhRyrQuDutd7nQwmZsK6OlhqMF+a5fPPmDGvNJ7S0kXW+SrMrz95PiTcezBdLaVndyD3/O8pKPTHGPjltYdGMgLggKE/weRe9HUaYUaQkELVeCuh4HgA5kdobRAtOjEoL0qfvBPpWe8kBfrFyfXPiInynIICydRajmkyDHeWbg92qF9zI70b8pekcRB1Rx1bJu8NQBE5Annr/yy6rjuxIHgxDLhAScK33s8Ncc8gA2bCqYVKstxLhD4/dCQo093EGclAyeUmn+wt3GIKkAiANeuWZhATDWNGpiKwTWO4xdIhGtysc2rchKnWmAtWJXPBJnKBMepWHkZksJx0UoGVbYhdPsWUF6pFSX5FDIm2PNDcx8i9YLN39PSv18PpG744116gk7f5djvuDbuagny1N5ziL7WPqKriO9cGDm1G6SPd+QL89sLh5eNrLCqO48h8zS6/77z90dat35Qy8ASlmRRp8MSyB2V2IZeNpactfWA8ja9tKD5UiF9EHgxsbN2+A/VMKsg7te+6bDDpOR89XUBY17HsDvxdqJWkacv+fa2z2QBauizoL7m73EIC/KxogQb3F1CtIHeGKHkFpT/nMhR1MXkDOtYy0jeNZ8knxdYeDFAN89esEVxthRBonZ0Y8rjreb0uEY3IfpnePo57DTIbAVxwA5MAsbAab9S59kTrdLSxiRQ1g3sMuWzmxD69fFk/o8faGQH0Rjim5u3oNWzQpoZdPJAh2qN0qcJIOrTjFxmeQw912JO2YvMzcLwhpCefeYdlTK9H6johS7zbbg71cX1pUCKq6ThLMW1WnypHY1edDxwa+vpGDc+qyY8y4Lc9yys/1iSf7VCTZfQRVS5UaLYYqPPY+eVx2FPgDvbfkap2DYyO2BjSmj2emS/Qb267YwkaVO0VwDSkP6qzFe0+e+cO3QFlDZlJnXROxFLTZBT1PiirCU2uGRTiZFKOo2lIr8ooSCiituQjIZk9pbHJs0k7C8lKI+cH9CwSuw1QXZ1Gsy+11Zc+xpdb3AoqEZwQ4WvhglqNd0mYP+kZMM3cCQ+ZAfZ3TtSxaugST94F/Y/TBmomQPM7PIfP8+KBBRzsmDhKj4Xi73d4ouGHXTpSKHKjRRNjdbX+JnFnOR2bs1kuNNg=="
	signature := "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"
	serverPubKeyDataRemote := "dzyetSzZ7c1dyJ9URuTmK+R7u0sBYgAfBe12h11sZEfxEYoJs+UDbMNMBCPce+Z3mCN1rQ0U4sDMqMX+Su1bGUUT6qah128Fny2pg5d33jt2jIVS6tajeo/fWJm8EXUEW4lP55eF7UIbYo/WgkuzC3PIfLH/BJuDmFanhsUBtB6LDyUcAE48M5mXPDnCf0+DOtFXERT3Zg7CU/ioNf2GQVxCNdy+9bRQY46/VcgJOgEe2x8nrkgYgIm9fBSyTz31hEYhPYV9BS00idl8BgDq8HuCsqJJuYGXfSLoX6Nbq3C6+UKLiSuYJ/tS20t0sPckoZEDgfC4nAcvNPzdS2vRS3XyoJA3Hxp3k5uDI0oqJVK4bo2P5H5PyazXzjYQTPuFS+FS5iONDIPxXSCqJnouzZTjX3kHvMn7Suql7QDIMKaxBAwxEV8+BO0+Pt6LVU3B8U+KK9A+p3OBQ7oQapYLQkwIlNNnh5ScinRylnHo4HybItdvtbMOznl2mrJFNAp0/BuFfTSPBvi4NXNs7StoK59Z6G+Pm8cDKI7dXH3NlNe09r9hsHkxjEc+zYYKZwbp9HbkawGvMjAF0vpdBCQ8DFfZNbv1maDZ23Wa/0+nPsxYGv3E12/GB5EyBiyygGpRc5ADdD4CIYKgHGtZy0hOsRvSUQvu1TtHQT5u4LC025p23l20LjKOAHa7bO63mKNalrLHN1ryIUhsb9zt6X17iwQ8c0I/Seqk+7+uaDUvAlpy/OE02tCB3lfA9fLF5stGqAUXv+eaeMEVBXKUy3e7KzcvlcA0r/eb5a8oEE3JHELIPLif05PDLEloMikDnYd9rwDcvztKdFoaiiHWN6Y33hTpJw0pGzmc2Sjh2/R7mm7Brtivte92eWJSc6kjAYiQsXZMrev7nyeIGK2zhanDdkLE+APCgqdT2ZOu4EghB56KNKceprlm/VzLwnRxDUfIUEk5udjvuxsQUcxcZdRSbDb19sQM13t8+KdVGYeViyXBKy3ds+nVFm/Vmf1P7VokhvPi42nTRSxZKQKvssz0qgHlePMKvqUqj6rfNY6DTLj3hbaGr57EVruxStUFZacswiPGPOXc1UNe9MDvb5saQG5SMz+B39Pgzgbu2ZhYXydC3X6TQ81iTDLGgrfdulyLQ63qCdXiM3zqH8+/9iqU/TkqUWam1scJ8lPX+EO7iMaVgBGGTkm109uEyP4HJRfqUIC+7aQ1VZRG1utzoNkn2FtdA00Sq2C4Hg4IEfg7dQ2yJXW73gZGUVV8CqpxRAUtJo9iG7S5DbG76+nkAg1mWsHL+eDT68q+57UNipOIDxsK4VPtMWKhnXajPeQ4O0ABKT3XY7VodRKW3KtVz6HqSPbEZ7FM7uhZDAY2IgvyTRBA0DI7BGS0b+DC+gdCaVhxl2s6FS/yzU82W5ZtHaMxKxH3Koa8xyH7UqMKjpDCwuggSrmdwkxZUOxEvZWTqTEyC70Gk0HCiBIGsL6/FW2kUIs98D5OIu1oPMRGNMjz1w2Mmtpr40LIzje51nTMcsOZ2c/ouQhTzBheksiP5UpTe5Ow0Wc0vKplGwFJQX+JzirLx5SCk7vJj591nOPRh8v5NmYm3RashOgX8NTqd+IvKfiyy0rH+2Dm5zv3Zmy15+pqFLuXl2WOSXIs2XuOwwluixn1KvAQfsWBd3Edog0kQ6gyFwHTGpoGqSsNYhuaAh+L8pJEjTwcQL6bANpSriPGo4MLulFBi2/nPK33rDpA8FeYXwCu+eb9RQK29TTovO3flmbiYKRYZ2n6W7kOhA+xL32Iz6X9BeJAhyqEXuBjlqT6dHeCZEe9ApAlL64PKrSfiVPVvq3fe4hJCnFugT/m0s4WTHK7jcNIfENG9Balsw=="

	serverSigningPublicKeyBytes, err := base64.StdEncoding.DecodeString(serverSigningPublicKey)
	if err != nil {
		t.Fatalf("failed")
	}
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		t.Fatalf("failed")
	}
	fmt.Println("signatureBytes", len(signatureBytes))
	serverPubKeyDataRemoteBytes, err := base64.StdEncoding.DecodeString(serverPubKeyDataRemote)
	if err != nil {
		t.Fatalf("failed")
	}

	pubKeyExpected, err := sig.DecodePublicKey(serverSigningPublicKeyBytes)
	if err != nil {
		t.Fatalf("failed")
	}
	addressExpected, err := sig.PublicKeyToAddress(pubKeyExpected)
	if err != nil {
		t.Fatalf("failed")
	}

	pubKeyActual, err := sig.DecodePublicKey(serverPubKeyDataRemoteBytes)
	if err != nil {
		t.Fatalf("failed")
	}
	addressActual, err := sig.PublicKeyToAddress(pubKeyActual)
	if err != nil {
		t.Fatalf("failed")
	}

	fmt.Println("addressExpected", addressExpected, "addressActual", addressActual)
}

func TestBase64(t *testing.T) {
	testBase64(t)
}

func TestCompactAndFullInterop(t *testing.T) {
	var sigFull signaturealgorithm.SignatureAlgorithm
	sigFull = CreateHybridedsfullSig()

	var sigCompact signaturealgorithm.SignatureAlgorithm
	sigCompact = hybrideds.CreateHybridedsSig(true)

	keyCompact, err := sigCompact.GenerateKey()
	if err != nil {
		t.Fatalf(err.Error())
	}

	serializedCompact, err := sigCompact.SerializePrivateKey(keyCompact)
	if err != nil {
		t.Fatalf(err.Error())
	}

	keyFull, err := sigFull.DeserializePrivateKey(serializedCompact)
	if err != nil {
		t.Fatalf(err.Error())
	}

	addrCompact, err := sigCompact.PublicKeyToAddress(&keyCompact.PublicKey)
	if err != nil {
		t.Fatalf(err.Error())
	}

	addrFull, err := sigFull.PublicKeyToAddress(&keyFull.PublicKey)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if addrFull.IsEqualTo(addrCompact) == false {
		t.Fatalf("failed")
	}

	if bytes.Compare(keyCompact.PubData, keyFull.PubData) != 0 {
		t.Fatalf("failed")
	}

	digestHash1 := []byte(testmsg1)
	signature1, err := sigFull.Sign(digestHash1, keyFull)
	if err != nil {
		fmt.Println(err)
		t.Fatal("Sign failed")
	}

	if sigFull.Verify(keyCompact.PubData, digestHash1, signature1) != true { //compact pub
		t.Fatal("Verify failed")
	}

	signature2, err := sigCompact.Sign(digestHash1, keyFull)
	if err != nil {
		fmt.Println(err)
		t.Fatal("Sign failed")
	}

	if sigCompact.Verify(keyFull.PubData, digestHash1, signature2) != true { //full pub
		t.Fatal("Verify failed")
	}
}