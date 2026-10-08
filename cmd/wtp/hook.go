package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v3"
)

// NewHookCommand creates the hook command definition
func NewHookCommand() *cli.Command {
	return &cli.Command{
		Name:  "hook",
		Usage: "Generate shell hook for cd functionality",
		Description: "Generate shell hook scripts that enable the 'wtp cd' command to change directories. " +
			"This provides a seamless navigation experience without needing subshells.\n\n" +
			"To enable the hook, add the following to your shell config:\n" +
			"  Bash (~/.bashrc):         eval \"$(wtp hook bash)\"\n" +
			"  Zsh (~/.zshrc):           eval \"$(wtp hook zsh)\"\n" +
			"  Fish (~/.config/fish/config.fish): wtp hook fish | source\n" +
			"  PowerShell ($PROFILE):    wtp hook powershell | Out-String | Invoke-Expression\n" +
			"  PowerShell 7 ($PROFILE):  wtp hook pwsh | Out-String | Invoke-Expression",
		Commands: []*cli.Command{
			{
				Name:        "bash",
				Usage:       "Generate bash hook script",
				Description: "Generate bash hook script for cd functionality",
				Action:      hookBash,
			},
			{
				Name:        "zsh",
				Usage:       "Generate zsh hook script",
				Description: "Generate zsh hook script for cd functionality",
				Action:      hookZsh,
			},
			{
				Name:        "fish",
				Usage:       "Generate fish hook script",
				Description: "Generate fish hook script for cd functionality",
				Action:      hookFish,
			},
			{
				Name:        "powershell",
				Usage:       "Generate Windows PowerShell hook script",
				Description: "Generate Windows PowerShell hook script for cd functionality",
				Action:      hookPowerShell,
			},
			{
				Name:        "pwsh",
				Usage:       "Generate PowerShell 7 hook script",
				Description: "Generate PowerShell 7 hook script for cd functionality",
				Action:      hookPowerShell,
			},
		},
	}
}

func hookBash(_ context.Context, cmd *cli.Command) error {
	w := cmd.Root().Writer
	if w == nil {
		w = os.Stdout
	}
	return printBashHook(w)
}

func hookZsh(_ context.Context, cmd *cli.Command) error {
	w := cmd.Root().Writer
	if w == nil {
		w = os.Stdout
	}
	return printZshHook(w)
}

func hookFish(_ context.Context, cmd *cli.Command) error {
	w := cmd.Root().Writer
	if w == nil {
		w = os.Stdout
	}
	return printFishHook(w)
}

func hookPowerShell(_ context.Context, cmd *cli.Command) error {
	w := cmd.Root().Writer
	if w == nil {
		w = os.Stdout
	}
	return printPowerShellHook(w)
}

func printBashHook(w io.Writer) error {
	_, err := fmt.Fprintln(w, `# wtp cd command hook for bash
wtp() {
    for arg in "$@"; do
        if [[ "$arg" == "--generate-shell-completion" ]]; then
            command wtp "$@"
            return $?
        fi
    done
    if [[ "$1" == "cd" ]]; then
        local target_dir
        if [[ -z "$2" ]]; then
            target_dir=$(command wtp cd 2>/dev/null)
        else
            target_dir=$(command wtp cd "$2" 2>/dev/null)
        fi
        if [[ $? -eq 0 && -n "$target_dir" ]]; then
            cd "$target_dir"
        else
            if [[ -z "$2" ]]; then
                command wtp cd
            else
                command wtp cd "$2"
            fi
        fi
    elif [[ "$1" == "add" ]]; then
        for arg in "$@"; do
            if [[ "$arg" == "--help" || "$arg" == "-h" ]]; then
                command wtp "$@"
                return $?
            fi
        done

        if [[ ! -t 1 ]]; then
            command wtp "$@"
            return $?
        fi

        local target_dir
        target_dir=$(command wtp "$@" --quiet)
        local wtp_status=$?
        if [[ $wtp_status -eq 0 && -n "$target_dir" ]]; then
            cd "$target_dir" || return $?
        fi
        return $wtp_status
    else
        command wtp "$@"
    fi
}`)

	return err
}

