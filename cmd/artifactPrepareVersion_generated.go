// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/piperenv"
	"github.com/SAP/jenkins-library/pkg/splunk"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/SAP/jenkins-library/pkg/validation"
	"github.com/spf13/cobra"
)

type artifactPrepareVersionOptions struct {
	AdditionalTargetTools       []string `json:"additionalTargetTools,omitempty" validate:"possible-values=custom docker dub golang gradle helm maven mta npm pip sbt yarn"`
	AdditionalTargetDescriptors []string `json:"additionalTargetDescriptors,omitempty"`
	BuildTool                   string   `json:"buildTool,omitempty" validate:"possible-values=custom docker dub golang gradle helm maven mta npm pip sbt yarn"`
	CommitUserName              string   `json:"commitUserName,omitempty"`
	CustomVersionField          string   `json:"customVersionField,omitempty"`
	CustomVersionSection        string   `json:"customVersionSection,omitempty"`
	CustomVersioningScheme      string   `json:"customVersioningScheme,omitempty" validate:"possible-values=docker maven pep440 semver2"`
	DockerVersionSource         string   `json:"dockerVersionSource,omitempty"`
	FetchCoordinates            bool     `json:"fetchCoordinates,omitempty"`
	FilePath                    string   `json:"filePath,omitempty"`
	GlobalSettingsFile          string   `json:"globalSettingsFile,omitempty"`
	IncludeCommitID             bool     `json:"includeCommitId,omitempty"`
	IsOptimizedAndScheduled     bool     `json:"isOptimizedAndScheduled,omitempty"`
	M2Path                      string   `json:"m2Path,omitempty"`
	Password                    string   `json:"password,omitempty"`
	ProjectSettingsFile         string   `json:"projectSettingsFile,omitempty"`
	ShortCommitID               bool     `json:"shortCommitId,omitempty"`
	TagPrefix                   string   `json:"tagPrefix,omitempty"`
	UnixTimestamp               bool     `json:"unixTimestamp,omitempty"`
	Username                    string   `json:"username,omitempty"`
	VersioningTemplate          string   `json:"versioningTemplate,omitempty"`
	VersioningType              string   `json:"versioningType,omitempty" validate:"possible-values=cloud cloud_noTag library"`
	CustomTLSCertificateLinks   []string `json:"customTlsCertificateLinks,omitempty"`
}

type artifactPrepareVersionCommonPipelineEnvironment struct {
	artifactVersion         string
	originalArtifactVersion string
	artifactID              string
	groupID                 string
	packaging               string
	git                     struct {
		commitID      string
		headCommitID  string
		commitMessage string
	}
}

func (p *artifactPrepareVersionCommonPipelineEnvironment) persist(path, resourceName string) {
	content := []struct {
		category string
		name     string
		value    interface{}
	}{
		{category: "", name: "artifactVersion", value: p.artifactVersion},
		{category: "", name: "originalArtifactVersion", value: p.originalArtifactVersion},
		{category: "", name: "artifactId", value: p.artifactID},
		{category: "", name: "groupId", value: p.groupID},
		{category: "", name: "packaging", value: p.packaging},
		{category: "git", name: "commitId", value: p.git.commitID},
		{category: "git", name: "headCommitId", value: p.git.headCommitID},
		{category: "git", name: "commitMessage", value: p.git.commitMessage},
	}

	errCount := 0
	for _, param := range content {
		err := piperenv.SetResourceParameter(path, resourceName, filepath.Join(param.category, param.name), param.value)
		if err != nil {
			log.Entry().WithError(err).Error("Error persisting piper environment.")
			errCount++
		}
	}
	if errCount > 0 {
		log.Entry().Error("failed to persist Piper environment")
	}
}

