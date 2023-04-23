// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/cespare/cp"
)

// These tests are 'smoke tests' for the account related
// subcommands and flags.
//
// For most tests, the test files from package accounts
// are copied into a temporary keystore directory.

func tmpDatadirWithKeystore(t *testing.T) string {
	datadir := tmpdir(t)
	keystore := filepath.Join(datadir, "keystore")
	source := filepath.Join("..", "..", "accounts", "keystore", "testdata", "keystore")
	if err := cp.CopyAll(keystore, source); err != nil {
		t.Fatal(err)
	}
	return datadir
}

func TestAccountListEmpty(t *testing.T) {
	geth := runGeth(t, "account", "list")
	geth.ExpectExit()
}

func TestAccountList(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	geth := runGeth(t, "account", "list", "--datadir", datadir)
	defer geth.ExpectExit()
	if runtime.GOOS == "windows" {
		geth.Expect(`
Account #0: {ff8fd0e9064bcbc5462763229b8a4314f07cf000} keystore://{{.Datadir}}\keystore\UTC--2023-01-17T00-37-32.112011400Z--ff8fd0e9064bcbc5462763229b8a4314f07cf000
Account #1: {9279739dce240363860f3672afa9e4c1e1f86ccc} keystore://{{.Datadir}}\keystore\aaa
Account #2: {31cdf7786b0b26b733b545988701c6ed0601f848} keystore://{{.Datadir}}\keystore\zzz
`)
	} else {
		geth.Expect(`
Account #0: {ff8fd0e9064bcbc5462763229b8a4314f07cf000} keystore://{{.Datadir}}/keystore/UTC--2023-01-17T00-37-32.112011400Z--ff8fd0e9064bcbc5462763229b8a4314f07cf000
Account #1: {9279739dce240363860f3672afa9e4c1e1f86ccc} keystore://{{.Datadir}}/keystore/aaa
Account #2: {31cdf7786b0b26b733b545988701c6ed0601f848} keystore://{{.Datadir}}/keystore/zzz
`)
	}
}

func TestAccountNew(t *testing.T) {
	geth := runGeth(t, "account", "new", "--lightkdf")
	defer geth.ExpectExit()
	geth.Expect(`
Your new account is locked with a password. Please give a password. Do not forget this password.
!! Unsupported terminal, password will be echoed.
Password: {{.InputLine "foobar"}}
Repeat password: {{.InputLine "foobar"}}

Your new key was generated
`)
	geth.ExpectRegexp(`
Public address of the key:   0x[0-9a-fA-F]{40}
Path of the secret key file: .*UTC--.+--[0-9a-f]{40}

- You can share your public address with anyone. Others need it to interact with you.
- You must NEVER share the secret key with anyone! The key controls access to your funds!
- You must BACKUP your key file! Without the key, it's impossible to access account funds!
- You must REMEMBER your password! Without the password, it's impossible to decrypt the key!
`)
}

func TestAccountImport(t *testing.T) {
	key1, err := cryptobase.SigAlg.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	privkeyhex, err := cryptobase.SigAlg.PrivateKeyToHex(key1)
	if err != nil {
		t.Fatal(err)
	}

	addr, err := cryptobase.SigAlg.PublicKeyToAddress(&key1.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct{ name, key, output string }{
		{
			name:   "correct account",
			key:    privkeyhex,
			output: "Address: {" + strings.ToLower(addr.StringNoHex()) + "}\n",
		},
		{
			name:   "invalid character",
			key:    privkeyhex + "1",
			output: "Fatal: Failed to load the private key from file: invalid character '1' at end of key file\n",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			importAccountWithExpect(t, test.key, test.output)
		})
	}
}

func importAccountWithExpect(t *testing.T, key string, expected string) {
	dir := tmpdir(t)
	keyfile := filepath.Join(dir, "key.prv")
	if err := ioutil.WriteFile(keyfile, []byte(key), 0600); err != nil {
		t.Error(err)
	}
	passwordFile := filepath.Join(dir, "password.txt")
	if err := ioutil.WriteFile(passwordFile, []byte("foobar"), 0600); err != nil {
		t.Error(err)
	}
	geth := runGeth(t, "account", "import", keyfile, "-password", passwordFile)
	defer geth.ExpectExit()
	geth.Expect(expected)
}

