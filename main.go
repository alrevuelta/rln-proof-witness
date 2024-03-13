package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/waku-org/go-zerokit-rln/rln"
)

// See this branch with endpoint to retrieve merkle proof:
// https://github.com/waku-org/go-waku/compare/master...merkle-proof-provider
// Deployed in sandbox machine as proof of concept.
const MerkleProofProdiver = "http://65.21.94.244:30312"

func main() {
	rlnInstance, err := rln.NewRLN()
	if err != nil {
		log.Fatal(err)
	}

	// Proof generation for this membership:
	// https://sepolia.etherscan.io/tx/0x039ded260a587ad14262a1690e604adaf4a7326cbeb028095dc081f772c0cb44
	idCommitment, err := rln.ToBytes32LE("df5b05a4d6c3f5a3aee2ea4664a1ccd529933058a2d82931153b6d0c0a0df32e")
	if err != nil {
		log.Fatal(err)
	}
	idCommitmentBig := new(big.Int).SetBytes(idCommitment[:])

	// Only the owner will know this secret
	idSecretHash, err := rln.ToBytes32LE("bad944102c2798ce77e8feb37d264f64b033e847f2e9ac29acb527543ea8f008")
	if err != nil {
		log.Fatal(err)
	}

	merkleProof, err := GetMerkleProof(idCommitmentBig)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("idCommitmentBig: ", idCommitmentBig, " index: ", merkleProof.LeafIndex, " merkle root: ", merkleProof.MerkleRoot)

	// RLN proof generation
	someMessage := []byte("some message") // Your message
	epoch := rln.ToEpoch(0)               // Epoch timestsamp in seconds

	// We create a witness with the merkle proof and our secret
	witness := rln.CreateWitness(idSecretHash, someMessage, epoch, rln.MerkleProof{
		PathIndexes:  ConvertStringToUint8(merkleProof.MerkePathIndexes),
		PathElements: ConvertStringToBytes(merkleProof.MerkePathElements),
	})

	// We can generate a RLN proof without having the whole tree
	rlnProof, err := rlnInstance.GenerateRLNProofWithWitness(witness)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("rln proof is: ", rlnProof)
}

// TODO: Duplicated from go code
// See: https://github.com/waku-org/go-waku/compare/master...merkle-proof-provider
type MerkleProofResponse struct {
	MerkleRoot        string   `json:"root"`
	MerkePathElements []string `json:"pathElements"`
	MerkePathIndexes  []string `json:"pathIndexes"`
	LeafIndex         uint64   `json:"leafIndex"`
	CommitmentId      string   `json:"commitmentId"`
}

func GetMerkleProof(publicCommitment *big.Int) (*MerkleProofResponse, error) {
	url := fmt.Sprintf("%s/debug/v1/merkleProof/%s", MerkleProofProdiver, publicCommitment.String())
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get merkle proof")
	}
	defer resp.Body.Close()
	merkleProof := &MerkleProofResponse{}

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read response body")
		}

		if err = json.Unmarshal(bodyBytes, merkleProof); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal response body")
		}
	}

	return merkleProof, err
}

func ConvertStringToUint8(pathIndexesStr []string) []uint8 {
	// Convert merkle proof from string to uint8
	pathIndexes := make([]uint8, 0)
	for _, pathIndex := range pathIndexesStr {
		pathIndexUint8, err := strconv.ParseUint(pathIndex, 10, 8)
		if err != nil {
			log.Fatal(err)
		}
		pathIndexes = append(pathIndexes, uint8(pathIndexUint8))
	}
	return pathIndexes
}

func ConvertStringToBytes(pathElementsStr []string) [][32]byte {
	// TODO: check commitment id match
	// TODO: check merkle proof is valid against root.

	// Convert merkle proof from string to [32]byte
	pathElements := make([][32]byte, 0)
	for _, pathElement := range pathElementsStr {
		pathElementBytes, err := hex.DecodeString(pathElement)
		if err != nil {
			log.Fatal(err)
		}
		var pathElement32 [32]byte
		copy(pathElement32[:], pathElementBytes)
		pathElements = append(pathElements, pathElement32)
	}
	return pathElements
}
