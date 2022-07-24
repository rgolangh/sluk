package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rgolangh/sluk/data"
	"github.com/spf13/cobra"
)

var (
	verbose                 bool
	exact                   bool
	printUnicodeValue       bool
	printUnicodeDescription bool
	dbFile                  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sluk [search-term]",
	Short: "sluk (symbol look up) will lookup a unicode symbol by name and print",
	Long: `Sluk will lookup (fuzzy search) the symbol description in the unicode
map extracted from unicode.org and will print it the terminal. For example:
$ sluk heavy check mark
✔
✅
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a search term")
		}
		return nil
	},
	Run: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

type match struct {
	rank        int
	unicode     string
	description string
}

func run(cmd *cobra.Command, args []string) {
	var results []match
	searchTerm := strings.ToUpper(strings.Join(args, " "))

	var reader io.Reader
	if dbFile != "" {
		fmt.Printf("db file %v\n", dbFile)
		fd, err := os.Open(dbFile)
		cobra.CheckErr(err)
		defer fd.Close()
		reader = fd
	} else {
		reader = strings.NewReader(data.DB)
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		split := strings.Split(text, ";")
		if len(split) < 2 {
			cobra.CheckErr(fmt.Errorf("failed to parse the search term %s", searchTerm))
		}
		k, v := strings.TrimSpace(split[0]), strings.TrimSpace(split[1])
		if exact {
			if v == searchTerm {
				results = append(results, match{unicode: k, description: v})
			}
		} else {
			words := strings.Split(split[1], " ")
			matches := fuzzy.RankFind(searchTerm, words)
			sort.Sort(matches)

			closestMatch := matches[0]
			if len(matches) > 0 && closestMatch.Distance > -1 && closestMatch.Distance < 10 {
				results = append(results, match{rank: closestMatch.Distance, unicode: k, description: v})
			}
		}
	}

	sort.SliceStable(results, func(i, j int) bool { return results[i].rank < results[j].rank })

	for _, m := range results {
		// pad the unicode value with zeros so it will be valid unicode hex value
		s := fmt.Sprintf("'\\U%08s'", m.unicode)
		unquote, err := strconv.Unquote(s)
		cobra.CheckErr(err)
		fmt.Printf("%s", unquote)
		if printUnicodeValue {
			fmt.Printf("\t%s", s)
		}
		if printUnicodeDescription {
			fmt.Printf("\t%s", m.description)
		}
		fmt.Println()
		if verbose {
			fmt.Printf("%+v\n", m)
		}
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print verbose and debug info")
	rootCmd.Flags().BoolVarP(&exact, "exact-match", "e", false, "Exact match the search term")
	rootCmd.Flags().BoolVarP(&printUnicodeValue, "print-unicode", "p", false, "Print the unicode value")
	rootCmd.Flags().BoolVarP(&printUnicodeDescription, "print-description", "d", false, "Print the unicode description")
	rootCmd.Flags().StringVarP(&dbFile, "db-file", "f", "", "A file containing extracted unicode mapping in the form of [unicode] ; [description]")
}
