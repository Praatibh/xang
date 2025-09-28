# 🚀 Xang 💬 - AI Powered Terminal Assistant

[![build](https://github.com/Praatibh/xang/actions/workflows/build.yml/badge.svg)](https://github.com/Praatibh/xang/actions/workflows/build.yml)
[![release](https://github.com/Praatibh/xang/actions/workflows/release.yml/badge.svg?branch=main)](https://github.com/Praatibh/xang/releases)
[![License](https://img.shields.io/github/license/Praatibh/xang)](https://github.com/Praatibh/xang/blob/master/LICENSE)

> Unleash the power of Google Gemini AI to streamline your command line experience with an anime-style interactive terminal assistant.

## What is Xang?

`Xang` is an intelligent terminal assistant powered by [Google Gemini](https://gemini.google.com/) that helps you build and execute commands using natural language. Simply describe what you want to do in everyday language, and Xang will generate the appropriate terminal commands for you.

### Key Features

- 🤖 **AI-Powered Command Generation**: Convert natural language to terminal commands
- 💬 **Interactive Chat Mode**: Ask questions and get AI-powered responses without leaving your terminal  
- 🎨 **Anime Character Interface**: Features a reactive ASCII art character that responds to different states
- 🔧 **System Awareness**: Automatically detects your OS, shell, editor, and other system preferences
- 📝 **Command History**: Navigate through your previous commands with arrow keys
- ⚡ **Multiple Modes**: Switch between exec mode (🚀) and chat mode (💬) with Tab
- 🔐 **Secure Configuration**: Uses Gemini API key stored locally in your config

Xang is already aware of your:
- Operating system & distribution
- Username, shell & home directory  
- Preferred editor
- Custom user preferences

## Quick Start

### Installation

Install Xang with a single command:

```shell
curl -sS https://raw.githubusercontent.com/Praatibh/xang/main/install.sh | bash
```

### First Run Setup

1. Run `xang` to start the interactive REPL mode
2. On first run, you'll be prompted to enter your [Gemini API key](https://aistudio.google.com/app/apikey)
3. Get your free API key from [Google AI Studio](https://aistudio.google.com/app/apikey)
4. The configuration will be saved to `~/.config/xang.json`

### Usage Examples

```shell
# Interactive REPL mode with anime character
xang

# Execute single command
xang "list all files in current directory"
xang "find large files over 100MB"
xang "create a backup of my documents folder"

# Process piped input
echo "analyze this data" | xang
ls -la | xang "explain what these files are"
```

## Interface Modes

### 🚀 Exec Mode (Default)
- Generates executable terminal commands from natural language
- Shows command preview with explanation before execution
- Confirms before running potentially destructive commands

### 💬 Chat Mode  
- General AI conversation and assistance
- Ask programming questions, get explanations
- No command execution, just helpful responses

Switch between modes by pressing `Tab`.

## Anime Character Reactions

Xang features a reactive ASCII art character that changes expressions based on the current state:

- **Idle**: 😊 Gentle breathing and occasional blinking
- **Thinking**: 🤔 Pondering your request  
- **Processing**: 🔄 Working on command generation
- **Success**: ⭐ Happy when commands succeed
- **Error**: ❌ Shows concern when errors occur

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Tab` | Switch between exec (🚀) and chat (💬) modes |
| `↑/↓` | Navigate command history |
| `Ctrl+H` | Show help |
| `Ctrl+L` | Clear terminal (keep history) |
| `Ctrl+R` | Reset terminal and clear history |
| `Ctrl+S` | Edit settings |
| `Ctrl+C` | Exit or interrupt |

## Configuration

Configuration is stored in `~/.config/xang.json`:

```json
{
  "gemini_key": "your-api-key-here",
  "gemini_model": "gemini-2.5-flash",
  "user_default_prompt_mode": "exec",
  "user_preferences": "I prefer verbose output and detailed explanations"
}
```

Edit your preferences to customize Xang's behavior for your specific needs.

## Building from Source

```shell
git clone https://github.com/Praatibh/xang.git
cd xang
go build -o xang .
./xang
```

### Requirements
- Go 1.19 or later
- Valid Gemini API key

## Uninstalling

To remove Xang from your system:

```shell
curl -sS https://raw.githubusercontent.com/Praatibh/xang/main/uninstall.sh | bash
```

This will remove:
- The binary from `/usr/local/bin/xang`
- Configuration file from `~/.config/xang.json`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Google Gemini API](https://ai.google.dev/)
- Inspired by the original [Yai project](https://github.com/ekkinox/yai)
- Terminal UI powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Styling with [Lip Gloss](https://github.com/charmbracelet/lipgloss)

## Support

- 🐛 [Report bugs](https://github.com/Praatibh/xang/issues)
- 💡 [Request features](https://github.com/Praatibh/xang/issues)
- 📖 [Documentation](https://github.com/Praatibh/xang/wiki)

---

**Happy coding with Xang!** 🚀
