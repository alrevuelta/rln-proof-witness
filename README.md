# rln-proof-witness

This repo showcases how to create RLN zk proofs with a custom witness, which is, the Merkle proof that a given leaf belongs to the tree. The main advantage is that generating such proofs no longer requires the whole tree, but just the Merkle proof, which has a fixed size of 20 elements, assuming a tree of depth 20. This is a proof of concept not meant to be used in production.

This feature is useful for light clients since they no longer need to sync the whole tree. The Merkle proof used to generate the RLN proof can be fetched from:
* 1) From other waku nodes: Any waku node offering this service (see [branch](https://github.com/waku-org/go-waku/compare/master...merkle-proof-provider))
* 2) From the contract: See [ongoing work](https://github.com/privacy-scaling-explorations/zk-kit/issues/123)

## i) From other waku nodes

A service like [this one](https://github.com/waku-org/go-waku/compare/master...merkle-proof-provider) can be offered by any waku node, and returns the Merkle proof and root of a given commitment in the tree. For example, lets say that a light client with public commitment `21235824865182058647676208664819393617711826061506661756580202797779091020767` wants to create a RLN zk proof for a given message. It can use the above fork as:

```
curl http://65.21.94.244:30304/debug/v1/merkleProof/21235824865182058647676208664819393617711826061506661756580202797779091020767
```

Which returns all needed elements to generate a RLN proof. With this plus the message, epoch and secret, the light client can generate the proof. No tree needed.

```
{"root":"f46c2005a3c47c0ded1707b0396262e44d2b73b1e6fe3478858af96f62167b05","pathElements":["4693987cdebc611af84591309753f4de85671a9700d980e4b07e04cca9664a1d","ba5539dd1fa12981146ca436ae05579bee166e15ec97f585bba724c7847b0409","4d4d9b70341b80acb9ecfe1b57f8f1294786f45966a48f43be840f68298e5f1c","38d256b8b27ed528d51d3750ea6e7c460621f7508d753d2eafe27e533133f418","409972be02123b9b7c3aa33931f211aa4831fb3e86dfd94aa113ccb712e93a1a","084e147f355fd170c063dfd0a5b5bc646668d6ad19d62b2136749bf62556b910","4c27f1cfef26fc37bdd76ffd0aa928a8784588884f1a130b4cc4319fa6d03903","78e433d9574e23708f16083c46f5ad72bac80054371700f9a8260938ead20705","61ccf3993abe4c441a21414a272e6b612a47644586ec1b50a627608ff1e5a52f","47d7fc14a656213eab28e2e3cc7a5ee4661f949e3880b7ec21fdd8d07643880e","f20a19dae57561de33357157f99258f969b42ea5d17a71281e4f4972da01721b","36767dcefa6bbcbeb5080865e4e1e6a619982401b2c0005238365e7222888d1f","5af8b571049a87d0a888cf2aa1b06261fbfc8cba891570b9af4b916cf6825d2c","d0bfbfe070f2586464f413a1aac4f54e13a13fdf5a7f9520b80b94a04841c514","0ce8ebf44b8e1116d489ad8c5825be11afb9d844eec0101e966f982fb1330d19","926ce0259364b3a50a51af9665ae6711ed73ad14493517ac524170cea98af922","2373ba8bd353b7f8eecc6ec6296f525a576abf728d226f9f0b88e56c9b7c7c2a","92b9363f64dd754d958b98c2c9430047fc3f464dc1f97ac6c18e6958e586812e","0ff11f1c9d24463527927364ad6eef8a94ae0d05cfc8e249ab4e9a1e57c5570f","ca2cf73461e39c3ce4467d6910e378fe1c0e8088433df6d54a55fbb567ee3018"],"pathIndexes":"AQEBAAEBAQEAAAAAAAAAAAAAAAA=","leafIndex":247,"commitmentId":"21235824865182058647676208664819393617711826061506661756580202797779091020767"}
```

Run the code as:
```
go run main.go
```

Notes:
* New memberships registered by others modify the tree, so said Merkle proof is not valid forever.
* Use an endpoint that you trust, otherwise you leak your ip/commitment.

## ii) From the contract

Under development