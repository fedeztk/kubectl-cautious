package plugin

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/fedeztk/kubectl-cautious/pkg/config"
)

var ErrKubectl = errors.New("kubectl error")

func RunPlugin(conf *config.Config, args []string) error {
	kubeconfig, err := getKubeconfigPath()
	if err != nil {
		return err
	}

	currentCtx, err := getContext(kubeconfig)
	if err != nil {
		return err
	}

	confirm := true
	if actions := getActionsForContextInArgs(currentCtx, conf, args); actions != nil {
		log.Warn("Operating in", "context", currentCtx)
		for _, action := range actions {
			if action.DryRun {
				log.Info("Performing dry run")
				err := execKubectl(append(args, "--dry-run=client"))
				if err != nil {
					return err
				}
			}
		}
		huh.NewConfirm().
			Title("Would you like to proceed?").
			Value(&confirm).WithTheme(huh.ThemeBase16()).Run()
	} // else no actions for this context

	if !confirm { // do not run if user does not confirm
		return nil
	}

	return execKubectl(args)
}

// returns actions associated with the current context
func getActionsForContext(currentCtx string, conf *config.Config) []config.Action {
	for _, ctx := range conf.Contexts {
		// chekc if regex ctx.Name matches currentCtx
		if regexp.MustCompile(ctx.Name).MatchString(currentCtx) {
			return ctx.Actions
		}
	}
	return nil
}

// returns actions present in args
func checkActionInArgs(action string, args []string) bool {
	for _, arg := range args {
		if arg == action {
			return true
		}
	}
	return false
}

// getActionsForContextInArgs returns actions for the current context that are in args
func getActionsForContextInArgs(currentCtx string, conf *config.Config, args []string) []config.Action {
	actionsForCtx := getActionsForContext(currentCtx, conf)
	var actionsInArgs []config.Action
	for _, action := range actionsForCtx {
		if checkActionInArgs(action.Name, args) {
			actionsInArgs = append(actionsInArgs, action)
		}
	}
	return actionsInArgs
}

func execKubectl(args []string) error {
	cmd := exec.Command("kubectl", args...)
	// preserve stdin so that kubectl can apply -f -
	cmd.Stdin = os.Stdin
	output, err := cmd.CombinedOutput()
	fmt.Print(string(output))

	if err != nil {
		return ErrKubectl
	}
	return nil
}
