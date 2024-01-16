package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

type Section struct {
	Name          string    `json:"name,omitempty"`
	Childs        []Setting `json:"childs,omitempty"`
	SubSections   []Section `json:"subsection,omitempty"`
	ParentSection *Section  `json:"-"`
}

type Setting struct {
	Name           string `json:"name,omitempty"`
	Variable       string `json:"variable,omitempty"`
	Description    string `json:"description,omitempty"`
	TypeOfSetting  string `json:"type,omitempty"`
	DefaultSetting string `json:"default,omitempty"`
	Range          string `json:"range,omitempty"`
}

type entryPoint struct {
	url          string
	startAt      string
	skipSections []string
	getFirstOnly bool
	section      string
}

func getName(variable string) (name string) {
	name = strings.ReplaceAll(variable, "_", " ")
	name = strings.ReplaceAll(name, "col.", "Color ")
	t := []rune(strings.ToLower(name))
	t[0] = unicode.ToUpper(t[0])
	name = string(t)
	return
}

func formatDescription(description string) (descriptionFormated string) {
	t := []rune(strings.ToLower(description))
	t[0] = unicode.ToUpper(t[0])
	descriptionFormated = string(t)
	if !strings.HasSuffix(descriptionFormated, ".") {
		descriptionFormated += "."
	}
	return
}

func getRangeOrOption(description string) (rangeOrOption string) {
	re := regexp.MustCompile(`\[([^]]+)\]$`)
	matches := re.FindStringSubmatch(description)

	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

func parsePage(scanner *bufio.Scanner, startAt string, skipSections []string, getFirstOnly bool, name string, section *Section) {

	var afterSections bool = false
	var hint bool = false
	var inTable bool = false
	var currentSection *Section = section
	var depth int = 1
	var ntable int = 0

	if getFirstOnly {
		currentSection.SubSections = append(currentSection.SubSections, Section{Name: name, ParentSection: currentSection})
		currentSection = &currentSection.SubSections[len(currentSection.SubSections)-1]
	}

	for scanner.Scan() {
		line := scanner.Text()
		if !afterSections && strings.HasPrefix(line, startAt) {
			afterSections = true
		} else if !afterSections {
			continue
		}

		// Skip Hints
		if strings.HasPrefix(line, "{{") {
			hint = !hint
			continue
		} else if hint {
			continue
		}

		// Detect tables
		if strings.HasPrefix(line, "|---|") {
			continue
		} else if strings.HasPrefix(line, "|") {
			if !inTable {
				inTable = true
				continue
			}
		} else {
			if inTable {
				ntable++
			}
			inTable = false
		}

		if ntable == 1 && getFirstOnly {
			break
		}

		// Remove empty sections
		skip := false
		for _, skipSection := range skipSections {
			if strings.HasPrefix(line, skipSection) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		// Handle sections and subsections
		if strings.HasPrefix(line, "###") {
			if depth > 1 && currentSection.ParentSection != nil {
				currentSection = currentSection.ParentSection
			}
			sectionName := strings.TrimSpace(strings.TrimPrefix(line, "###"))
			currentSection.SubSections = append(currentSection.SubSections, Section{Name: sectionName, ParentSection: currentSection})
			currentSection = &currentSection.SubSections[len(currentSection.SubSections)-1]
			depth = 2
		} else if strings.HasPrefix(line, "##") {
			if depth > 1 && currentSection.ParentSection != nil {
				currentSection = currentSection.ParentSection
				currentSection = currentSection.ParentSection
			} else if currentSection.ParentSection != nil {
				currentSection = currentSection.ParentSection
			}

			sectionName := strings.TrimSpace(strings.TrimPrefix(line, "##"))
			currentSection.SubSections = append(currentSection.SubSections, Section{Name: sectionName, ParentSection: currentSection})
			currentSection = &currentSection.SubSections[len(currentSection.SubSections)-1]
			depth = 1
		}

		if !inTable {
			continue
		}

		temp := strings.Split(line, "|")
		setting := Setting{
			Name:           getName(strings.TrimSpace(temp[1])),
			Variable:       strings.TrimSpace(temp[1]),
			Description:    formatDescription(strings.TrimSpace(temp[2])),
			Range:          getRangeOrOption(strings.TrimSpace(temp[2])),
			TypeOfSetting:  strings.TrimSpace(temp[3]),
			DefaultSetting: strings.TrimSpace(temp[4]),
		}

		currentSection.Childs = append(currentSection.Childs, setting)
	}

}

func main() {

	var section Section = Section{Name: "Hyprland Wiki Variables"}
	entryPoints := []entryPoint{
		{
			url:          "https://raw.githubusercontent.com/hyprwm/hyprland-wiki/main/pages/Configuring/Variables.md",
			startAt:      "# Section",
			skipSections: []string{"## More", "## Per-device"},
		},
		{
			url:          "https://raw.githubusercontent.com/hyprwm/hyprland-wiki/main/pages/Configuring/Dwindle-Layout.md",
			startAt:      "# Config",
			getFirstOnly: true,
			section:      "dwindle",
		},
		{
			url:          "https://raw.githubusercontent.com/hyprwm/hyprland-wiki/main/pages/Configuring/Master-Layout.md",
			startAt:      "# Config",
			getFirstOnly: true,
			section:      "master",
		},
	}

	for _, entryPoint := range entryPoints {
		response, err := http.Get(entryPoint.url)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		scanner := bufio.NewScanner(response.Body)
		parsePage(scanner, entryPoint.startAt, entryPoint.skipSections, entryPoint.getFirstOnly, entryPoint.section, &section)

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}

	jsonData, err := json.Marshal(section)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonData))
}
