package namegen

import (
	"fmt"
	"math/rand"
	"time"
)

var adjectives = []string{
	"amber", "azure", "bright", "calm", "clear", "crisp", "dark", "deep",
	"dense", "eager", "early", "fair", "fast", "firm", "fleet", "fond",
	"free", "fresh", "glad", "gold", "good", "grand", "great", "green",
	"grey", "high", "keen", "kind", "large", "late", "lean", "light",
	"lofty", "long", "loud", "lunar", "mild", "mist", "neat", "noble",
	"north", "opal", "open", "pale", "plain", "prime", "proud", "pure",
	"quick", "quiet", "rapid", "rich", "rosy", "round", "royal", "rust",
	"sage", "sharp", "sheer", "silk", "slim", "smart", "soft", "solar",
	"solid", "stark", "still", "stone", "storm", "stout", "sure", "swift",
	"tall", "tame", "teal", "thin", "tidy", "true", "vast", "warm",
	"white", "wide", "wild", "wise", "bold", "brave", "cool",
}

var nouns = []string{
	"arrow", "atlas", "bay", "beacon", "bear", "bird", "blade", "blaze",
	"bloom", "brook", "canyon", "cedar", "cloud", "coast", "coral", "crane",
	"creek", "delta", "drift", "dune", "eagle", "ember", "falcon", "fern",
	"field", "fjord", "flame", "flash", "fleet", "forest", "frost", "gale",
	"gate", "glade", "glen", "grove", "gust", "haven", "hawk", "heath",
	"hill", "hollow", "horizon", "island", "jade", "lake", "lark", "leaf",
	"ledge", "light", "marsh", "meadow", "mesa", "mist", "moon", "moss",
	"mountain", "oak", "ocean", "peak", "pine", "plain", "pond", "quartz",
	"rain", "rapid", "raven", "reef", "ridge", "river", "rock", "sage",
	"shore", "sierra", "sky", "slate", "snow", "sol", "spark", "spring",
	"star", "stone", "storm", "stream", "summit", "surf", "tide", "timber",
	"trail", "vale", "valley", "wave", "wind", "wolf", "wood",
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// Generate returns a random "adjective-noun" name prefixed with the repo name.
// e.g. "myrepo-bright-falcon"
func Generate(repoName string) string {
	adj := adjectives[rng.Intn(len(adjectives))]
	noun := nouns[rng.Intn(len(nouns))]
	return fmt.Sprintf("%s-%s-%s", repoName, adj, noun)
}