func printZshHook(w io.Writer) error {
	_, err := fmt.Fprintln(w, `# wtp cd command hook for zsh
wtp() {
    for arg in "$@"; do
        if [[ "$arg" == "--generate-shell-completion" ]]; then
            command wtp "$@"
            return $?
        fi
    done
    if [[ "$1" == "cd" ]]; then
        local target_dir
        if [[ -z "$2" ]]; then
            target_dir=$(command wtp cd 2>/dev/null)
        else
            target_dir=$(command wtp cd "$2" 2>/dev/null)
        fi
        if [[ $? -eq 0 && -n "$target_dir" ]]; then
            cd "$target_dir"
        else
            if [[ -z "$2" ]]; then
                command wtp cd
            else
                command wtp cd "$2"
            fi
        fi
    elif [[ "$1" == "add" ]]; then
        for arg in "$@"; do
            if [[ "$arg" == "--help" || "$arg" == "-h" ]]; then
                command wtp "$@"
                return $?
            fi
        done

        if [[ ! -t 1 ]]; then
            command wtp "$@"
            return $?
        fi

        local target_dir
        target_dir=$(command wtp "$@" --quiet)
        local wtp_status=$?
        if [[ $wtp_status -eq 0 && -n "$target_dir" ]]; then
            cd "$target_dir" || return $?
        fi
        return $wtp_status
    else
        command wtp "$@"
    fi
}`)

	return err
}

func printFishHook(w io.Writer) error {
	_, err := fmt.Fprintln(w, `# wtp cd command hook for fish
function wtp
    for arg in $argv
        if test "$arg" = "--generate-shell-completion"
            command wtp $argv
            return $status
        end
    end
    if test "$argv[1]" = "cd"
        set -l target_dir
        if test -z "$argv[2]"
            set target_dir (command wtp cd 2>/dev/null)
        else
            set target_dir (command wtp cd $argv[2] 2>/dev/null)
        end
        if test $status -eq 0 -a -n "$target_dir"
            cd "$target_dir"
        else
            if test -z "$argv[2]"
                command wtp cd
            else
                command wtp cd $argv[2]
            end
        end
    else if test "$argv[1]" = "add"
        for arg in $argv
            if test "$arg" = "--help"; or test "$arg" = "-h"
                command wtp $argv
                return $status
            end
        end

        if not isatty stdout
            command wtp $argv
            return $status
        end

        set -l target_dir (command wtp $argv --quiet)
        set -l wtp_status $status
        if test $wtp_status -eq 0 -a -n "$target_dir"
            cd "$target_dir"
            or return $status
        end
        return $wtp_status
    else
        command wtp $argv
    end
end`)

	return err
}

func printPowerShellHook(w io.Writer) error {
	_, err := fmt.Fprintln(w, `# wtp cd command hook for PowerShell
function wtp {
    param(
        [Parameter(ValueFromRemainingArguments = $true)]
        [string[]]$WtpArgs
    )

    $wtpCommand = Get-Command wtp.exe -CommandType Application -ErrorAction SilentlyContinue
    if ($null -eq $wtpCommand) {
        $wtpCommand = Get-Command wtp -CommandType Application -ErrorAction SilentlyContinue
    }
    if ($null -eq $wtpCommand) {
        Write-Error "wtp executable not found in PATH"
        return
    }
    $wtpExe = $wtpCommand.Source

    foreach ($arg in $WtpArgs) {
        if ($arg -eq "--generate-shell-completion") {
            & $wtpExe @WtpArgs
            return
        }
    }

    if ($WtpArgs.Count -gt 0 -and $WtpArgs[0] -eq "cd") {
        if ($WtpArgs.Count -gt 1) {
            $targetDir = & $wtpExe cd $WtpArgs[1] 2>$null
        } else {
            $targetDir = & $wtpExe cd 2>$null
        }

        if ($LASTEXITCODE -eq 0 -and -not [string]::IsNullOrWhiteSpace($targetDir)) {
            Set-Location -LiteralPath $targetDir
        } else {
            & $wtpExe @WtpArgs
        }
        return
    }

    if ($WtpArgs.Count -gt 0 -and $WtpArgs[0] -eq "add") {
        foreach ($arg in $WtpArgs) {
            if ($arg -eq "--help" -or $arg -eq "-h") {
                & $wtpExe @WtpArgs
                return
            }
        }

        $targetDir = & $wtpExe @WtpArgs --quiet
        $wtpStatus = $LASTEXITCODE
        if ($wtpStatus -eq 0 -and -not [string]::IsNullOrWhiteSpace($targetDir)) {
            Set-Location -LiteralPath $targetDir
        }
        return
    }

    & $wtpExe @WtpArgs
}`)

	return err
}
