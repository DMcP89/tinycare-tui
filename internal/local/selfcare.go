// ChatGPT Prompts https://chat.openai.com/share/fd3771d8-0bde-460d-bb52-4300a81cdb50
package local

import (
	"fmt"
	"math/rand"
	"time"
)

var advices = []string{
	"Don't forget to drink enough water!",
	"Have you gotten outside today?",
	"Take a break and do something you enjoy.",
	"Get a good night's sleep for better well-being.",
	"Practice deep breathing or meditation.",
	"Engage in regular physical exercise.",
	"Spend time with loved ones or friends.",
}

var emoji = []string{
	"ᕕ(⌐■_■)ᕗ ♪♬",
	"༼ つ ◕_◕ ༽つ",
	"(◕ᴥ◕ʋ)",
	"(ﾉ◕ヮ◕)ﾉ*:・ﾟ✧",
	"ᕕ( ᐛ )ᕗ",
	"L(° O °L)",
}

var rng *rand.Rand

func init() {
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rng = rand.New(src)
}

func GetSelfCareAdvice() string {
	randomAdviceIndex := rng.Intn(len(advices))
	randomEmojiIndex := rng.Intn(len(emoji))
	return fmt.Sprintf("%s \n %s", advices[randomAdviceIndex], emoji[randomEmojiIndex])
}
