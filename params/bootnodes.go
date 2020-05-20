// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ethereum network.
var MainnetBootnodes = []string{

	// storeman boot nodes
	"enode://ed214e8ce499d92a2085e7e6041b4f081c7d29d8770057fc705a131d2918fcdb737e23980bdd11fa86f5d824ea1f8a35333ac6f99246464dd4d19adac9da21d1@127.0.0.1:27717",
	"enode://9c6d6f351a3ede10ed994f7f6b754b391745bba7677b74063ff1c58597ad52095df8e95f736d42033eee568dfa94c5a7689a9b83cc33bf919ff6763ae7f46f8d@127.0.0.1:17717",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
var TestnetBootnodes = []string{
	//"enode://ce45ef8eb81b66c6bfdd255cd0cad8e7f75bcfea39edd130012ca6f437a9518b9d61447adab25e189fa39673a33107efa1978e332851fcb407f2fbee2fe61961@34.208.125.29:17717", //nk050
	//"enode://23208f45ddc7c96f635459a1740011746abfc3a74c02794512d84537adfcf51e8de89391227e84a09e2a3fdc954b247a3307788bc7cb576b80c6ef3bc6e456b3@35.167.106.109:17717",//nk049
	//"enode://01b31508ed81c2c6f1bb941faee1f6b0c5445994311a8bcc916c04f5693e19d8a4fdc7b4bfbd1aa65bcf32dc4a222fa5709e52e5fc624b81ef8edd79ca281e88@52.89.169.52:17717", //nk048
	//"enode://84c15ab07550c70b5b96d1e40d2347a718bb08500f804a4ca621a07ad66bb78503bd1e1711ff9ae585ad5bf154d72f314aed1b49bc42f65614a4ebadeb8cba0d@35.160.61.159:17717",//nk047
	"enode://1618570be0da74f7c9dabe674a85cdc2f13d87808a0f51e69400a74d42180413669e5898c0a8caa7a40226a9d6c858f1419e240ecb30f14d9eb42f5bf690f356@13.52.155.165:17717",
	"enode://b7b52602bebb302b386ff9c419e5c09c9bb3fb43cf19cf81b47cbccc252e2b67b3826ee102f5bca71f36e5856c2901f30cc3fa123824b76453b20e25a69ed0ec@52.52.142.120:17717",
}

var InternalBootnodes = []string{
	//"enode://dea09d1ded799044d3b8b5c66e28e584ea3fdaae12e0e39bb3491ac99424cc6c098f32e978c4aef1c3382c3c4492d7a33d720eabdee78cddba28541d6bef1bdc@52.53.224.4:17717", //
	//"enode://9e41c167954d33f5f5b7740a0f6a03b90ddab423cfd4e1fc6c844feff32e3a5d82e76c20d1823915676b58505efb6d33ea1fa6f7e6e22812b1d7ae7a90874881@118.190.33.102:17717",
	"enode://f0b604a19b711d20e60912fa2acc3b4966b3f1d8339bbcdd134d51e6ef7245927c09804aba1b8f94be5f5a77f062952ae58760ee7426299c16f9d19bf47be732@54.183.96.28:17717",
}

// PlutoBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Pluto test network.
var PlutoBootnodes = []string{
	"enode://81ffab14284d29f9a87737780717719666af814d78057ec4b6799b9d275c41e8041887c733a95fba13ae9fba4fb3026d5f53993143b83ab6648cae3b1e5e9c35@18.236.236.189:17717",
	"enode://ca5496aa6eda6403f4ac41e7841d1ae6d963a321afcb8c59c0f2935f837bd2300ff258ab94bf4db375257f29898e5d0ea5903c28a0e9a41a4aba4e100b4b2ed0@34.212.171.224:17717",
	"enode://86989aacffbc22640dee74864ac0f17fb4987ab0b6792a6fd14801557e7f7ff6447d77945d64463d5e6e0bed5ac257c842d9631abe8e68110e1aa9233ad4e3a1@54.184.26.209:27717",
}
var PlutoV5ootnodes = []string{
	"enode://81ffab14284d29f9a87737780717719666af814d78057ec4b6799b9d275c41e8041887c733a95fba13ae9fba4fb3026d5f53993143b83ab6648cae3b1e5e9c35@18.236.236.189:17717?discport=17718",
}

var InternalV5Bootnodes = []string{
	"enode://81ffab14284d29f9a87737780717719666af814d78057ec4b6799b9d275c41e8041887c733a95fba13ae9fba4fb3026d5f53993143b83ab6648cae3b1e5e9c35@18.236.236.189:17717",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
	"enode://dfa95c2be31b3541895df452678355ef4b38988863959fe56d1a217cd9fdeee27024cb13688cab56373ac597968aa2faf0da8cd87f4238366ddb41f03fc78884@118.190.33.68:17717", // IE
	"enode://9e41c167954d33f5f5b7740a0f6a03b90ddab423cfd4e1fc6c844feff32e3a5d82e76c20d1823915676b58505efb6d33ea1fa6f7e6e22812b1d7ae7a90874881@118.190.33.102:17717",
}
