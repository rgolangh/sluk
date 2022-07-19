package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/rgolangh/sluk/data"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	exact                   *bool
	printUnicodeValue       *bool
	printUnicodeDescription *bool
	//db                      = "/home/rgolan/Downloads/UCD/extracted/DerivedName.txt"
	dbFile *string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sluk [search-term]",
	Short: "sluk (simbol look up) will lookup a unicode simbol by name and print",
	Long: `Sluk will lookup a simbol in the unicode map extracted from unicode.org
and will print it the terminal. For example:
sluk white heavy check mark
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

func run(cmd *cobra.Command, args []string) {
	results := map[string]string{}
	searchTerm := strings.ToUpper(strings.Join(args, " "))

	var reader io.Reader
	if *dbFile != "" {
		reader, err := os.Open(*dbFile)
		cobra.CheckErr(err)
		defer reader.Close()
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
		if *exact {
			if v == searchTerm {
				results[k] = v
			}
		} else {
			if strings.Contains(v, searchTerm) {
				results[k] = v
			}
		}
	}

	for code, desc := range results {
		// pad the unicode value with zeros so it will be valid unicode hex value
		s := fmt.Sprintf("'\\U%08s'", code)
		unquote, err := strconv.Unquote(s)
		cobra.CheckErr(err)
		fmt.Printf("%s", unquote)
		if *printUnicodeValue {
			fmt.Printf("\t%s", s)
		}
		if *printUnicodeDescription {
			fmt.Printf("\t%s", desc)
		}
		fmt.Println()
	}
}

func init() {
	exact = rootCmd.Flags().BoolP("exact-match", "e", false, "Exact match the search term")
	printUnicodeValue = rootCmd.Flags().BoolP("print-unicode", "p", false, "Print the unicode value")
	printUnicodeDescription = rootCmd.Flags().BoolP("print-description", "d", false, "Print the unicode description")
	dbFile = rootCmd.Flags().StringP("db-file", "f", "", "A file containing extracted unicode mapping in the form of [unicode] ; [description]")
}
