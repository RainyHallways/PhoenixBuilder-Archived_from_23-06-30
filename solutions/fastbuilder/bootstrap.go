package fastbuilder

import (
	"fmt"
	"os"
	"os/exec"
	"phoenixbuilder/fastbuilder/args"
	"phoenixbuilder/fastbuilder/credentials"
	I18n "phoenixbuilder/fastbuilder/i18n"
	script_bridge "phoenixbuilder/fastbuilder/script_engine/bridge"
	"phoenixbuilder/fastbuilder/utils"

	"github.com/pterm/pterm"

	"phoenixbuilder/fastbuilder/readline"
	_ "phoenixbuilder/plantform_specific/fix_timer"
)

func setup() {
	pterm.Error.Prefix = pterm.Prefix{
		Text:  "ERROR",
		Style: pterm.NewStyle(pterm.BgBlack, pterm.FgRed),
	}

	I18n.Init()
	// iSH.app specific, for foreground ability
	if _, err := os.Stat("/dev/location"); err == nil {
		// Call location service
		pterm.Println(pterm.Yellow(I18n.T(I18n.Notice_iSH_Location_Service)))
		cmd := exec.Command("ash", "-c", "cat /dev/location > /dev/null &")
		err := cmd.Start()
		if err != nil {
			fmt.Println(err)
		}
	}

	if !args.NoReadline() {
		readline.InitReadline()
	}
}

func display_info() {
	pterm.DefaultBox.Println(pterm.LightCyan("https://github.com/rnhws-Team/PhoenixBuilder"))
	pterm.Println(pterm.Yellow("PhoenixBuilder Dev " + args.GetFBVersion()))
        pterm.Println(pterm.Red("注意: 此版本为第三方预览版本，请勿用于实际生产用途！！！"))
	if I18n.ShouldDisplaySpecial() {
		fmt.Printf("%s", I18n.T(I18n.Special_Startup))
	}
}

func hash_check() {
	fmt.Printf(I18n.T(I18n.Notice_CheckUpdate))
	hasUpdate, latestVersion := utils.CheckUpdate(args.GetFBVersion())
	fmt.Printf(I18n.T(I18n.Notice_OK))
	if hasUpdate {
		fmt.Printf(I18n.T(I18n.Notice_UpdateAvailable), latestVersion)
		fmt.Printf(I18n.T(I18n.Notice_UpdateNotice))
		// To ensure user won't ignore it directly, can be suppressed by command line argument.
		os.Exit(0)
	}
}

func Bootstrap() {
	args.ParseArgs()
	if len(args.PackScripts()) != 0 {
		os.Exit(script_bridge.MakePackage(args.PackScripts(), args.PackScriptsOut()))
	}
	setup()
	display_info()
	defer Fatal()
	if args.DebugMode() {
		init_and_run_debug_client()
		return
	}
	if !args.ShouldDisableHashCheck() {
		hash_check()
	}

	token := loadFBTokenOrAskFBCredential()
	runInteractiveClient(token)
}

func runInteractiveClient(token string) {
	var code, serverPasswd string
	var err error
	if !args.SpecifiedServer() {
		code, serverPasswd, err = credentials.GetRentalServerCode()
	} else {
		code = args.ServerCode()
		serverPasswd = args.ServerPassword()
	}

	if err != nil {
		fmt.Println(err)
		return
	}
	env := ConfigRealEnvironment(token, code, serverPasswd)
	ptoken_succ := credentials.ProcessTokenDefault(env)
	//init_and_run_client(token, code, serverPasswd)
	if !ptoken_succ {
		panic("Failed to load token")
	}
	EstablishConnectionAndInitEnv(env)
	go EnterReadlineThread(env, nil)
	defer DestroyEnv(env)
	EnterWorkerThread(env, nil)
}

func init_and_run_debug_client() {
	env := ConfigDebugEnvironment()
	EstablishConnectionAndInitEnv(env)
	go EnterReadlineThread(env, nil)
	defer DestroyEnv(env)
	EnterWorkerThread(env, nil)
}
