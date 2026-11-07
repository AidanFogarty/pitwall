package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/AidanFogarty/pitwall/internal/auth"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "log in to your f1 account",
	Long:  "log in to your f1 account to use features that require one, such as importing sessions.",
	RunE:  runLogin,
}

func runLogin(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	out := cmd.OutOrStdout()

	browserPath, err := auth.FindBrowser()
	if err != nil {
		return err
	}

	fmt.Fprintln(out, "Pitwall will open your browser so you can sign in on the official Formula 1 website.")
	fmt.Fprintln(out)

	proceed, err := confirm(cmd.InOrStdin(), out, "Open your browser to sign in?", true)
	if err != nil {
		return err
	}
	if !proceed {
		fmt.Fprintln(out, "Login cancelled. Run `pitwall login` when you are ready.")
		return nil
	}

	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
	defer stop()

	fmt.Fprintln(out, "\nOpening your browser...")
	fmt.Fprintln(out, "Waiting for you to sign in. Press Ctrl+C to cancel.")

	session, err := auth.Login(ctx, browserPath)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			fmt.Fprintln(out, "\nLogin cancelled.")
			return nil
		}
		return err
	}

	store, err := auth.DefaultStore()
	if err != nil {
		return err
	}
	if err := store.Save(session); err != nil {
		return err
	}

	printLoggedIn(out, session)
	return nil
}

func printLoggedIn(out io.Writer, session *auth.Session) {
	fmt.Fprintln(out, "\n✓ Logged in to Formula 1")

	if !session.HasActiveSubscription() {
		fmt.Fprintln(out, "\nNote: this account has no active F1 TV subscription, so some features may not work.")
	}
}

func confirm(in io.Reader, out io.Writer, question string, defaultYes bool) (bool, error) {
	hint := "[y/N]"
	if defaultYes {
		hint = "[Y/n]"
	}

	reader := bufio.NewReader(in)
	for {
		fmt.Fprintf(out, "%s %s ", question, hint)

		line, err := reader.ReadString('\n')
		if err != nil && line == "" {
			return false, fmt.Errorf("reading answer: %w", err)
		}

		switch strings.ToLower(strings.TrimSpace(line)) {
		case "":
			return defaultYes, nil
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		}

		fmt.Fprintln(out, "Please answer y or n.")
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