// ArtifactPrepareVersionCommand Prepares and potentially updates the artifact's version before building the artifact.
func ArtifactPrepareVersionCommand() *cobra.Command {
	const STEP_NAME = "artifactPrepareVersion"

	metadata := artifactPrepareVersionMetadata()
	var stepConfig artifactPrepareVersionOptions
	var startTime time.Time
	var commonPipelineEnvironment artifactPrepareVersionCommonPipelineEnvironment
	var logCollector *log.CollectorHook
	var splunkClient *splunk.Splunk
	telemetryClient := &telemetry.Telemetry{}

	var createArtifactPrepareVersionCmd = &cobra.Command{
		Use:   STEP_NAME,
		Short: "Prepares and potentially updates the artifact's version before building the artifact.",
		Long: `Prepares and potentially updates the artifact's version before building the artifact.

The continuous delivery process requires that each build is done with a unique version number.
There are two common patterns found:

### 1. Continuous Deployment pattern with automatic versioning

The team has full authority on ` + "`" + `<major>.<minor>.<patch>` + "`" + ` and can increase any part whenever required.
Nonetheless, the automatic versioning makes sure that every build will create a unique version by appending ` + "`" + `<major>.<minor>.<patch>` + "`" + ` with a buildversion (we use a timestamp) and optionally the commitId.

In order to represent this version also in the version control system the new unique version will be pushed with a dedicated tag (` + "`" + `<tagPrefix><major>.<minor>.<patch><unique extension>` + "`" + `).

Depending on the build tool used and thus the allowed versioning format the ` + "`" + `<unique extension>` + "`" + ` varies.

**Remarks:**

* There is no commit to master since this would create a perpetuum mobile and just trigger the next automatic build with automatic versioning, and so on ...
* Not creating a tag would lead to a loss of the final artifact version in scm which often is not acceptable
* You need to ensure that your CI/CD system can push back to your SCM (via providing ssh or HTTP(s) credentials)

**This pattern is the default** behavior (` + "`" + `versioningType: cloud` + "`" + `) since this is suitable for most cloud deliveries.

It is possible to use ` + "`" + `versioningType: cloud_noTag` + "`" + ` which has a slightly different behavior than described above:

* The new version will NOT be written as tag into the SCM but it is only available in the corresponding CI/CD workspace
* IMPORTANT NOTICE: Using the option ` + "`" + `cloud_noTag` + "`" + ` should not be picked in case you need to ensure a fully traceable path from SCM commit to your build artifact.

### 2. Pure version ` + "`" + `<major>.<minor>.<patch>` + "`" + `

This pattern is often used by teams that have cloud deliveries with no fully automated procedure, e.g. delivery after each takt.
Another typical use-case is development of a library with regular releases where the versioning pattern should be consumable and thus ideally complies to a ` + "`" + `<major>.<minor>.<patch>` + "`" + ` pattern.

The version is then either manually set by the team in the course of the development process or automatically pushed to master after a successful release.

Unlike for the _Continuous Deloyment_ pattern described above, in this case there is no dedicated tagging required for the build process since the version is already available in the repository.

Configuration of this pattern is done via ` + "`" + `versioningType: library` + "`" + `.

### Support of additional build tools

Besides the ` + "`" + `buildTools` + "`" + ` provided out of the box (like ` + "`" + `maven` + "`" + `, ` + "`" + `mta` + "`" + `, ` + "`" + `npm` + "`" + `, ...) it is possible to set ` + "`" + `buildTool: custom` + "`" + `.

This allows you to provide automatic versioning for tools using a:

#### file with the version as only content:

Define ` + "`" + `buildTool: custom` + "`" + ` as well as ` + "`" + `filePath: <path to your file>` + "`" + `

**Please note:** ` + "`" + `<path to your file>` + "`" + ` need to point either to a ` + "`" + `*.txt` + "`" + ` file or to a file without extension.

#### ` + "`" + `ini` + "`" + ` file containing the version:

Define ` + "`" + `buildTool: custom` + "`" + `, ` + "`" + `filePath: <path to your ini-file>` + "`" + ` as well as parameters ` + "`" + `versionSection` + "`" + ` and ` + "`" + `versionSource` + "`" + ` to point to the version location (section & parameter name) within the file.

**Please note:** ` + "`" + `<path to your file>` + "`" + ` need to point either to a ` + "`" + `*.cfg` + "`" + ` or a ` + "`" + `*.ini` + "`" + ` file.

#### ` + "`" + `json` + "`" + ` file containing the version:

Define ` + "`" + `buildTool: custom` + "`" + `, ` + "`" + `filePath: <path to your *.json file` + "`" + ` as well as parameter ` + "`" + `versionSource` + "`" + ` to point to the parameter containing the version.

#### ` + "`" + `yaml` + "`" + ` file containing the version

Define ` + "`" + `buildTool: custom` + "`" + `, ` + "`" + `filePath: <path to your *.yml/*.yaml file` + "`" + ` as well as parameter ` + "`" + `versionSource` + "`" + ` to point to the parameter containing the version.`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			startTime = time.Now()
			log.SetStepName(STEP_NAME)
			log.SetVerbose(GeneralConfig.Verbose)

			GeneralConfig.GitHubAccessTokens = ResolveAccessTokens(GeneralConfig.GitHubTokens)

			path, _ := os.Getwd()
			fatalHook := &log.FatalHook{CorrelationID: GeneralConfig.CorrelationID, Path: path}
			log.RegisterHook(fatalHook)

			err := PrepareConfig(cmd, &metadata, STEP_NAME, &stepConfig, config.OpenPiperFile)
			if err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}
			log.RegisterSecret(stepConfig.Password)
			log.RegisterSecret(stepConfig.Username)

			if len(GeneralConfig.HookConfig.SentryConfig.Dsn) > 0 {
				sentryHook := log.NewSentryHook(GeneralConfig.HookConfig.SentryConfig.Dsn, GeneralConfig.CorrelationID)
				log.RegisterHook(&sentryHook)
			}

			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 || len(GeneralConfig.HookConfig.SplunkConfig.ProdCriblEndpoint) > 0 {
				splunkClient = &splunk.Splunk{}
				logCollector = &log.CollectorHook{CorrelationID: GeneralConfig.CorrelationID}
				log.RegisterHook(logCollector)
			}

			if err = log.RegisterANSHookIfConfigured(GeneralConfig.CorrelationID); err != nil {
				log.Entry().WithError(err).Warn("failed to set up SAP Alert Notification Service log hook")
			}

			validation, err := validation.New(validation.WithJSONNamesForStructFields(), validation.WithPredefinedErrorMessages())
			if err != nil {
				return err
			}
			if err = validation.ValidateStruct(stepConfig); err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}

			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			stepTelemetryData := telemetry.CustomData{}
			stepTelemetryData.ErrorCode = "1"
			handler := func() {
				commonPipelineEnvironment.persist(GeneralConfig.EnvRootPath, "commonPipelineEnvironment")
				config.RemoveVaultSecretFiles()
				stepTelemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				stepTelemetryData.ErrorCategory = log.GetErrorCategory().String()
				stepTelemetryData.PiperCommitHash = GitCommit
				telemetryClient.SetData(&stepTelemetryData)
				telemetryClient.Send()
				if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
					splunkClient.Initialize(GeneralConfig.CorrelationID,
						GeneralConfig.HookConfig.SplunkConfig.Dsn,
						GeneralConfig.HookConfig.SplunkConfig.Token,
						GeneralConfig.HookConfig.SplunkConfig.Index,
						GeneralConfig.HookConfig.SplunkConfig.SendLogs)
					splunkClient.Send(telemetryClient.GetData(), logCollector)
				}
				if len(GeneralConfig.HookConfig.SplunkConfig.ProdCriblEndpoint) > 0 {
					splunkClient.Initialize(GeneralConfig.CorrelationID,
						GeneralConfig.HookConfig.SplunkConfig.ProdCriblEndpoint,
						GeneralConfig.HookConfig.SplunkConfig.ProdCriblToken,
						GeneralConfig.HookConfig.SplunkConfig.ProdCriblIndex,
						GeneralConfig.HookConfig.SplunkConfig.SendLogs)
					splunkClient.Send(telemetryClient.GetData(), logCollector)
				}
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetryClient.Initialize(GeneralConfig.NoTelemetry, STEP_NAME)
			artifactPrepareVersion(stepConfig, &stepTelemetryData, &commonPipelineEnvironment)
			stepTelemetryData.ErrorCode = "0"
			log.Entry().Info("SUCCESS")
		},
	}

	addArtifactPrepareVersionFlags(createArtifactPrepareVersionCmd, &stepConfig)
	return createArtifactPrepareVersionCmd
}