func TestAccountNewBadRepeat(t *testing.T) {
	geth := runGeth(t, "account", "new", "--lightkdf")
	defer geth.ExpectExit()
	geth.Expect(`
Your new account is locked with a password. Please give a password. Do not forget this password.
!! Unsupported terminal, password will be echoed.
Password: {{.InputLine "something"}}
Repeat password: {{.InputLine "something else"}}
Fatal: Passwords do not match
`)
}

func TestAccountUpdate(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	geth := runGeth(t, "account", "update",
		"--datadir", datadir, "--lightkdf",
		"ff8fd0e9064bcbc5462763229b8a4314f07cf000")
	defer geth.ExpectExit()
	geth.Expect(`
Unlocking account ff8fd0e9064bcbc5462763229b8a4314f07cf000 | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Password: {{.InputLine "foobar"}}
Please give a new password. Do not forget this password.
Password: {{.InputLine "foobar2"}}
Repeat password: {{.InputLine "foobar2"}}
`)
}

func TestWalletImport(t *testing.T) {
	geth := runGeth(t, "wallet", "import", "--lightkdf", "testdata/guswallet.json")
	defer geth.ExpectExit()
	geth.Expect(`
!! Unsupported terminal, password will be echoed.
Password: {{.InputLine "foo"}}
Address: {7522bAf8770939E65CEF366b7A30204f1dbc023f}
`)

	files, err := ioutil.ReadDir(filepath.Join(geth.Datadir, "keystore"))
	if len(files) != 1 {
		t.Errorf("expected one key file in keystore directory, found %d files (error: %v)", len(files), err)
	}
}

func TestWalletImportBadPassword(t *testing.T) {
	geth := runGeth(t, "wallet", "import", "--lightkdf", "testdata/guswallet.json")
	defer geth.ExpectExit()
	geth.Expect(`
!! Unsupported terminal, password will be echoed.
Password: {{.InputLine "wrong"}}
Fatal: could not decrypt key with given password
`)
}

func TestUnlockFlag(t *testing.T) {
	geth := runMinimalGeth(t, "--port", "0", "--ipcdisable", "--datadir", tmpDatadirWithKeystore(t),
		"--unlock", "9279739dce240363860f3672afa9e4c1e1f86ccc", "js", "testdata/empty.js")
	geth.Expect(`
Unlocking account 9279739dce240363860f3672afa9e4c1e1f86ccc | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Password: {{.InputLine "foobar"}}
`)
	geth.ExpectExit()

	wantMessages := []string{
		"Unlocked account",
		"=0x9279739dce240363860f3672afa9e4c1e1f86ccc",
	}
	for _, m := range wantMessages {
		if !strings.Contains(strings.ToLower(geth.StderrText()), strings.ToLower(m)) {
			t.Errorf("stderr text does not contain %q", m)
		}
	}
}

func TestUnlockFlagWrongPassword(t *testing.T) {
	geth := runMinimalGeth(t, "--port", "0", "--ipcdisable", "--datadir", tmpDatadirWithKeystore(t),
		"--unlock", "9279739dce240363860f3672afa9e4c1e1f86ccc", "js", "testdata/empty.js")

	defer geth.ExpectExit()
	geth.Expect(`
Unlocking account 9279739dce240363860f3672afa9e4c1e1f86ccc | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Password: {{.InputLine "wrong1"}}
Unlocking account 9279739dce240363860f3672afa9e4c1e1f86ccc | Attempt 2/3
Password: {{.InputLine "wrong2"}}
Unlocking account 9279739dce240363860f3672afa9e4c1e1f86ccc | Attempt 3/3
Password: {{.InputLine "wrong3"}}
Fatal: Failed to unlock account 9279739dce240363860f3672afa9e4c1e1f86ccc (could not decrypt key with given password)
`)
}

// https://github.com/ethereum/go-ethereum/issues/1785
func TestUnlockFlagMultiIndex(t *testing.T) {
	geth := runMinimalGeth(t, "--port", "0", "--ipcdisable", "--datadir", tmpDatadirWithKeystore(t),
		"--unlock", "9279739dce240363860f3672afa9e4c1e1f86ccc", "--unlock", "0,2", "js", "testdata/empty.js")

	geth.Expect(`
Unlocking account 0 | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Password: {{.InputLine "foobar"}}
Unlocking account 2 | Attempt 1/3
Password: {{.InputLine "foobar"}}
`)
	geth.ExpectExit()

	wantMessages := []string{
		"Unlocked account",
		"=0xFF8fd0e9064bCBc5462763229b8a4314F07cF000",
		"=0x31cDF7786B0B26b733B545988701c6Ed0601f848",
	}
	for _, m := range wantMessages {
		if !strings.Contains(geth.StderrText(), m) {
			t.Errorf("stderr text does not contain %q", m)
		}
	}
}

