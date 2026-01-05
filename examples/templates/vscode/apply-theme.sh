#!/bin/bash
# VS Code Live Theme - applies colors without restart
# Requires: jq

settings_file="$HOME/.config/Code/User/settings.json"

# Require jq for safe JSON manipulation
if ! command -v jq >/dev/null 2>&1; then
    echo "Error: jq required for VS Code theming"
    exit 1
fi

# Ensure settings directory and file exist
mkdir -p "$(dirname "$settings_file")"
[[ ! -f "$settings_file" ]] && echo '{}' > "$settings_file"

# Validate existing JSON (backup and recreate if malformed)
if ! jq empty "$settings_file" 2>/dev/null; then
    cp "$settings_file" "$settings_file.backup"
    echo '{}' > "$settings_file"
    echo "Warning: Fixed malformed settings.json (backup created)"
fi

# Colors from peachy template variables (using .strip to remove # prefix)
black="{black.strip}"
red="{red.strip}"
green="{green.strip}"
yellow="{yellow.strip}"
blue="{blue.strip}"
magenta="{magenta.strip}"
cyan="{cyan.strip}"
white="{white.strip}"
bright_black="{bright_black.strip}"
bright_red="{bright_red.strip}"
bright_green="{bright_green.strip}"
bright_yellow="{bright_yellow.strip}"
bright_blue="{bright_blue.strip}"
bright_magenta="{bright_magenta.strip}"
bright_cyan="{bright_cyan.strip}"
bright_white="{bright_white.strip}"
foreground="{foreground.strip}"
background="{background.strip}"

