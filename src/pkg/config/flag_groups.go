package config

import (
	"fmt"
	"github.com/amit7itz/goset"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sort"
	"strings"
)

var validFlagCombinationsByCommand = make(map[*cobra.Command][][]string)

func MarkValidFlagCombinations(cmd *cobra.Command, favoriteFlagCombination []string, validFlagCombinations ...[]string) {
	// favoriteFlagCombination is always first in the slice.
	// if no flags are given, interactive mode will ask for favorite combination flags
	validFlagCombinationsByCommand[cmd] = append(validFlagCombinationsByCommand[cmd], favoriteFlagCombination)
	validFlagCombinationsByCommand[cmd] = append(validFlagCombinationsByCommand[cmd], validFlagCombinations...)
}

func ValidateFlagCombination(cmd *cobra.Command) error {
	if len(validFlagCombinationsByCommand[cmd]) == 0 {
		// no valid combinations provided
		return nil
	}
	// figure out what are all the flags, and convert the combinations to Sets
	allFlags := goset.NewSet[string]()
	flagCombinationSets := make([]*goset.Set[string], 0, len(validFlagCombinationsByCommand[cmd]))
	for _, flagCombination := range validFlagCombinationsByCommand[cmd] {
		set := goset.FromSlice(flagCombination)
		flagCombinationSets = append(flagCombinationSets, set)
		allFlags.Update(set)
	}

	// figure out what flags were used
	usedFlags := goset.NewSet[string]()
	for _, flag := range allFlags.Items() {
		if viper.IsSet(flag) {
			usedFlags.Add(flag)
		}
	}

	// find out what flags are missing
	largestIntersection := goset.NewSet[string]()
	missingFlagsOptions := make([]*goset.Set[string], 0)
	for _, flagCombinationSet := range flagCombinationSets {
		if usedFlags.Equal(flagCombinationSet) {
			// yay! a valid combination of flags was used
			return nil
		}
		if usedFlags.Difference(flagCombinationSet).Len() > 0 {
			// they used a flag that is not in this combination
			continue
		}
		intersection := flagCombinationSet.Intersection(usedFlags)
		if intersection.Len() > largestIntersection.Len() {
			largestIntersection = intersection
			missingFlagsOptions = []*goset.Set[string]{flagCombinationSet.Difference(intersection)}
		} else if intersection.Len() == largestIntersection.Len() {
			// we have the same flag intersection as a previous combination, so we add the missing flags as another option
			missingFlagsOptions = append(missingFlagsOptions, flagCombinationSet.Difference(intersection))
		}
	}
	if len(missingFlagsOptions) == 0 || missingFlagsOptions[0].IsEmpty() {
		// at this point we know it's not any of the valid combinations,
		// so if there are no missing flags, it means we have too much of them.
		return fmt.Errorf("invalid combination of flags, you should probably remove one or more flags")
	}
	return newMissingFlagsError(missingFlagsOptions)
}

func addDashes(flags []string) []string {
	return lo.Map(flags, func(flag string, _ int) string {
		return "--" + flag
	})
}

func newMissingFlagsError(missingFlagsOptions []*goset.Set[string]) error {
	errorString := "missing flags: "
	sortedMissingFlagOptions := make([][]string, 0, len(missingFlagsOptions))
	for i, missingFlagsOption := range missingFlagsOptions {
		sortedMissingFlagsOption := missingFlagsOption.Items()
		sort.Strings(sortedMissingFlagsOption)
		sortedMissingFlagOptions = append(sortedMissingFlagOptions, sortedMissingFlagsOption)
		if i > 0 {
			errorString += ", or "
		}
		errorString += strings.Join(addDashes(sortedMissingFlagsOption), " + ")
	}
	return MissingFlagsError{MissingFlagsOptions: sortedMissingFlagOptions, Err: errorString}
}

type MissingFlagsError struct {
	MissingFlagsOptions [][]string
	Err                 string
}

func (m MissingFlagsError) Error() string {
	return m.Err
}
