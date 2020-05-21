package platforms

import (
	"fmt"
	"github.com/s0nerik/goloc/goloc"
	"github.com/s0nerik/goloc/goloc/re"
	"github.com/s0nerik/goloc/registry"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func init() {
	registry.RegisterPlatform(&flutter{})
}

type flutter struct{}

func (flutter) Names() []string {
	return []string{
		"flutter",
		"Flutter",
	}
}

func (flutter) LocalizationFilePath(lang goloc.Lang, resDir goloc.ResDir) string {
	return filepath.Join(resDir, fmt.Sprintf("localizations_%s.g.dart", lang))
}

func (flutter) Header(args *goloc.HeaderArgs) string {
	return fmt.Sprintf(`// DO NOT EDIT. This is code generated via https://github.com/s0nerik/goloc
// This is a library that provides messages for a en locale. All the
// messages from the main program should be duplicated here with the same
// function name.

// ignore_for_file: annotate_overrides, prefer_single_quotes, lines_longer_than_80_chars, non_constant_identifier_names, avoid_escaping_inner_quotes

part of 'localizations.dart';

class AppLocalizations%s implements AppLocalizations {
  final AppLocalizations fallback;

  AppLocalizations%s(this.fallback);

`, strings.Title(args.Lang), strings.Title(args.Lang))
}

func (flutter) LocalizedString(args *goloc.LocalizedStringArgs) string {
	if len(args.FormatArgs) > 0 {
		fArgs := buildFormatArgsList(args.FormatArgs, nil)
		return fmt.Sprintf("  String %s(%s) => sprintf(\"%s\", [%s]);\n", args.Key, fArgs, args.Value, fArgs)
	} else {
		return fmt.Sprintf("  String get %s => \"%s\";\n", args.Key, args.Value)
	}
}

func (flutter) FallbackString(args *goloc.LocalizedStringArgs) string {
	if len(args.FormatArgs) > 0 {
		fArgs := buildFormatArgsList(args.FormatArgs, nil)
		return fmt.Sprintf("  String %s(%s) => fallback?.%s(%s);\n", args.Key, fArgs, args.Key, fArgs)
	} else {
		return fmt.Sprintf("  String get %s => fallback?.%s;\n", args.Key, args.Key)
	}
}

func (flutter) FormatString(args *goloc.FormatStringArgs) string {
	return fmt.Sprintf("%%%s", args.Format)
}

func (flutter) Footer(args *goloc.FooterArgs) string {
	return "}"
}

func (flutter) ValidateFormat(format string) error {
	return nil
}

func (flutter) ReplacementChars() map[string]string {
	return map[string]string{
		"\n": `\n`,
		"\t": `\t`,
		`"`:  `\"`,
		`\`:  `\\`,
	}
}

func LocalizationsContent(args goloc.PreprocessArgs) string {
	contentFmt := `// DO NOT EDIT. This is code generated via https://github.com/s0nerik/goloc
// This is a library that provides messages for a en locale. All the
// messages from the main program should be duplicated here with the same
// function name.

// ignore_for_file: annotate_overrides, prefer_single_quotes, lines_longer_than_80_chars, non_constant_identifier_names, avoid_escaping_inner_quotes

import 'dart:async';
import 'dart:ui';

import 'package:flutter/widgets.dart';
import 'package:sprintf/sprintf.dart';

%s
abstract class AppLocalizations {
  static AppLocalizations of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations);
  }

%s}

class AppLocalizationsDelegate extends LocalizationsDelegate<AppLocalizations> {
  static const supportedLanguages = %s;

  @override
  bool isSupported(Locale locale) {
    return supportedLanguages.contains(locale.languageCode);
  }

  @override
  Future<AppLocalizations> load(Locale locale) {
    switch (locale.languageCode) {
%s    }
  }

  @override
  bool shouldReload(LocalizationsDelegate<AppLocalizations> old) {
    return false;
  }
}
`

	// Parts/supported locales
	var partsBuilder strings.Builder
	for _, loc := range args.Localizations.Locales() {
		partsBuilder.WriteString(fmt.Sprintf("part 'localizations_%s.g.dart';\n", loc))
	}

	// Localized strings
	var locBuilder strings.Builder
	for _, key := range args.Localizations.SortedKeys() {
		fArgs := args.FormatArgs[key]
		if len(fArgs) <= 0 {
			str := fmt.Sprintf("  String get %s;\n", key)
			locBuilder.WriteString(str)
		} else {
			typedArgsStr := buildFormatArgsList(fArgs, args.Formats)
			str := fmt.Sprintf("  String %s(%s);\n", key, typedArgsStr)
			locBuilder.WriteString(str)
		}
	}

	// Supported locales
	supportedLocales := `["` + strings.Join(args.Localizations.Locales(), `", "`) + `"]`

	// AppLocalizations loading
	var loadBuilder strings.Builder
	for _, loc := range args.Localizations.Locales() {
		fallback := "null"
		if loc != args.DefaultLocalization {
			fallback = fmt.Sprintf("AppLocalizations%s(null)", strings.Title(args.DefaultLocalization))
		}
		loadBuilder.WriteString(fmt.Sprintf(`      case '%s':
        return Future.value(AppLocalizations%s(%s));
`, loc, strings.Title(loc), fallback))
	}
	loadBuilder.WriteString(fmt.Sprintf(`      default:
        return Future.value(AppLocalizations%s(null));
`, strings.Title(args.DefaultLocalization)))

	return fmt.Sprintf(contentFmt, partsBuilder.String(), locBuilder.String(), supportedLocales, loadBuilder.String())
}

// buildFormatArgsList returns a ready-to-use list of format arguments for Dart.
// If `formats` are specified - returns a list of typed elements, otherwise - untyped
func buildFormatArgsList(fArgs []goloc.FormatKey, formats goloc.Formats) string {
	var argsListBuilder strings.Builder
	for i, fKey := range fArgs {
		var argNameBuilder strings.Builder

		f := formats[fKey]
		matches := re.SprintfRegexp().FindStringSubmatch("%" + f)
		if len(matches) >= 5 {
			valueType := matches[5]
			switch valueType {
			case "s":
				argNameBuilder.WriteString("String ")
			case "i", "d", "x", "X", "o", "O":
				argNameBuilder.WriteString("int ")
			case "e", "E", "f", "F", "g", "G":
				argNameBuilder.WriteString("double ")
			}
		}

		argNameBuilder.WriteString(fmt.Sprintf("arg%v", i))
		if i < len(fArgs)-1 {
			argNameBuilder.WriteString(", ")
		}
		argsListBuilder.WriteString(argNameBuilder.String())
	}
	return argsListBuilder.String()
}

func (flutter) Preprocess(args goloc.PreprocessArgs) (err error) {
	locFileName := filepath.Join(args.ResDir, "localizations.dart")
	err = ioutil.WriteFile(locFileName, []byte(LocalizationsContent(args)), 0644)
	return
}
