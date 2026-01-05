#!/usr/bin/env bash
#
# Peachy Template Installer
# Fetches and installs custom templates from the Peachy repository
#
# Usage: curl -fsSL https://raw.githubusercontent.com/bjarneo/peachy/main/docs/install-templates.sh | bash
#    or: ./install-templates.sh
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m' # No Color

# Configuration
REPO_URL="https://raw.githubusercontent.com/bjarneo/peachy/main/examples/templates"
TEMPLATES_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/peachy/templates"

# Available templates with descriptions
declare -A TEMPLATES=(
    ["alacritty"]="GPU-accelerated terminal emulator"
    ["btop"]="Resource monitor with colorful interface"
    ["cava"]="Console audio visualizer"
    ["dunst"]="Lightweight notification daemon"
    ["foot"]="Fast Wayland terminal emulator"
    ["ghostty"]="Zig-based GPU-accelerated terminal"
    ["gtk"]="GTK3/GTK4 application theming"
    ["hyprland"]="Dynamic tiling Wayland compositor"
    ["hyprlock"]="Screen locker for Hyprland"
    ["iterm2"]="macOS terminal emulator (macOS only)"
    ["kitty"]="GPU-based terminal emulator"
    ["mako"]="Wayland notification daemon"
    ["neovim"]="Neovim colorscheme (aether.nvim)"
    ["rofi"]="Application launcher"
    ["swayosd"]="On-screen display for Sway/Hyprland"
    ["vencord"]="Discord client mod theme"
    ["walker"]="Wayland application launcher"
    ["warp"]="Modern terminal with AI features"
    ["waybar"]="Wayland status bar"
    ["wofi"]="Wayland application launcher"
)

# Platform-specific templates
LINUX_ONLY=("hyprland" "hyprlock" "mako" "swayosd" "waybar" "walker" "wofi" "foot")
MACOS_ONLY=("iterm2")

print_header() {
    echo -e "${BOLD}${CYAN}"
    echo "  ___  ___  ___  ___| |_  _ _ "
    echo " | . || -_||.'||  _||   || | |"
    echo " |  _||___||__,||___||_|_||_  |"
    echo " |_|                     |___|"
    echo -e "${NC}"
    echo -e "${BOLD}Template Installer${NC}"
    echo ""
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

detect_platform() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "macos" ;;
        *)       echo "unknown" ;;
    esac
}

is_installed() {
    local template="$1"
    [[ -d "$TEMPLATES_DIR/$template" ]]
}

check_command() {
    command -v "$1" &> /dev/null
}

# Check if template is relevant for current platform
is_relevant_for_platform() {
    local template="$1"
    local platform=$(detect_platform)

    # Check if it's a macOS-only template on Linux
    if [[ "$platform" == "linux" ]]; then
        for t in "${MACOS_ONLY[@]}"; do
            [[ "$template" == "$t" ]] && return 1
        done
    fi

    # Check if it's a Linux-only template on macOS
    if [[ "$platform" == "macos" ]]; then
        for t in "${LINUX_ONLY[@]}"; do
            [[ "$template" == "$t" ]] && return 1
        done
    fi

    return 0
}

download_template() {
    local template="$1"
    local dest_dir="$TEMPLATES_DIR/$template"
    local temp_dir=$(mktemp -d)

    # Files to download for each template
    local files=("template.toml")

    # Get the template file list based on template type
    case "$template" in
        alacritty) files+=("colors.toml" "post-apply") ;;
        btop) files+=("peachy.theme") ;;
        cava) files+=("theme.ini" "post-apply") ;;
        dunst) files+=("colors.conf" "post-apply") ;;
        foot) files+=("colors.ini") ;;
        ghostty) files+=("colors.conf") ;;
        gtk) files+=("colors.css") ;;
        hyprland) files+=("colors.conf" "post-apply") ;;
        hyprlock) files+=("colors.conf") ;;
        iterm2) files+=("peachy.itermcolors" "post-apply") ;;
        kitty) files+=("colors.conf" "post-apply") ;;
        mako) files+=("config" "post-apply") ;;
        neovim) files+=("aether.lua") ;;
        rofi) files+=("colors.rasi") ;;
        swayosd) files+=("style.css") ;;
        vencord) files+=("peachy.theme.css") ;;
        walker) files+=("colors.css") ;;
        warp) files+=("peachy.yaml") ;;
        waybar) files+=("colors.css" "post-apply") ;;
        wofi) files+=("colors.css") ;;
    esac

    mkdir -p "$dest_dir"

    local success=true
    for file in "${files[@]}"; do
        local url="$REPO_URL/$template/$file"
        local dest="$dest_dir/$file"

        if curl -fsSL "$url" -o "$dest" 2>/dev/null; then
            # Make post-apply scripts executable
            if [[ "$file" == "post-apply" ]]; then
                chmod +x "$dest"
            fi
        else
            # Some files are optional (like post-apply)
            if [[ "$file" != "post-apply" ]]; then
                success=false
            fi
        fi
    done

    rm -rf "$temp_dir"

    $success
}

