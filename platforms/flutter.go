package platforms

import (
	"fmt"
	"github.com/s0nerik/goloc/goloc"
	"github.com/s0nerik/goloc/platforms/registry"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func init() {
	registry.Register(&flutter{})
}

type flutter struct{}

func (flutter) Names() []string {
	return []string{
		"flutter",
		"Flutter",
	}
}

func (flutter) LocalizationFilePath(lang goloc.Lang, resDir goloc.ResDir) string {
	return filepath.Join(resDir, fmt.Sprintf("messages_%s.dart", lang))
}

func (flutter) Header(args *goloc.HeaderArgs) string {
	return fmt.Sprintf(`// DO NOT EDIT. This is code generated via https://github.com/s0nerik/goloc
// This is a library that provides messages for a en locale. All the
// messages from the main program should be duplicated here with the same
// function name.

import 'package:intl/intl.dart';
import 'package:intl/message_lookup_by_library.dart';
import 'package:sprintf/sprintf.dart';

// ignore: unnecessary_new
final messages = new MessageLookup();

// ignore: unused_element
final _keepAnalysisHappy = Intl.defaultLocale;

// ignore: non_constant_identifier_names
typedef MessageIfAbsent(String message_str, List args);

class MessageLookup extends MessageLookupByLibrary {
  get localeName => '%s';

  final messages = _notInlinedMessages(_notInlinedMessages);
  static _notInlinedMessages(_) => <String, Function>{
`, args.Lang)
}

func (flutter) LocalizedString(args *goloc.LocalizedStringArgs) string {
	if len(args.FormatArgs) > 0 {
		fArgs := buildFormatArgsList(args.FormatArgs)
		return fmt.Sprintf("        \"%s\": (%s) => sprintf(\"%s\", [%s]),\n", args.Key, fArgs, args.Value, fArgs)
	} else {
		return fmt.Sprintf("        \"%s\": MessageLookupByLibrary.simpleMessage(\"%s\"),\n", args.Key, args.Value)
	}
}

func (flutter) FormatString(args *goloc.FormatStringArgs) string {
	return fmt.Sprintf("%%%s", args.Format)
}

func (flutter) Footer(args *goloc.FooterArgs) string {
	return "      };\n}\n"
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

func PostprocessAllLocalizationsContent(args goloc.PostprocessArgs) string {
	locales := args.Localizations.Locales()
	defaultLocale := args.DefaultLocalization

	var builder strings.Builder

	// Header
	builder.WriteString(`// DO NOT EDIT. This is code generated via https://github.com/s0nerik/goloc
// This is a library that looks up messages for specific locales by
// delegating to the appropriate library.

import 'dart:async';

import 'package:intl/intl.dart';
import 'package:intl/message_lookup_by_library.dart';
import 'package:intl/src/intl_helpers.dart';

`)

	// Per-locale messages imports
	for _, loc := range locales {
		builder.WriteString(fmt.Sprintf("import 'messages_%s.dart' as messages_%s;\n", loc, loc))
	}

	builder.WriteString("// ignore: implementation_imports\n")

	// Deferred library initialization
	builder.WriteString(`
typedef Future<dynamic> LibraryLoader();
Map<String, LibraryLoader> _deferredLibraries = {
`)
	for _, loc := range locales {
		line := fmt.Sprintf(`// ignore: unnecessary_new
  '%s': () => new Future.value(null),
`, loc)
		builder.WriteString(line)
	}
	builder.WriteString("};\n")

	// _findExact implementation
	builder.WriteString(`
MessageLookupByLibrary _findExact(localeName) {
  switch (localeName) {
`)
	for _, loc := range locales {
		line := fmt.Sprintf(`    case '%s':
      return messages_%s.messages;
`, loc, loc)
		builder.WriteString(line)
	}
	builder.WriteString(fmt.Sprintf(`    default:
      return messages_%s.messages;`, defaultLocale))

	// Footer
	builder.WriteString(`
  }
}

/// User programs should call this before using [localeName] for messages.
Future<bool> initializeMessages(String localeName) async {
  var availableLocale =
      Intl.verifiedLocale(localeName, (locale) => _deferredLibraries[locale] != null, onFailure: (_) => null);
  if (availableLocale == null) {
    // ignore: unnecessary_new
    return new Future.value(false);
  }
  var lib = _deferredLibraries[availableLocale];
  // ignore: unnecessary_new
  await (lib == null ? new Future.value(false) : lib());
  // ignore: unnecessary_new
  initializeInternalMessageLookup(() => new CompositeMessageLookup());
  messageLookup.addLocale(availableLocale, _findGeneratedMessagesFor);
  // ignore: unnecessary_new
  return new Future.value(true);
}

bool _messagesExistFor(String locale) {
  try {
    return _findExact(locale) != null;
  } catch (e) {
    return false;
  }
}

MessageLookupByLibrary _findGeneratedMessagesFor(locale) {
  var actualLocale = Intl.verifiedLocale(locale, _messagesExistFor, onFailure: (_) => null);
  if (actualLocale == null) return null;
  return _findExact(actualLocale);
}
`)

	return builder.String()
}

func PostprocessLocalizationsContent(args goloc.PostprocessArgs) string {
	var builder strings.Builder

	// Header
	builder.WriteString(`// DO NOT EDIT. This is code generated via https://github.com/s0nerik/goloc
// This is a library that provides messages for a en locale. All the
// messages from the main program should be duplicated here with the same
// function name.

import 'dart:async';
import 'dart:ui';

import 'package:flutter/widgets.dart';
import 'package:intl/intl.dart';

import 'messages_all.dart';

class AppLocalizations {
  static Future<AppLocalizations> load(Locale locale) {
    final String name = locale.countryCode == null ? locale.languageCode : locale.toString();
    final String localeName = Intl.canonicalizedLocale(name);

    return initializeMessages(localeName).then((bool _) {
      Intl.defaultLocale = localeName;
      return new AppLocalizations();
    });
  }

  static AppLocalizations of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations);
  }

`)

	// Localized strings
	for key := range args.Localizations {
		fArgs := args.FormatArgs[key]
		if len(fArgs) <= 0 {
			str := fmt.Sprintf("  String get %s => Intl.message('', name: '%s');\n", key, key)
			builder.WriteString(str)
		} else {
			argsStr := buildFormatArgsList(fArgs)
			str := fmt.Sprintf("  String %s(%s) => Intl.message('', name: '%s', args: [%s]);\n", key, argsStr, key, argsStr)
			builder.WriteString(str)
		}
	}

	// Footer
	builder.WriteString(`
}

class AppLocalizationsDelegate extends LocalizationsDelegate<AppLocalizations> {
  final List<String> supportedLocales;

  const AppLocalizationsDelegate(this.supportedLocales);

  @override
  bool isSupported(Locale locale) {
    return supportedLocales.contains(locale.languageCode);
  }

  @override
  Future<AppLocalizations> load(Locale locale) {
    return AppLocalizations.load(locale);
  }

  @override
  bool shouldReload(LocalizationsDelegate<AppLocalizations> old) {
    return false;
  }
}
`)

	return builder.String()
}

func buildFormatArgsList(fArgs []string) string {
	var argsListBuilder strings.Builder
	for i := range fArgs {
		var argName string
		if i < len(fArgs)-1 {
			argName = fmt.Sprintf("arg%v, ", i)
		} else {
			argName = fmt.Sprintf("arg%v", i)
		}
		argsListBuilder.WriteString(argName)
	}
	return argsListBuilder.String()
}

func (flutter) Postprocess(args goloc.PostprocessArgs) (err error) {
	allMessagesFileName := filepath.Join(args.ResDir, "messages_all.dart")
	err = ioutil.WriteFile(allMessagesFileName, []byte(PostprocessAllLocalizationsContent(args)), 0644)
	if err != nil {
		return
	}
	locFileName := filepath.Join(args.ResDir, "localizations.dart")
	err = ioutil.WriteFile(locFileName, []byte(PostprocessLocalizationsContent(args)), 0644)
	return
}