func addArtifactPrepareVersionFlags(cmd *cobra.Command, stepConfig *artifactPrepareVersionOptions) {
	cmd.Flags().StringSliceVar(&stepConfig.AdditionalTargetTools, "additionalTargetTools", []string{}, "Additional buildTool targets where descriptors need to be updated besides the main `buildTool`.")
	cmd.Flags().StringSliceVar(&stepConfig.AdditionalTargetDescriptors, "additionalTargetDescriptors", []string{}, "Defines patterns for build descriptors which should be used for option [`additionalTargetTools`](additionaltargettools).")
	cmd.Flags().StringVar(&stepConfig.BuildTool, "buildTool", os.Getenv("PIPER_buildTool"), "Defines the tool which is used for building the artifact.")
	cmd.Flags().StringVar(&stepConfig.CommitUserName, "commitUserName", `Project Piper`, "Defines the user name which appears in version control for the versioning update (in case `versioningType: cloud`).")
	cmd.Flags().StringVar(&stepConfig.CustomVersionField, "customVersionField", os.Getenv("PIPER_customVersionField"), "For `buildTool: custom`: Defines the field which contains the version in the descriptor file.")
	cmd.Flags().StringVar(&stepConfig.CustomVersionSection, "customVersionSection", os.Getenv("PIPER_customVersionSection"), "For `buildTool: custom`: Defines the section for version retrieval in vase a *.ini/*.cfg file is used.")
	cmd.Flags().StringVar(&stepConfig.CustomVersioningScheme, "customVersioningScheme", `maven`, "For `buildTool: custom`: Defines the versioning scheme to be used.")
	cmd.Flags().StringVar(&stepConfig.DockerVersionSource, "dockerVersionSource", os.Getenv("PIPER_dockerVersionSource"), "For `buildTool: docker`: Defines the source of the version. Can be `FROM`, any supported _buildTool_ or an environment variable name.")
	cmd.Flags().BoolVar(&stepConfig.FetchCoordinates, "fetchCoordinates", false, "If set to `true` the step will retreive artifact coordinates and store them in the common pipeline environment.")
	cmd.Flags().StringVar(&stepConfig.FilePath, "filePath", os.Getenv("PIPER_filePath"), "Defines a custom path to the descriptor file. Build tool specific defaults are used (e.g. `maven: pom.xml`, `npm: package.json`, `mta: mta.yaml`).")
	cmd.Flags().StringVar(&stepConfig.GlobalSettingsFile, "globalSettingsFile", os.Getenv("PIPER_globalSettingsFile"), "Maven only - Path to the mvn settings file that should be used as global settings file.")
	cmd.Flags().BoolVar(&stepConfig.IncludeCommitID, "includeCommitId", true, "Defines if the automatically generated version (`versioningType: cloud`) should include the commit id hash.")
	cmd.Flags().BoolVar(&stepConfig.IsOptimizedAndScheduled, "isOptimizedAndScheduled", false, "Whether the pipeline runs in optimized mode and the current execution is a scheduled one")
	cmd.Flags().StringVar(&stepConfig.M2Path, "m2Path", os.Getenv("PIPER_m2Path"), "Maven only - Path to the location of the local repository that should be used.")
	cmd.Flags().StringVar(&stepConfig.Password, "password", os.Getenv("PIPER_password"), "Password/token for git authentication.")
	cmd.Flags().StringVar(&stepConfig.ProjectSettingsFile, "projectSettingsFile", os.Getenv("PIPER_projectSettingsFile"), "Maven only - Path to the mvn settings file that should be used as project settings file.")
	cmd.Flags().BoolVar(&stepConfig.ShortCommitID, "shortCommitId", false, "Defines if a short version of the commitId should be used. GitHub format is used (first 7 characters).")
	cmd.Flags().StringVar(&stepConfig.TagPrefix, "tagPrefix", `build_`, "Defines the prefix which is used for the git tag which is written during the versioning run (only `versioningType: cloud`).")
	cmd.Flags().BoolVar(&stepConfig.UnixTimestamp, "unixTimestamp", false, "Defines if the Unix timestamp number should be used as build number instead of the standard date format.")
	cmd.Flags().StringVar(&stepConfig.Username, "username", os.Getenv("PIPER_username"), "User name for git authentication")
	cmd.Flags().StringVar(&stepConfig.VersioningTemplate, "versioningTemplate", os.Getenv("PIPER_versioningTemplate"), "DEPRECATED: Defines the template for the automatic version which will be created")
	cmd.Flags().StringVar(&stepConfig.VersioningType, "versioningType", `cloud`, "Defines the type of versioning")
	cmd.Flags().StringSliceVar(&stepConfig.CustomTLSCertificateLinks, "customTlsCertificateLinks", []string{}, "List containing download links of custom TLS certificates. This is required to ensure trusted connections to registries with custom certificates.")

	cmd.MarkFlagRequired("buildTool")
}

