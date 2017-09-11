package v2

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/cli/actor/v2action"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/flag"
	"code.cloudfoundry.org/cli/command/v2/shared"
)

//go:generate counterfeiter . CreateAppManifestActor

type CreateAppManifestActor interface {
	CreateApplicationManifestByNameAndSpace(appName string, spaceGUID string, filePath string) (v2action.Warnings, error)
}

type CreateAppManifestCommand struct {
	RequiredArgs    flag.AppName `positional-args:"yes"`
	FilePath        flag.Path    `short:"p" description:"Specify a path for file creation. If path not specified, manifest file is created in current working directory."`
	usage           interface{}  `usage:"CF_NAME create-app-manifest APP_NAME [-p /path/to/<app-name>_manifest.yml]"`
	relatedCommands interface{}  `related_commands:"apps, push"`

	UI          command.UI
	Config      command.Config
	SharedActor command.SharedActor
	Actor       CreateAppManifestActor
}

func (CreateAppManifestCommand) Setup(config command.Config, ui command.UI) error {
	return nil
}

func (cmd CreateAppManifestCommand) Execute(args []string) error {
	err := cmd.SharedActor.CheckTarget(cmd.Config, true, true)
	if err != nil {
		return shared.HandleError(err)
	}

	user, err := cmd.Config.CurrentUser()
	if err != nil {
		return shared.HandleError(err)
	}

	cmd.UI.DisplayText("Creating an app manifest from current settings of app {{.AppName}} in org {{.OrgName}} / space {{.SpaceName}} as {{.UserName}}...", map[string]interface{}{
		"AppName":   cmd.RequiredArgs.AppName,
		"OrgName":   cmd.Config.TargetedOrganization().Name,
		"SpaceName": cmd.Config.TargetedSpace().Name,
		"Username":  user.Name,
	})

	manifestPath := cmd.FilePath.String()
	if manifestPath == "" {
		manifestPath = fmt.Sprintf(".%s%s_manifest.yml", string(os.PathSeparator), cmd.RequiredArgs.AppName)
	}
	warnings, err := cmd.Actor.CreateApplicationManifestByNameAndSpace(cmd.RequiredArgs.AppName, cmd.Config.TargetedSpace().GUID, manifestPath)

	cmd.UI.DisplayWarnings(warnings)
	if err != nil {
		return shared.HandleError(err)
	}

	cmd.UI.DisplayOK()
	cmd.UI.DisplayText("Manifest file created successfully at {{.FilePath}}", map[string]interface{}{
		"FilePath": manifestPath,
	})

	return nil
}