# Interactive checkbox selector
select_templates_interactive() {
    local platform=$(detect_platform)
    local template_list=()
    local selected=()
    local cursor=0

    # Build list of relevant templates
    for template in $(echo "${!TEMPLATES[@]}" | tr ' ' '\n' | sort); do
        if is_relevant_for_platform "$template"; then
            template_list+=("$template")
            selected+=(0)  # 0 = not selected
        fi
    done

    local total=${#template_list[@]}

    # Hide cursor
    tput civis 2>/dev/null >/dev/tty || true

    # Cleanup on exit
    cleanup() {
        tput cnorm 2>/dev/null >/dev/tty || true
        stty echo 2>/dev/null || true
    }
    trap cleanup EXIT

    # Draw the menu (output to /dev/tty so it shows on screen)
    draw_menu() {
        # Move cursor up to redraw (only after first draw)
        if [[ $1 == "redraw" ]]; then
            tput cuu $((total + 5)) >/dev/tty 2>/dev/null || true
        fi

        echo -e "${BOLD}Select templates to install:${NC}" >/dev/tty
        echo -e "${DIM}Use ↑/↓ or j/k to move, Space to toggle, Enter to confirm, a to select all${NC}" >/dev/tty
        echo "" >/dev/tty

        for i in "${!template_list[@]}"; do
            local template="${template_list[$i]}"
            local desc="${TEMPLATES[$template]}"
            local checkbox
            local line_prefix="  "
            local line_suffix=""

            if [[ ${selected[$i]} -eq 1 ]]; then
                checkbox="${GREEN}[✓]${NC}"
            else
                checkbox="[ ]"
            fi

            if is_installed "$template"; then
                line_suffix=" ${DIM}(installed)${NC}"
            fi

            if [[ $i -eq $cursor ]]; then
                echo -e "${CYAN}▸${NC} $checkbox ${BOLD}$template${NC} - $desc$line_suffix" >/dev/tty
            else
                echo -e "  $checkbox $template - $desc$line_suffix" >/dev/tty
            fi
        done

        echo "" >/dev/tty
        local count=0
        for s in "${selected[@]}"; do
            ((count += s))
        done
        echo -e "${DIM}$count template(s) selected${NC}" >/dev/tty
    }

    # Initial draw
    draw_menu

    # Read keys from /dev/tty
    while true; do
        # Read a single character
        IFS= read -rsn1 key </dev/tty

        # Handle escape sequences (arrow keys)
        if [[ $key == $'\x1b' ]]; then
            read -rsn2 -t 0.1 key2 </dev/tty
            key+="$key2"
        fi

        case "$key" in
            # Up arrow or k
            $'\x1b[A'|k)
                if [[ $cursor -gt 0 ]]; then
                    ((cursor--))
                fi
                draw_menu redraw
                ;;
            # Down arrow or j
            $'\x1b[B'|j)
                if [[ $cursor -lt $((total - 1)) ]]; then
                    ((cursor++))
                fi
                draw_menu redraw
                ;;
            # Space - toggle selection
            ' ')
                if [[ ${selected[$cursor]} -eq 1 ]]; then
                    selected[$cursor]=0
                else
                    selected[$cursor]=1
                fi
                draw_menu redraw
                ;;
            # Enter - confirm
            '')
                break
                ;;
            # 'a' - select all
            a)
                local all_selected=1
                for s in "${selected[@]}"; do
                    if [[ $s -eq 0 ]]; then
                        all_selected=0
                        break
                    fi
                done
                if [[ $all_selected -eq 1 ]]; then
                    # Deselect all
                    for i in "${!selected[@]}"; do
                        selected[$i]=0
                    done
                else
                    # Select all
                    for i in "${!selected[@]}"; do
                        selected[$i]=1
                    done
                fi
                draw_menu redraw
                ;;
            # 'q' - quit
            q)
                echo "" >/dev/tty
                print_info "Installation cancelled." >/dev/tty
                exit 0
                ;;
        esac
    done

    # Show cursor again
    tput cnorm 2>/dev/null >/dev/tty || true
    echo "" >/dev/tty

    # Return selected templates to stdout (this is what gets captured)
    local result=()
    for i in "${!template_list[@]}"; do
        if [[ ${selected[$i]} -eq 1 ]]; then
            result+=("${template_list[$i]}")
        fi
    done

    echo "${result[@]}"
}

