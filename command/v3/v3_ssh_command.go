package v3

import (
	"errors"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"code.cloudfoundry.org/cli/actor/sharedaction"
	"code.cloudfoundry.org/cli/actor/v2action"
	"code.cloudfoundry.org/cli/actor/v3action"
	sshcmd "code.cloudfoundry.org/cli/cf/ssh"
	sshoptions "code.cloudfoundry.org/cli/cf/ssh/options"
	sshterminal "code.cloudfoundry.org/cli/cf/ssh/terminal"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/flag"
	sharedV2 "code.cloudfoundry.org/cli/command/v2/shared"
	"code.cloudfoundry.org/cli/command/v3/shared"
)

type V3SSHActor interface {
	GetApplicationSummaryByNameAndSpace(appName string, spaceGUID string) (v3action.ApplicationSummary, v3action.Warnings, error)
}

type V2SSHActor interface {
	GetSSHInfo() (v2action.SSHInfo, error)
}

type V3SSHCommand struct {
	RequiredArgs        flag.AppName `positional-args:"yes"`
	ProcessType         string       `short:"p" description:"The process name" required:"true"`
	ProcessIndex        int          `short:"i" description:"The process index" required:"true"`
	Command             string       `short:"c" description:"command" required:"false"`
	DisablePseudoTTY    bool         `short:"T" description:"disable pseudo-tty" required:"false"`
	ForcePseudoTTY      bool         `short:"F" description:"force pseudo-tty" required:"false"`
	RequestPseudoTTY    bool         `short:"t" description:"request pseudo-tty" required:"false"`
	Forward             []string     `short:"L" description:"forward" required:"false"`
	SkipHostValidation  bool         `short:"k" description:"skip host validation" required:"false"`
	SkipRemoteExecution bool         `short:"N" description:"skip remote execution" required:"false"`

	usage interface{} `usage:"CF_NAME v3-ssh APP_NAME"`

	UI          command.UI
	Config      command.Config
	SharedActor command.SharedActor
	Actor       V3SSHActor
	V2SSHActor  V2SSHActor
}

func (cmd *V3SSHCommand) Setup(config command.Config, ui command.UI) error {
	cmd.UI = ui
	cmd.Config = config
	cmd.SharedActor = sharedaction.NewActor()

	ccClient, _, err := shared.NewClients(config, ui, true)
	if err != nil {
		return err
	}
	cmd.Actor = v3action.NewActor(ccClient, config)

	ccClientV2, uaaClientV2, err := sharedV2.NewClients(config, ui, true)
	if err != nil {
		return err
	}
	cmd.V2SSHActor = v2action.NewActor(ccClientV2, uaaClientV2, config)

	return nil
}

func (cmd V3SSHCommand) Execute(args []string) error {
	err := cmd.SharedActor.CheckTarget(cmd.Config, true, true)
	if err != nil {
		return shared.HandleError(err)
	}

	user, err := cmd.Config.CurrentUser()
	if err != nil {
		return shared.HandleError(err)
	}

	summary, warnings, err := cmd.Actor.GetApplicationSummaryByNameAndSpace(cmd.RequiredArgs.AppName, cmd.Config.TargetedSpace().GUID)
	cmd.UI.DisplayWarnings(warnings)
	if err != nil {
		return shared.HandleError(err)
	}

	cmd.UI.DisplayTextWithFlavor("Sshing into app {{.AppName}} in org {{.OrgName}} / space {{.SpaceName}} as {{.Username}}...", map[string]interface{}{
		"AppName":   cmd.RequiredArgs.AppName,
		"OrgName":   cmd.Config.TargetedOrganization().Name,
		"SpaceName": cmd.Config.TargetedSpace().Name,
		"Username":  user.Name,
	})

	var processGUID string
	for _, process := range summary.Processes {
		if process.Type == cmd.ProcessType {
			processGUID = process.GUID
		}
	}

	if processGUID == "" {
		return errors.New("process does not exist")
	}

	sshInfo, err := cmd.V2SSHActor.GetSSHInfo()
	if err != nil {
		return shared.HandleError(err)
	}

	secureShell := sshcmd.NewSecureShell(
		sshcmd.DefaultSecureDialer(),
		sshterminal.DefaultHelper(),
		sshcmd.DefaultListenerFactory(),
		30*time.Second,
		summary.State,
		processGUID,
		true,
		sshInfo.SSHHostKeyFingerprint,
		sshInfo.SSHEndpoint,
		sshInfo.SSHPasscode,
	)

	var command []string
	if cmd.Command != "" {
		command = strings.Split(cmd.Command, " ")
	}

	opts := sshoptions.SSHOptions{
		AppName:             cmd.RequiredArgs.AppName,
		Command:             command,
		Index:               uint(cmd.ProcessIndex),
		SkipHostValidation:  cmd.SkipHostValidation,
		SkipRemoteExecution: cmd.SkipRemoteExecution,
	}

	switch {
	case cmd.DisablePseudoTTY:
		opts.TerminalRequest = sshoptions.RequestTTYNo
	case cmd.ForcePseudoTTY:
		opts.TerminalRequest = sshoptions.RequestTTYForce
	case cmd.RequestPseudoTTY:
		opts.TerminalRequest = sshoptions.RequestTTYYes
	default:
		opts.TerminalRequest = sshoptions.RequestTTYAuto
	}

	if len(cmd.Forward) > 0 {
		for _, arg := range cmd.Forward {
			forwardSpec, err := opts.ParseLocalForwardingSpec(arg)
			if err != nil {
				return err
			}
			opts.ForwardSpecs = append(opts.ForwardSpecs, *forwardSpec)
		}
	}

	err = secureShell.Connect(&opts)
	if err != nil {
		return errors.New("Error opening SSH connection: " + err.Error())
	}
	defer secureShell.Close()

	err = secureShell.LocalPortForward()
	if err != nil {
		return errors.New("Error forwarding port: " + err.Error())
	}

	if opts.SkipRemoteExecution {
		err = secureShell.Wait()
	} else {
		err = secureShell.InteractiveSession()
	}

	if err != nil {
		if exitError, ok := err.(*ssh.ExitError); ok {
			exitStatus := exitError.ExitStatus()
			if sig := exitError.Signal(); sig != "" {
				cmd.UI.DisplayText("Process terminated by signal: {{.Signal}}. Exited with {{.ExitCode}}", map[string]interface{}{
					"Signal":   sig,
					"ExitCode": exitStatus,
				})
			}
			os.Exit(exitStatus)
		} else {
			return errors.New("Error: " + err.Error())
		}
	}

	return nil
}
