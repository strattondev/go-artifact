package deckcode

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

const (
	s_nCurrentVersion = 2
	sm_rgchEncodedPrefix = "ADC"
)

type Deck struct {
	Heroes []Hero
	Cards []Card
	Name string
}

type Hero struct {
	Id int
	Turn int
}

type Card struct {
	Id int
	Count int
}

func ParseDeck(strDeckCode string) (Deck, error) {
	if deckBytes, err := decodeDeckString(strDeckCode); err != nil {
		return Deck{}, err
	} else {
		return parseDeckInternal(strDeckCode, deckBytes)
	}
}

func decodeDeckString(strDeckCode string) ([]byte, error) {
	if !strings.HasPrefix(strDeckCode, sm_rgchEncodedPrefix) {
		return nil, errors.New("Deck code does not contain proper prefix")
	}

	strDeckCode = strDeckCode[3:]
	strDeckCode = strings.Replace(strDeckCode, "-", "/", -1)
	strDeckCode = strings.Replace(strDeckCode, "_", "=", -1)

	return base64.StdEncoding.DecodeString(strDeckCode)
}

func parseDeckInternal(strDeckCode string, deckBytes []byte) (Deck, error) {
	var nCurrentByteIndex = 0
	var nTotalBytes = len(deckBytes)

	var nVersionAndHeroes = int(deckBytes[nCurrentByteIndex])
	nCurrentByteIndex++
	var version = nVersionAndHeroes >> 4

	if s_nCurrentVersion != version && version != 1 {
		return Deck{}, errors.New(fmt.Sprintf("Expected version %d but received version %d", s_nCurrentVersion, version))
	}

	var nCheckSum = int(deckBytes[nCurrentByteIndex])
	nCurrentByteIndex++
	var nStringLength = 0

	if version > 1 {
		nStringLength = int(deckBytes[nCurrentByteIndex])
		nCurrentByteIndex++
	}

	var nTotalCardBytes = nTotalBytes - nStringLength

	{
		var nComputedCheckSum = 0

		for i := nCurrentByteIndex; i < nTotalCardBytes; i++ {
			nComputedCheckSum += int(deckBytes[i])
		}

		var masked = nComputedCheckSum & 0xFF

		if nCheckSum != masked {
			return Deck{}, errors.New("Checksum did not match mask")
		}
	}

	var nNumHeroes = 0

	if !readVarEncodedUint32(nVersionAndHeroes, 3, deckBytes, &nCurrentByteIndex, nTotalCardBytes, &nNumHeroes) {
		return Deck{}, errors.New("Could not read in hero count")
	}

	var heroes = make([]Hero, 0)
	var nPrevCardBase = 0
	{
		for nCurrHero := 0; nCurrHero < nNumHeroes; nCurrHero++ {
			var nHeroTurn = 0
			var nHeroCardID = 0

			if !readSerializedCard(deckBytes, &nCurrentByteIndex, nTotalCardBytes, &nPrevCardBase, &nHeroTurn, &nHeroCardID) {
				return Deck{}, errors.New("Could not read in hero card")
			}

			heroes = append(heroes, Hero{
				Id: nHeroCardID,
				Turn: nHeroTurn,
			})
		}
	}

	var cards = make([]Card, 0)
	nPrevCardBase = 0
	for nCurrentByteIndex < nTotalCardBytes {
		var nCardCount = 0
		var nCardID = 0

		if !readSerializedCard(deckBytes, &nCurrentByteIndex, nTotalCardBytes, &nPrevCardBase, &nCardCount, &nCardID) {
			return Deck{}, errors.New("Could not read in non-hero card")
		}

		cards = append(cards, Card{
			Id: nCardID,
			Count: nCardCount,
		})
	}

	var name = ""
	if nCurrentByteIndex <= nTotalCardBytes {
		var bytes = deckBytes[len(deckBytes) - nStringLength:]
		name = string(bytes)
	}

	return Deck{
		Heroes: heroes,
		Cards: cards,
		Name: name,
	}, nil
}

func readSerializedCard(data []byte, indexStart *int, indexEnd int, nPrevCardBase *int, nOutCount *int, nOutCardID *int) bool {
	if *indexStart == indexEnd {
		return true
	}

	if *indexStart > indexEnd {
		return false
	}

	var nHeader = int(data[*indexStart])
	*indexStart++
	var bHasExtendedCount = (nHeader >> 6) == 0x03
	var nCardDelta = 0

	if !readVarEncodedUint32(nHeader, 5, data, indexStart, indexEnd, &nCardDelta) {
		return false
	}

	*nOutCardID = *nPrevCardBase + nCardDelta

	if bHasExtendedCount {
		if !readVarEncodedUint32(0, 0, data, indexStart, indexEnd, nOutCount) {
			return false
		}
	} else {
		*nOutCount = (nHeader >> 6) + 1
	}

	*nPrevCardBase = *nOutCardID
	return true
}

func readVarEncodedUint32(nBaseValue int, nBaseBits int, data []byte, indexStart *int, indexEnd int, outValue *int) bool {
	*outValue = 0
	var nDeltaShift = 0

	if nBaseBits == 0 || readBitsChunk(nBaseValue, nBaseBits, nDeltaShift, outValue) {
		nDeltaShift += nBaseBits

		for {
			if *indexStart > indexEnd {
				return false
			}

			if *indexStart == indexEnd {
				return true
			}

			var nNextByte = data[*indexStart]
			*indexStart++

			if !readBitsChunk(int(nNextByte), 7, nDeltaShift, outValue) {
				break
			}

			nDeltaShift += 7
		}
	}

	return true
}

func readBitsChunk(nChunk int, nNumBits int, nCurrShift int, nOutBits *int) bool {
	var nContinueBit = 1 << uint(nNumBits)
	var nNewBits = nChunk & (nContinueBit - 1)
	*nOutBits |= nNewBits << uint(nCurrShift)

	return (nChunk & nContinueBit) != 0
}