install_templates() {
    local templates=("$@")

    if [[ ${#templates[@]} -eq 0 ]]; then
        print_warning "No templates selected."
        return
    fi

    echo ""
    print_info "Installing ${#templates[@]} template(s)..."
    echo ""

    # Ensure templates directory exists
    mkdir -p "$TEMPLATES_DIR"

    local installed=0
    local failed=0

    for template in "${templates[@]}"; do
        printf "  Installing %-12s ... " "$template"

        if download_template "$template"; then
            echo -e "${GREEN}done${NC}"
            ((installed++))
        else
            echo -e "${RED}failed${NC}"
            ((failed++))
        fi
    done

    echo ""
    print_success "Installed $installed template(s)"

    if [[ $failed -gt 0 ]]; then
        print_warning "$failed template(s) failed to install"
    fi

    echo ""
    print_info "Templates installed to: $TEMPLATES_DIR"
    echo ""
    print_info "To apply templates, run:"
    echo -e "    ${CYAN}peachy apply <theme>${NC}"
    echo ""
}

list_installed() {
    echo -e "${BOLD}Installed templates:${NC}"
    echo ""

    if [[ ! -d "$TEMPLATES_DIR" ]]; then
        print_info "No templates installed yet."
        return
    fi

    local count=0
    for dir in "$TEMPLATES_DIR"/*/; do
        if [[ -d "$dir" ]]; then
            local name=$(basename "$dir")
            local desc="${TEMPLATES[$name]:-Custom template}"
            printf "  ${GREEN}✓${NC} %-12s %s\n" "$name" "$desc"
            ((count++))
        fi
    done

    if [[ $count -eq 0 ]]; then
        print_info "No templates installed yet."
    else
        echo ""
        print_info "$count template(s) installed"
    fi
}

show_help() {
    echo "Usage: $0 [OPTIONS] [TEMPLATES...]"
    echo ""
    echo "Options:"
    echo "  -h, --help      Show this help message"
    echo "  -l, --list      List installed templates"
    echo "  -a, --all       Install all templates"
    echo "  -i, --interactive  Interactive selection (default)"
    echo ""
    echo "Examples:"
    echo "  $0                    # Interactive mode"
    echo "  $0 kitty alacritty    # Install specific templates"
    echo "  $0 --all              # Install all templates"
    echo "  $0 --list             # List installed templates"
    echo ""
    echo "Available templates:"
    for template in $(echo "${!TEMPLATES[@]}" | tr ' ' '\n' | sort); do
        printf "  %-12s %s\n" "$template" "${TEMPLATES[$template]}"
    done
}

main() {
    local mode="interactive"
    local templates=()

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case "$1" in
            -h|--help)
                show_help
                exit 0
                ;;
            -l|--list)
                list_installed
                exit 0
                ;;
            -a|--all)
                mode="all"
                shift
                ;;
            -i|--interactive)
                mode="interactive"
                shift
                ;;
            -*)
                print_error "Unknown option: $1"
                echo "Use --help for usage information."
                exit 1
                ;;
            *)
                templates+=("$1")
                shift
                ;;
        esac
    done

    print_header

    local platform=$(detect_platform)
    print_info "Detected platform: ${BOLD}$platform${NC}"
    print_info "Templates directory: ${BOLD}$TEMPLATES_DIR${NC}"
    echo ""

    # Check for required tools
    if ! check_command curl; then
        print_error "curl is required but not installed."
        exit 1
    fi

    if [[ ${#templates[@]} -gt 0 ]]; then
        # Install specified templates
        install_templates "${templates[@]}"
    elif [[ "$mode" == "all" ]]; then
        # Install all templates relevant for this platform
        local all_templates=()
        for template in "${!TEMPLATES[@]}"; do
            if is_relevant_for_platform "$template"; then
                all_templates+=("$template")
            fi
        done
        install_templates "${all_templates[@]}"
    else
        # Interactive mode with checkboxes
        IFS=' ' read -r -a selected <<< "$(select_templates_interactive)"
        install_templates "${selected[@]}"
    fi
}

main "$@"
