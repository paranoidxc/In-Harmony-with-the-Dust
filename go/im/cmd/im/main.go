package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"im/cmd/internal/sourceflags"
	"im/pkg/ime"
)

func main() {
	fcitxPath := flag.String("fcitx", "", "override path to fcitx pinyin dictionary; empty uses ime default")
	sogouPath := flag.String("sogou", "", "override path to sogou pinyin dictionary; empty uses ime default")
	limit := flag.Int("n", 10, "max candidates to print")
	var sources sourceflags.MultiValue
	flag.Var(&sources, "source", "dictionary source spec, repeatable: fcitx=path, sogou-syllables=path, sogou=path")
	flag.Parse()

	loadStart := time.Now()
	var (
		engine *ime.Engine
		err    error
	)
	if len(sources) == 0 && *fcitxPath == "" && *sogouPath == "" {
		engine, err = ime.LoadDefaultEngine()
	} else {
		config, buildErr := sourceflags.BuildConfig(sources, *fcitxPath, *sogouPath)
		if buildErr != nil {
			log.Fatalf("build source config: %v", buildErr)
		}
		engine, err = ime.LoadEngineWithConfig(config)
	}
	if err != nil {
		log.Fatalf("load dictionary: %v", err)
	}
	fmt.Printf("loaded dictionaries in %s\n", time.Since(loadStart).Round(time.Millisecond))
	fmt.Println("input pinyin or initials, press Enter to search, ':quit' to exit")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == ":quit" || line == ":q" || line == "quit" || line == "exit" {
			return
		}

		result := engine.Search(line, *limit)
		if result.Log.Normalized == "" {
			fmt.Println("empty query after normalization")
			continue
		}
		if len(result.Items) == 0 {
			fmt.Printf("no matches for %q\n", line)
			printLog(result.Log)
			continue
		}

		for i, item := range result.Items {
			extra := item.Key
			if item.Pattern != "" {
				extra = item.Pattern
			}
			fmt.Printf("%2d. %s\t[%s:%s] match=%s\n", i+1, item.Word, item.Source, extra, item.Match)
		}
		printLog(result.Log)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("read stdin: %v", err)
	}
}

func printLog(match ime.MatchLog) {
	segments := "-"
	if len(match.Segmentations) > 0 {
		segments = strings.Join(match.Segmentations, ",")
	}
	fmt.Printf(
		"log query=%q normalized=%q segments=%q exact_full=%d exact_initial=%d derived_initial=%d combined=%d tail_composed=%d tail_fallback=%d prefix_full=%d prefix_initial=%d returned=%d elapsed=%s\n",
		match.Query,
		match.Normalized,
		segments,
		match.ExactFullMatches,
		match.ExactInitialMatches,
		match.DerivedInitialMatches,
		match.CombinedMatches,
		match.TailComposedMatches,
		match.TailFallbackMatches,
		match.PrefixFullMatches,
		match.PrefixInitialMatches,
		match.Returned,
		match.Elapsed.Round(time.Microsecond),
	)
}