# Generate UI colors JSON
ui_colors=$(cat <<EOF
{
    "foreground": "#${white}",
    "disabledForeground": "#${bright_black}",
    "widget.shadow": "#${black}",
    "selection.background": "#${blue}",
    "descriptionForeground": "#${bright_black}",
    "errorForeground": "#${red}",
    "icon.foreground": "#${bright_black}",

    "textBlockQuote.background": "#${black}",
    "textBlockQuote.border": "#${blue}",
    "textCodeBlock.background": "#${black}",
    "textLink.activeForeground": "#${cyan}",
    "textLink.foreground": "#${blue}",
    "textPreformat.foreground": "#${blue}",

    "toolbar.hoverBackground": "#${black}",
    "toolbar.activeBackground": "#${bright_black}",

    "button.background": "#${blue}",
    "button.foreground": "#${bright_white}",
    "button.hoverBackground": "#${bright_black}",
    "button.secondaryForeground": "#${bright_white}",
    "button.secondaryBackground": "#${magenta}",
    "button.secondaryHoverBackground": "#${bright_black}",
    "checkbox.background": "#${black}",
    "checkbox.foreground": "#${white}",

    "dropdown.background": "#${black}",
    "dropdown.listBackground": "#${black}",
    "dropdown.foreground": "#${white}",

    "input.background": "#${black}",
    "input.foreground": "#${white}",
    "input.placeholderForeground": "#${bright_black}",
    "inputOption.activeBackground": "#${black}",
    "inputOption.activeBorder": "#${bright_red}",
    "inputOption.activeForeground": "#${white}",
    "inputValidation.errorBackground": "#${red}",
    "inputValidation.errorForeground": "#${white}",
    "inputValidation.errorBorder": "#${red}",
    "inputValidation.infoBackground": "#${blue}",
    "inputValidation.infoForeground": "#${white}",
    "inputValidation.infoBorder": "#${blue}",
    "inputValidation.warningBackground": "#${yellow}",
    "inputValidation.warningForeground": "#${white}",
    "inputValidation.warningBorder": "#${yellow}",

    "scrollbar.shadow": "#${black}",
    "scrollbarSlider.activeBackground": "#${bright_black}6f",
    "scrollbarSlider.background": "#${black}6f",
    "scrollbarSlider.hoverBackground": "#${bright_black}6f",

    "badge.background": "#${black}",
    "badge.foreground": "#${white}",

    "progressBar.background": "#${bright_black}",

    "list.activeSelectionBackground": "#${black}",
    "list.activeSelectionForeground": "#${white}",
    "list.dropBackground": "#${bright_white}",
    "list.focusBackground": "#${black}",
    "list.focusForeground": "#${white}",
    "list.highlightForeground": "#${bright_white}",
    "list.hoverBackground": "#${bright_black}",
    "list.hoverForeground": "#${white}",
    "list.inactiveSelectionBackground": "#${black}",
    "list.inactiveSelectionForeground": "#${white}",
    "list.inactiveFocusBackground": "#${black}",
    "list.invalidItemForeground": "#${red}",
    "list.errorForeground": "#${red}",
    "list.warningForeground": "#${yellow}",
    "listFilterWidget.background": "#${black}",
    "listFilterWidget.noMatchesOutline": "#${red}",
    "list.filterMatchBackground": "#${black}",
    "tree.indentGuidesStroke": "#${white}",

    "activityBar.background": "#${black}",
    "activityBar.foreground": "#${white}",
    "activityBar.inactiveForeground": "#${bright_black}",
    "activityBarBadge.background": "#${blue}",
    "activityBarBadge.foreground": "#${bright_white}",
    "activityBar.activeBackground": "#${black}",

    "sideBar.background": "#${black}",
    "sideBar.foreground": "#${white}",
    "sideBar.dropBackground": "#${black}6f",
    "sideBarTitle.foreground": "#${white}",
    "sideBarSectionHeader.background": "#${bright_black}",
    "sideBarSectionHeader.foreground": "#${white}",

    "minimap.findMatchHighlight": "#${yellow}6f",
    "minimap.selectionHighlight": "#${black}6f",
    "minimap.errorHighlight": "#${red}",
    "minimap.warningHighlight": "#${yellow}",
    "minimap.background": "#${black}",
    "minimap.selectionOccurrenceHighlight": "#${bright_black}",
    "minimapGutter.addedBackground": "#${green}",
    "minimapGutter.modifiedBackground": "#${magenta}",
    "minimapGutter.deletedBackground": "#${red}",

    "editorGroup.dropBackground": "#${black}6f",
    "editorGroupHeader.noTabsBackground": "#${black}",
    "editorGroupHeader.tabsBackground": "#${black}",
    "editorGroup.emptyBackground": "#${black}",
    "editorGroup.dropIntoPromptForeground": "#${bright_white}",
    "editorGroup.dropIntoPromptBackground": "#${black}",
    "tab.activeBackground": "#${black}",
    "tab.unfocusedActiveBackground": "#${black}",
    "tab.activeForeground": "#${white}",
    "tab.inactiveBackground": "#${black}",
    "tab.inactiveForeground": "#${bright_black}",
    "tab.unfocusedActiveForeground": "#${bright_black}",
    "tab.unfocusedInactiveForeground": "#${bright_black}",
    "tab.hoverBackground": "#${black}",
    "tab.unfocusedHoverBackground": "#${black}",
    "tab.activeModifiedBorder": "#${blue}",
    "tab.inactiveModifiedBorder": "#${blue}",
    "tab.unfocusedActiveModifiedBorder": "#${blue}",
    "tab.unfocusedInactiveModifiedBorder": "#${blue}",
    "editorPane.background": "#${black}",

    "editor.background": "#${black}",
    "editor.foreground": "#${white}",
    "editorLineNumber.foreground": "#${bright_black}",
    "editorLineNumber.activeForeground": "#${bright_black}",
    "editorCursor.foreground": "#${white}",
    "editor.selectionBackground": "#${black}",
    "editor.inactiveSelectionBackground": "#${black}",
    "editor.selectionHighlightBackground": "#${black}",
    "editor.wordHighlightBackground": "#${black}6f",
    "editor.wordHighlightStrongBackground": "#${bright_black}6f",
    "editor.findMatchBackground": "#${yellow}6f",
    "editor.findMatchHighlightBackground": "#${bright_red}6f",
    "editor.findRangeHighlightBackground": "#${black}6f",
    "searchEditor.findMatchBackground": "#${yellow}6f",
    "editor.hoverHighlightBackground": "#${black}6f",
    "editor.lineHighlightBackground": "#${black}",
    "editorLink.activeForeground": "#${blue}",
    "editor.rangeHighlightBackground": "#${black}6f",
    "editorWhitespace.foreground": "#${bright_black}",
    "editorIndentGuide.background1": "#${bright_black}",
    "editorIndentGuide.activeBackground1": "#${bright_black}",
    "editorInlayHint.background": "#${black}",
    "editorInlayHint.foreground": "#${white}",
    "editorInlayHint.typeBackground": "#${black}",
    "editorInlayHint.typeForeground": "#${white}",
    "editorInlayHint.parameterBackground": "#${black}",
    "editorInlayHint.parameterForeground": "#${white}",
    "editorRuler.foreground": "#${bright_black}",

    "editorCodeLens.foreground": "#${black}",
    "editorLightBulb.foreground": "#${yellow}",
    "editorLightBulbAutoFix.foreground": "#${blue}",
    "editorBracketMatch.background": "#${black}",

    "editorBracketHighlight.foreground1": "#${red}",
    "editorBracketHighlight.foreground2": "#${bright_red}",
    "editorBracketHighlight.foreground3": "#${yellow}",
    "editorBracketHighlight.foreground4": "#${green}",
    "editorBracketHighlight.foreground5": "#${blue}",
    "editorBracketHighlight.foreground6": "#${magenta}",
    "editorBracketHighlight.unexpectedBracket.foreground": "#${bright_magenta}",
    "editorOverviewRuler.findMatchForeground": "#${yellow}6f",
    "editorOverviewRuler.rangeHighlightForeground": "#${bright_black}6f",
    "editorOverviewRuler.selectionHighlightForeground": "#${black}6f",
    "editorOverviewRuler.wordHighlightForeground": "#${bright_white}6f",
    "editorOverviewRuler.wordHighlightStrongForeground": "#${blue}6f",
    "editorOverviewRuler.modifiedForeground": "#${magenta}",
    "editorOverviewRuler.addedForeground": "#${green}",
    "editorOverviewRuler.deletedForeground": "#${red}",
    "editorOverviewRuler.errorForeground": "#${red}",
    "editorOverviewRuler.warningForeground": "#${yellow}",
    "editorOverviewRuler.infoForeground": "#${cyan}",
    "editorOverviewRuler.bracketMatchForeground": "#${bright_white}",

    "editorError.foreground": "#${red}",
    "editorWarning.foreground": "#${yellow}",
    "editorInfo.foreground": "#${cyan}",
    "editorHint.foreground": "#${blue}",
    "problemsErrorIcon.foreground": "#${red}",
    "problemsWarningIcon.foreground": "#${yellow}",
    "problemsInfoIcon.foreground": "#${cyan}",

    "editorGutter.background": "#${black}",
    "editorGutter.modifiedBackground": "#${magenta}",
    "editorGutter.addedBackground": "#${green}",
    "editorGutter.deletedBackground": "#${red}",
    "editorGutter.commentRangeForeground": "#${bright_black}",
    "editorGutter.foldingControlForeground": "#${white}",

    "diffEditor.insertedTextBackground": "#${green}20",
    "diffEditor.removedTextBackground": "#${red}20",
    "diffEditor.diagonalFill": "#${black}",

    "editorWidget.foreground": "#${white}",
    "editorWidget.background": "#${black}",
    "editorSuggestWidget.background": "#${black}",
    "editorSuggestWidget.foreground": "#${white}",
    "editorSuggestWidget.focusHighlightForeground": "#${bright_white}",
    "editorSuggestWidget.highlightForeground": "#${blue}",
    "editorSuggestWidget.selectedBackground": "#${black}",
    "editorSuggestWidget.selectedForeground": "#${bright_white}",
    "editorHoverWidget.foreground": "#${white}",
    "editorHoverWidget.background": "#${black}",
    "debugExceptionWidget.background": "#${black}",
    "editorMarkerNavigation.background": "#${black}",
    "editorMarkerNavigationError.background": "#${red}",
    "editorMarkerNavigationWarning.background": "#${yellow}",
    "editorMarkerNavigationInfo.background": "#${blue}",
    "editorMarkerNavigationError.headerBackground": "#${red}20",
    "editorMarkerNavigationWarning.headerBackground": "#${yellow}20",
    "editorMarkerNavigationInfo.headerBackground": "#${cyan}20",

    "peekViewEditor.background": "#${black}",
    "peekViewEditorGutter.background": "#${black}",
    "peekViewEditor.matchHighlightBackground": "#${bright_red}6f",
    "peekViewResult.background": "#${black}",
    "peekViewResult.fileForeground": "#${white}",
    "peekViewResult.lineForeground": "#${bright_black}",
    "peekViewResult.matchHighlightBackground": "#${bright_red}6f",
    "peekViewResult.selectionBackground": "#${black}",
    "peekViewResult.selectionForeground": "#${white}",
    "peekViewTitle.background": "#${black}",
    "peekViewTitleDescription.foreground": "#${bright_black}",
    "peekViewTitleLabel.foreground": "#${white}",

    "merge.currentContentBackground": "#${blue}40",
    "merge.currentHeaderBackground": "#${blue}40",
    "merge.incomingContentBackground": "#${green}60",
    "merge.incomingHeaderBackground": "#${green}60",
    "editorOverviewRuler.currentContentForeground": "#${blue}",
    "editorOverviewRuler.incomingContentForeground": "#${green}",
    "editorOverviewRuler.commonContentForeground": "#${bright_magenta}",

    "panel.background": "#${black}",
    "panel.dropBorder": "#${black}6f",
    "panelTitle.activeForeground": "#${white}",
    "panelTitle.inactiveForeground": "#${bright_black}",

    "statusBar.background": "#${blue}",
    "statusBar.foreground": "#${bright_white}",
    "statusBar.debuggingBackground": "#${bright_red}",
    "statusBar.debuggingForeground": "#${bright_white}",
    "statusBar.noFolderBackground": "#${magenta}",
    "statusBar.noFolderForeground": "#${bright_white}",
    "statusBarItem.activeBackground": "#${bright_black}",
    "statusBarItem.hoverBackground": "#${black}",
    "statusBarItem.prominentForeground": "#${bright_white}",
    "statusBarItem.prominentBackground": "#${magenta}",
    "statusBarItem.prominentHoverBackground": "#${red}",
    "statusBarItem.remoteBackground": "#${green}",
    "statusBarItem.remoteForeground": "#${bright_white}",
    "statusBarItem.errorBackground": "#${red}",
    "statusBarItem.errorForeground": "#${bright_white}",
    "statusBarItem.warningBackground": "#${yellow}",
    "statusBarItem.warningForeground": "#${bright_white}",

    "titleBar.activeBackground": "#${black}",
    "titleBar.activeForeground": "#${white}",
    "titleBar.inactiveBackground": "#${black}",
    "titleBar.inactiveForeground": "#${bright_black}",

    "menubar.selectionForeground": "#${white}",
    "menubar.selectionBackground": "#${black}",
    "menu.foreground": "#${white}",
    "menu.background": "#${black}",
    "menu.selectionForeground": "#${white}",
    "menu.selectionBackground": "#${black}",
    "menu.separatorBackground": "#${bright_white}",

    "commandCenter.foreground": "#${white}",
    "commandCenter.activeForeground": "#${bright_white}",
    "commandCenter.background": "#${black}",
    "commandCenter.activeBackground": "#${black}",

    "notificationCenterHeader.foreground": "#${white}",
    "notificationCenterHeader.background": "#${black}",
    "notifications.foreground": "#${white}",
    "notifications.background": "#${black}",
    "notificationLink.foreground": "#${blue}",
    "notificationsErrorIcon.foreground": "#${red}",
    "notificationsWarningIcon.foreground": "#${yellow}",
    "notificationsInfoIcon.foreground": "#${blue}",

    "banner.background": "#${black}",
    "banner.foreground": "#${white}",
    "banner.iconForeground": "#${blue}",

    "extensionButton.prominentBackground": "#${green}",
    "extensionButton.prominentForeground": "#${bright_white}",
    "extensionButton.prominentHoverBackground": "#${black}",
    "extensionBadge.remoteBackground": "#${bright_red}",
    "extensionBadge.remoteForeground": "#${bright_white}",
    "extensionIcon.starForeground": "#${yellow}",
    "extensionIcon.verifiedForeground": "#${blue}",
    "extensionIcon.preReleaseForeground": "#${bright_red}",

    "pickerGroup.foreground": "#${bright_black}",
    "quickInput.background": "#${black}",
    "quickInput.foreground": "#${white}",
    "quickInputList.focusBackground": "#${bright_black}",
    "quickInputList.focusForeground": "#${bright_white}",
    "quickInputList.focusIconForeground": "#${bright_white}",

    "keybindingLabel.background": "#${black}",
    "keybindingLabel.foreground": "#${white}",
    "keybindingTable.headerBackground": "#${black}",
    "keybindingTable.rowsBackground": "#${black}",

    "terminal.background": "#${black}",
    "terminal.foreground": "#${white}",
    "terminal.ansiBlack": "#${black}",
    "terminal.ansiRed": "#${red}",
    "terminal.ansiGreen": "#${green}",
    "terminal.ansiYellow": "#${yellow}",
    "terminal.ansiBlue": "#${blue}",
    "terminal.ansiMagenta": "#${magenta}",
    "terminal.ansiCyan": "#${cyan}",
    "terminal.ansiWhite": "#${white}",
    "terminal.ansiBrightBlack": "#${bright_black}",
    "terminal.ansiBrightRed": "#${bright_red}",
    "terminal.ansiBrightGreen": "#${bright_green}",
    "terminal.ansiBrightYellow": "#${bright_yellow}",
    "terminal.ansiBrightBlue": "#${bright_blue}",
    "terminal.ansiBrightMagenta": "#${bright_magenta}",
    "terminal.ansiBrightCyan": "#${bright_cyan}",
    "terminal.ansiBrightWhite": "#${bright_white}",
    "terminalCursor.foreground": "#${white}",

    "debugToolBar.background": "#${black}",
    "debugView.stateLabelForeground": "#${bright_white}",
    "debugView.stateLabelBackground": "#${blue}",
    "debugView.valueChangedHighlight": "#${blue}",
    "debugTokenExpression.name": "#${magenta}",
    "debugTokenExpression.value": "#${white}",
    "debugTokenExpression.string": "#${green}",
    "debugTokenExpression.boolean": "#${bright_red}",
    "debugTokenExpression.number": "#${bright_red}",
    "debugTokenExpression.error": "#${red}",

    "testing.iconFailed": "#${red}",
    "testing.iconErrored": "#${bright_magenta}",
    "testing.iconPassed": "#${green}",
    "testing.runAction": "#${bright_black}",
    "testing.iconQueued": "#${yellow}",
    "testing.iconUnset": "#${bright_black}",
    "testing.iconSkipped": "#${magenta}",
    "testing.peekHeaderBackground": "#${black}",
    "testing.message.error.decorationForeground": "#${white}",
    "testing.message.error.lineBackground": "#${red}20",
    "testing.message.info.decorationForeground": "#${white}",
    "testing.message.info.lineBackground": "#${blue}20",

    "welcomePage.background": "#${black}",
    "welcomePage.progress.background": "#${bright_black}",
    "welcomePage.progress.foreground": "#${blue}",
    "welcomePage.tileBackground": "#${black}",
    "welcomePage.tileHoverBackground": "#${black}",
    "walkThrough.embeddedEditorBackground": "#${black}",

    "gitDecoration.addedResourceForeground": "#${green}",
    "gitDecoration.modifiedResourceForeground": "#${magenta}",
    "gitDecoration.deletedResourceForeground": "#${red}",
    "gitDecoration.renamedResourceForeground": "#${cyan}",
    "gitDecoration.stageModifiedResourceForeground": "#${magenta}",
    "gitDecoration.stageDeletedResourceForeground": "#${red}",
    "gitDecoration.untrackedResourceForeground": "#${bright_red}",
    "gitDecoration.ignoredResourceForeground": "#${bright_black}",
    "gitDecoration.conflictingResourceForeground": "#${yellow}",
    "gitDecoration.submoduleResourceForeground": "#${bright_magenta}",

    "settings.headerForeground": "#${white}",
    "settings.modifiedItemIndicator": "#${blue}",
    "settings.dropdownBackground": "#${black}",
    "settings.dropdownForeground": "#${white}",
    "settings.checkboxBackground": "#${black}",
    "settings.checkboxForeground": "#${white}",
    "settings.rowHoverBackground": "#${black}",
    "settings.textInputBackground": "#${black}",
    "settings.textInputForeground": "#${white}",
    "settings.numberInputBackground": "#${black}",
    "settings.numberInputForeground": "#${white}",
    "settings.focusedRowBackground": "#${black}",
    "settings.headerBorder": "#${white}",
    "settings.sashBorder": "#${white}",

    "breadcrumb.foreground": "#${white}",
    "breadcrumb.background": "#${black}",
    "breadcrumb.focusForeground": "#${bright_white}",
    "breadcrumb.activeSelectionForeground": "#${bright_white}",
    "breadcrumbPicker.background": "#${black}",

    "editor.snippetTabstopHighlightBackground": "#${black}",
    "editor.snippetFinalTabstopHighlightBackground": "#${bright_black}",

    "symbolIcon.arrayForeground": "#${white}",
    "symbolIcon.booleanForeground": "#${bright_red}",
    "symbolIcon.classForeground": "#${yellow}",
    "symbolIcon.constantForeground": "#${bright_red}",
    "symbolIcon.constructorForeground": "#${blue}",
    "symbolIcon.enumeratorForeground": "#${bright_red}",
    "symbolIcon.enumeratorMemberForeground": "#${blue}",
    "symbolIcon.eventForeground": "#${yellow}",
    "symbolIcon.fieldForeground": "#${red}",
    "symbolIcon.fileForeground": "#${white}",
    "symbolIcon.folderForeground": "#${white}",
    "symbolIcon.functionForeground": "#${blue}",
    "symbolIcon.interfaceForeground": "#${blue}",
    "symbolIcon.keywordForeground": "#${magenta}",
    "symbolIcon.methodForeground": "#${blue}",
    "symbolIcon.moduleForeground": "#${white}",
    "symbolIcon.namespaceForeground": "#${white}",
    "symbolIcon.nullForeground": "#${bright_magenta}",
    "symbolIcon.numberForeground": "#${bright_red}",
    "symbolIcon.propertyForeground": "#${white}",
    "symbolIcon.snippetForeground": "#${white}",
    "symbolIcon.stringForeground": "#${green}",
    "symbolIcon.structForeground": "#${yellow}",
    "symbolIcon.textForeground": "#${white}",
    "symbolIcon.variableForeground": "#${red}",

    "debugIcon.breakpointForeground": "#${red}",
    "debugIcon.breakpointDisabledForeground": "#${bright_black}",
    "debugIcon.breakpointUnverifiedForeground": "#${black}",
    "debugIcon.breakpointCurrentStackframeForeground": "#${yellow}",
    "debugIcon.breakpointStackframeForeground": "#${bright_magenta}",
    "debugIcon.startForeground": "#${green}",
    "debugIcon.pauseForeground": "#${blue}",
    "debugIcon.stopForeground": "#${red}",
    "debugIcon.disconnectForeground": "#${red}",
    "debugIcon.restartForeground": "#${green}",
    "debugIcon.stepOverForeground": "#${blue}",
    "debugIcon.stepIntoForeground": "#${cyan}",
    "debugIcon.stepOutForeground": "#${magenta}",
    "debugIcon.continueForeground": "#${green}",
    "debugIcon.stepBackForeground": "#${bright_magenta}",
    "debugConsole.infoForeground": "#${white}",
    "debugConsole.warningForeground": "#${yellow}",
    "debugConsole.errorForeground": "#${red}",
    "debugConsole.sourceForeground": "#${white}",
    "debugConsoleInputIcon.foreground": "#${white}",

    "notebook.editorBackground": "#${black}",
    "notebook.cellBorderColor": "#${bright_black}",
    "notebook.cellHoverBackground": "#${black}",
    "notebook.cellToolbarSeparator": "#${black}",
    "notebook.cellEditorBackground": "#${black}",
    "notebook.focusedCellBackground": "#${black}",
    "notebook.focusedCellBorder": "#${blue}",
    "notebook.focusedEditorBorder": "#${blue}",
    "notebook.inactiveFocusedCellBorder": "#${bright_black}",
    "notebook.selectedCellBackground": "#${black}",
    "notebookStatusErrorIcon.foreground": "#${red}",
    "notebookStatusRunningIcon.foreground": "#${cyan}",
    "notebookStatusSuccessIcon.foreground": "#${green}",

    "charts.foreground": "#${white}",
    "charts.lines": "#${white}",
    "charts.red": "#${red}",
    "charts.blue": "#${blue}",
    "charts.yellow": "#${yellow}",
    "charts.orange": "#${bright_red}",
    "charts.green": "#${green}",
    "charts.purple": "#${magenta}",
    "ports.iconRunningProcessForeground": "#${bright_red}"
}
EOF
)