// retrieve step metadata
func artifactPrepareVersionMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:        "artifactPrepareVersion",
			Aliases:     []config.Alias{{Name: "artifactSetVersion", Deprecated: false}, {Name: "setVersion", Deprecated: true}},
			Description: "Prepares and potentially updates the artifact's version before building the artifact.",
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Secrets: []config.StepSecrets{
					{Name: "gitHttpsCredentialsId", Description: "Jenkins 'Username with password' credentials ID containing username/password for http access to your git repository.", Type: "jenkins"},
					{Name: "gitSshKeyCredentialsId", Description: "Jenkins 'SSH Username with private key' credentials ID ssh key for accessing your git repository. You can find details about how to generate an ssh key in the [GitHub documentation](https://docs.github.com/en/enterprise/2.15/user/articles/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent).", Type: "jenkins", Aliases: []config.Alias{{Name: "gitCredentialsId", Deprecated: true}}},
				},
				Parameters: []config.StepParameters{
					{
						Name:        "additionalTargetTools",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     []string{},
					},
					{
						Name:        "additionalTargetDescriptors",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     []string{},
					},
					{
						Name:        "buildTool",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_buildTool"),
					},
					{
						Name:        "commitUserName",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "gitUserName"}},
						Default:     `Project Piper`,
					},
					{
						Name:        "customVersionField",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_customVersionField"),
					},
					{
						Name:        "customVersionSection",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_customVersionSection"),
					},
					{
						Name:        "customVersioningScheme",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     `maven`,
					},
					{
						Name:        "dockerVersionSource",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_dockerVersionSource"),
					},
					{
						Name:        "fetchCoordinates",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     false,
					},
					{
						Name:        "filePath",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_filePath"),
					},
					{
						Name:        "globalSettingsFile",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "STEPS", "STAGES", "PARAMETERS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "maven/globalSettingsFile"}},
						Default:     os.Getenv("PIPER_globalSettingsFile"),
					},
					{
						Name:        "includeCommitId",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     true,
					},
					{
						Name: "isOptimizedAndScheduled",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "commonPipelineEnvironment",
								Param: "custom/isOptimizedAndScheduled",
							},
						},
						Scope:     []string{"PARAMETERS"},
						Type:      "bool",
						Mandatory: false,
						Aliases:   []config.Alias{},
						Default:   false,
					},
					{
						Name:        "m2Path",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "STEPS", "STAGES", "PARAMETERS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "maven/m2Path"}},
						Default:     os.Getenv("PIPER_m2Path"),
					},
					{
						Name: "password",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "gitHttpsCredentialsId",
								Param: "password",
								Type:  "secret",
							},

							{
								Name:    "gitHttpsCredentialVaultSecretName",
								Type:    "vaultSecret",
								Default: "gitHttpsCredential",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{{Name: "access_token"}},
						Default:   os.Getenv("PIPER_password"),
					},
					{
						Name:        "projectSettingsFile",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "STEPS", "STAGES", "PARAMETERS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "maven/projectSettingsFile"}},
						Default:     os.Getenv("PIPER_projectSettingsFile"),
					},
					{
						Name:        "shortCommitId",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"STEPS", "STAGES", "PARAMETERS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     false,
					},
					{
						Name:        "tagPrefix",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     `build_`,
					},
					{
						Name:        "unixTimestamp",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"STEPS", "STAGES", "PARAMETERS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     false,
					},
					{
						Name: "username",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "gitHttpsCredentialsId",
								Param: "username",
								Type:  "secret",
							},

							{
								Name:    "gitHttpsCredentialVaultSecretName",
								Type:    "vaultSecret",
								Default: "gitHttpsCredential",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_username"),
					},
					{
						Name:        "versioningTemplate",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_versioningTemplate"),
					},
					{
						Name:        "versioningType",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     `cloud`,
					},
					{
						Name:        "customTlsCertificateLinks",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     []string{},
						Conditions:  []config.Condition{{ConditionRef: "strings-equal", Params: []config.Param{{Name: "buildTool", Value: "maven"}, {Name: "buildTool", Value: "gradle"}}}},
					},
				},
			},
			Containers: []config.Container{
				{Image: "maven:3.6-jdk-8", Conditions: []config.Condition{{ConditionRef: "strings-equal", Params: []config.Param{{Name: "buildTool", Value: "maven"}}}}},
			},
			Outputs: config.StepOutputs{
				Resources: []config.StepResources{
					{
						Name: "commonPipelineEnvironment",
						Type: "piperEnvironment",
						Parameters: []map[string]interface{}{
							{"name": "artifactVersion"},
							{"name": "originalArtifactVersion"},
							{"name": "artifactId"},
							{"name": "groupId"},
							{"name": "packaging"},
							{"name": "git/commitId"},
							{"name": "git/headCommitId"},
							{"name": "git/commitMessage"},
						},
					},
				},
			},
		},
	}
	return theMetaData
}
