package security

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSha256AndWordList(t *testing.T) {
	testCases := []struct {
		str              string
		expectedHash     string
		expectedWordList string
	}{
		{
			str:              "test",
			expectedHash:     "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			expectedWordList: "quota letterhead stagnate inventive newborn disbelief klaxon glossary pupil combustion Trojan Orlando solo existence stagnate bifocals reform rebellion dropper bravado briefcase armistice miser Chicago stairway filament glucose bifocals ruffled upcoming allow antenna",
		},
		{
			str:              "",
			expectedHash:     "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			expectedWordList: "tissue phonetic snowslide December printer Wilmington befriend belowground pupil Wichita upshot retrieval prowler hemisphere sentence Capricorn brackish performance cranky tradition flytrap Norwegian playhouse disbelief regain Montana prowler bravado island enrollment select equipment",
		},
		{
			str:              "This is a longer test string",
			expectedHash:     "7ad5e82c3271d879e0515a1356ca1ad5d3c1ad0ec23fdbe4365c8b4b4e369b4e",
			expectedWordList: "keyboard specialist trauma Chicago checkup hideaway stormy inertia tapeworm enchanting enlist barbecue egghead revenue beehive specialist stapler recover ringbolt Atlantic snapshot customer suspense tradition Christmas fascinate obtuse disable drifter congregate puppy distortion",
		},
	}

	for _, testCase := range testCases {
		hash, wordList := GetSha256AndWordList(testCase.str)
		assert.Equal(t, testCase.expectedHash, hash, "Expected GetSha256AndWordList to return correct hash")
		assert.Equal(t, testCase.expectedWordList, wordList, "Expected GetSha256AndWordList to return correct word list")
	}
}

func TestBytesToPgpWords(t *testing.T) {
	testCases := []struct {
		bytes         []byte
		expectedWords []string
	}{
		{
			bytes:         []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			expectedWords: []string{"aardvark", "adviser", "accrue", "aggregate", "adrift", "almighty", "afflict", "amusement", "aimless", "applicant"},
		},
		{
			bytes:         []byte{255, 254, 253, 252, 251, 250, 249, 248, 247, 246},
			expectedWords: []string{"Zulu", "yesteryear", "willow", "Wilmington", "watchword", "whimsical", "waffle", "warranty", "virus", "vocalist"},
		},
		{
			bytes:         []byte{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
			expectedWords: []string{"allow", "belowground", "berserk", "cellulose", "checkup", "crossover", "cubic", "embezzle", "enlist", "getaway"},
		},
	}

	for _, testCase := range testCases {
		words := bytesToPgpWords(testCase.bytes)
		assert.Equal(t, testCase.expectedWords, words, "Expected bytesToPgpWords to return correct words")
	}
}