# Generate syntax token customizations
syntax_colors=$(cat <<EOF
{
    "comments": "#${bright_black}",
    "strings": "#${green}",
    "keywords": "#${magenta}",
    "numbers": "#${bright_red}",
    "types": "#${yellow}",
    "functions": "#${blue}",
    "variables": "#${red}",
    "textMateRules": [
        {
            "name": "Comment",
            "scope": ["comment", "punctuation.definition.comment"],
            "settings": { "fontStyle": "italic", "foreground": "#${bright_black}" }
        },
        {
            "name": "Variables",
            "scope": ["variable", "variable.parameter", "entity.name.variable"],
            "settings": { "foreground": "#${red}" }
        },
        {
            "name": "Properties",
            "scope": ["variable.other.object.property"],
            "settings": { "foreground": "#${blue}" }
        },
        {
            "name": "Keywords",
            "scope": ["keyword", "storage.modifier", "keyword.control"],
            "settings": { "foreground": "#${magenta}" }
        },
        {
            "name": "Types",
            "scope": ["keyword.type", "storage.type.primitive", "support.type"],
            "settings": { "foreground": "#${cyan}" }
        },
        {
            "name": "Functions",
            "scope": ["entity.name.function", "support.function", "variable.function"],
            "settings": { "foreground": "#${blue}" }
        },
        {
            "name": "Strings",
            "scope": ["string", "constant.other.symbol"],
            "settings": { "foreground": "#${green}" }
        },
        {
            "name": "Numbers",
            "scope": ["constant.numeric", "constant.language", "keyword.other.unit"],
            "settings": { "foreground": "#${bright_red}" }
        },
        {
            "name": "Classes",
            "scope": ["entity.name", "support.class", "entity.name.type"],
            "settings": { "foreground": "#${yellow}" }
        },
        {
            "name": "Tags",
            "scope": ["entity.name.tag"],
            "settings": { "foreground": "#${red}" }
        },
        {
            "name": "Attributes",
            "scope": ["entity.other.attribute-name"],
            "settings": { "foreground": "#${blue}" }
        },
        {
            "name": "Invalid",
            "scope": ["invalid", "invalid.illegal"],
            "settings": { "foreground": "#${red}" }
        },
        {
            "name": "Deprecated",
            "scope": ["invalid.deprecated"],
            "settings": { "foreground": "#${bright_magenta}" }
        },
        {
            "name": "Regex",
            "scope": ["string.regexp", "constant.character.escape"],
            "settings": { "foreground": "#${cyan}" }
        },
        {
            "name": "Embedded",
            "scope": ["punctuation.section.embedded"],
            "settings": { "foreground": "#${bright_magenta}" }
        }
    ]
}
EOF
)

# Merge into settings.json (preserves all other settings)
if jq --argjson ui "$ui_colors" \
   --argjson syntax "$syntax_colors" \
   '. * {
     "workbench.colorCustomizations": $ui,
     "editor.tokenColorCustomizations": $syntax
   }' "$settings_file" > "$settings_file.tmp" 2>/dev/null; then
    mv "$settings_file.tmp" "$settings_file"
    echo "VS Code theme applied"
else
    rm -f "$settings_file.tmp"
    echo "Error: Failed to update settings.json"
    exit 1
fi