func TestUnlockFlagPasswordFile(t *testing.T) {
	geth := runMinimalGeth(t, "--port", "0", "--ipcdisable", "--datadir", tmpDatadirWithKeystore(t),
		"--unlock", "9279739dce240363860f3672afa9e4c1e1f86ccc", "--password", "testdata/passwords.txt", "--unlock", "0,2", "js", "testdata/empty.js")

	geth.ExpectExit()

	wantMessages := []string{
		"Unlocked account",
		"=0xff8fd0e9064bcbc5462763229b8a4314f07cf000",
		"=0x31cdf7786b0b26b733b545988701c6ed0601f848",
	}
	for _, m := range wantMessages {
		if !strings.Contains(strings.ToLower(geth.StderrText()), strings.ToLower(m)) {
			t.Errorf("stderr text does not contain %q", m)
		}
	}
}

func TestUnlockFlagPasswordFileWrongPassword(t *testing.T) {
	geth := runMinimalGeth(t, "--port", "0", "--ipcdisable", "--datadir", tmpDatadirWithKeystore(t),
		"--unlock", "9279739dce240363860f3672afa9e4c1e1f86ccc", "--password",
		"testdata/wrong-passwords.txt", "--unlock", "0,2")
	defer geth.ExpectExit()
	geth.Expect(`
Fatal: Failed to unlock account 0 (could not decrypt key with given password)
`)
}

func TestUnlockFlagAmbiguous(t *testing.T) {
	store := filepath.Join("..", "..", "accounts", "keystore", "testdata", "dupes")
	geth := runMinimalGeth(t, "--port", "0", "--ipcdisable", "--datadir", tmpDatadirWithKeystore(t),
		"--unlock", "9279739dce240363860f3672afa9e4c1e1f86ccc", "--keystore",
		store, "--unlock", "9279739dce240363860f3672afa9e4c1e1f86ccc",
		"js", "testdata/empty.js")
	defer geth.ExpectExit()

	// Helper for the expect template, returns absolute keystore path.
	geth.SetTemplateFunc("keypath", func(file string) string {
		abs, _ := filepath.Abs(filepath.Join(store, file))
		return abs
	})
	geth.Expect(`
Unlocking account 9279739dce240363860f3672afa9e4c1e1f86ccc | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Password: {{.InputLine "foobar"}}
Multiple key files exist for address 9279739dce240363860f3672afa9e4c1e1f86ccc:
   keystore://{{keypath "1"}}
   keystore://{{keypath "2"}}
Testing your password against all of them...
Your password unlocked keystore://{{keypath "1"}}
In order to avoid this warning, you need to remove the following duplicate key files:
   keystore://{{keypath "2"}}
`)
	geth.ExpectExit()

	wantMessages := []string{
		"Unlocked account",
		"=0x9279739dce240363860f3672afa9e4c1e1f86ccc",
	}
	for _, m := range wantMessages {
		if !strings.Contains(strings.ToLower(geth.StderrText()), strings.ToLower(m)) {
			t.Errorf("stderr text does not contain %q", m)
		}
	}
}

func TestUnlockFlagAmbiguousWrongPassword(t *testing.T) {
	store := filepath.Join("..", "..", "accounts", "keystore", "testdata", "dupes")
	geth := runMinimalGeth(t, "--port", "0", "--ipcdisable", "--datadir", tmpDatadirWithKeystore(t),
		"--unlock", "9279739dce240363860f3672afa9e4c1e1f86ccc", "--keystore",
		store, "--unlock", "9279739dce240363860f3672afa9e4c1e1f86ccc")

	defer geth.ExpectExit()

	// Helper for the expect template, returns absolute keystore path.
	geth.SetTemplateFunc("keypath", func(file string) string {
		abs, _ := filepath.Abs(filepath.Join(store, file))
		return abs
	})
	geth.Expect(`
Unlocking account 9279739dce240363860f3672afa9e4c1e1f86ccc | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Password: {{.InputLine "wrong"}}
Multiple key files exist for address 9279739dce240363860f3672afa9e4c1e1f86ccc:
   keystore://{{keypath "1"}}
   keystore://{{keypath "2"}}
Testing your password against all of them...
Fatal: None of the listed files could be unlocked.
`)
	geth.ExpectExit()
}
