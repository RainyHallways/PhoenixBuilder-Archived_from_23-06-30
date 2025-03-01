package I18n

var I18nDict_eo map[uint16]string = map[uint16]string{
	ACME_FailedToGetCommand:             "Malsukcesis ricevi acme komandon.",
	ACME_FailedToSeek:                   "Nevalida acme-dosiero ĉar serĉado malsukcesis.",
	ACME_StructureErrorNotice:           "Nevalida strukturo",
	ACME_UnknownCommand:                 "Nekonata ACME-komando",
	Auth_BackendError:                   "Backend Eraro",
	Auth_FailedToRequestEntry:           "Malsukcesis peti eniron por via servilo, bonvolu kontroli ĉu la pasvorto estas ĝusta kaj bonvolu malŝalti la nivellimon",
	Auth_HelperNotCreated:               "Helpanto-uzanto ne estis kreita, bonvolu krei ĝin en FastBuilder Uzanto-Centro.",
	Auth_InvalidFBVersion:               "Nevalida version de FastBuilder, bonvolu ĝisdatigi.",
	Auth_InvalidHelperUsername:          "Nevalida uzantnomo por helpanto, bonvolu agordi ĝin en FastBuilder Uzantcentro.",
	Auth_InvalidToken:                   "Nevalida ensalutĵetono.",
	Auth_InvalidUser:                    "Nevalida uzanto por Uzantcentro de FastBuilde",
	Auth_ServerNotFound:                 "La servilo ne estas trovita, bonvolu kontroli la publikan staton de via servilo",
	Auth_UnauthorizedRentalServerNumber: "Neaŭtorizita luservila numero, bonvolu aldoni ĝin al via Uzantcentro de FastBuilder.",
	Auth_UserCombined:                   "Donita uzanto estis kombinita al alia konto, bonvolu ensaluti uzante la informojn de alia konto.",
	Auth_FailedToRequestEntry_TryAgain:  "Malsukcesis peti servilan eniron, bonvolu provi denove poste.",
	BDump_Author:                        "Aŭtoro",
	BDump_EarlyEOFRightWhenOpening:      "Malsukcesis legi dosieron, frua EOF? Dosiero povas esti koruptita",
	BDump_FailedToGetCmd1:               "Malsukcesis ricevi iun ajn argumenton por cmd[pos:0], dosiero povas esti koruptita.",
	BDump_FailedToGetCmd2:               "Malsukcesis ricevi iun ajn argumenton por cmd[pos1], dosiero eble difektiĝis",
	BDump_FailedToGetCmd4:               "Malsukcesis ricevi iun ajn argumenton por cmd[pos2], dosiero eble difektiĝis",
	BDump_FailedToGetCmd6:               "Malsukcesis ricevi iun ajn argumenton por cmd[pos3], dosiero eble difektiĝis",
	BDump_FailedToGetCmd7_0:             "Malsukcesis ricevi iun ajn argumenton por cmd[pos4], dosiero eble difektiĝis",
	BDump_FailedToGetCmd7_1:             "Malsukcesis ricevi iun ajn argumenton por cmd[pos5], dosiero eble difektiĝis",
	BDump_FailedToGetCmd10:              "Malsukcesis ricevi iun ajn argumenton por cmd[pos6], dosiero eble difektiĝis",
	BDump_FailedToGetCmd11:              "Malsukcesis ricevi iun ajn argumenton por cmd[pos7], dosiero eble difektiĝis",
	BDump_FailedToGetCmd12:              "Malsukcesis ricevi iun ajn argumenton por cmd[pos8], dosiero eble difektiĝis",
	BDump_FailedToGetConstructCmd:       "Malsukcesis ricevi konstrukomandojn, dosiero eble difektiĝis",
	BDump_FailedToReadAuthorInfo:        "Malsukcesis legi informojn pri aŭtoro, dosiero povas esti koruptita",
	BDump_FileNotSigned:                 "Dosiero ne estas subskribita",
	BDump_FileSigned:                    "Dosiero estas subskribita, subskribinto:%s",
	BDump_NotBDX_Invheader:              "Ne bdx-dosiero (Nevalida dosierkapo)",
	BDump_NotBDX_Invinnerheader:         "Ne bdx-dosiero (Nevalida interna dosierkapo)",
	BDump_SignedVerifying:               "Dosiero estas subskribita, kontrolante...",
	BDump_VerificationFailedFor:         "Malsukcesis kontroli la subskribon de la dosiero pro:%v",
	BDump_Warn_Reserved:                 "WARNO: BDump/Import: Uzo de rezervita komando\n",
	CommandNotFound:                     "Komando ne trovita.",
	ConnectionEstablished:               "Sukcese kreita Minecraft-telefonilo.",
	Copyright_Notice_Bouldev:            "Copyright (c) FastBuilder DevGroup, Bouldev2022",
	Copyright_Notice_Contrib:            "Kontribuantoj: Ruphane, CAIMEO, CMA2401PT",
	Crashed_No_Connection:               "konekto ne establita post tre longa tempo",
	Crashed_OS_Windows:                  "Premu ENTER por eliri.",
	Crashed_StackDump_And_Error:         "Stack-dump estis montrita supre, eraro:",
	Crashed_Tip:                         "Ho ne! FastBuilder Phoenix kraŝis!",
	CurrentDefaultDelayMode:             "Nuna defaŭlta prokrasta reĝimo",
	CurrentTasks:                        "Nunaj taskoj:",
	DelayModeSet:                        "Agordita prokrasta reĝimo",
	DelayModeSet_DelayAuto:              "Prokrasto aŭtomate agordita al: %d",
	DelayModeSet_ThresholdAuto:          "Prokrasta sojlo aŭtomate agordita al: %d",
	DelaySet:                            "Prokrasto fiksita",
	DelaySetUnavailableUnderNoneMode:    "[prokrasto fiksita] ne disponeblas kun prokrasta reĝimo: neniu",
	DelayThreshold_OnlyDiscrete:         "Prokrasta sojlo disponeblas nur kun prokrasta reĝimo: diskreta",
	DelayThreshold_Set:                  "Prokrasta sojlo agordita al: %d",
	ERRORStr:                            "ERARO",
	EnterPasswordForFBUC:                "Enigu vian pasvorton por Uzantcentro de FastBuilder: ",
	Enter_FBUC_Username:                 "Enigu vian uzantnomon de FastBuilder User Center: ",
	Enter_Rental_Server_Code:            "Bonvolu enigi vian luservilan numeron: ",
	Enter_Rental_Server_Password:        "Enigu Pasvorton (Premu [Enigu] se ne agordita, enigo ne estos ripetita): ",
	ErrorIgnored:                        "Eraro ignorita.",
	Error_MapY_Exceed:                   "En 3DMap, MapY devus esti en [20~255] (Via Enigo = %v)",
	FBUC_LoginFailed:                    "Malĝusta uzantnomo aŭ pasvorto",
	FBUC_Token_ErrOnCreate:              "Eraro dum kreado de simbola dosiero: ",
	FBUC_Token_ErrOnGen:                 "Malsukcesis generi tempan ĵetonon",
	FBUC_Token_ErrOnRemove:              "Malsukcesis forigi ĵetonan dosieron: %v",
	FBUC_Token_ErrOnSave:                "Eraro dum konservado de ĵetono: ",
	FileCorruptedError:                  "Dosiero estas koruptita",
	Get_Warning:                         "",
	IgnoredStr:                          "ignorita",
	InvalidFileError:                    "Nevalida dosiero",
	InvalidPosition:                     "Neniu pozicio akiris. (ignorebla)",
	Lang_Config_ErrOnCreate:             "Eraro dum kreado de lingvo-agorda dosiero: %v",
	Lang_Config_ErrOnSave:               "Eraro dum konservado de lingvo-agordo: %v",
	LanguageName:                        "Usona angla",
	LanguageUpdated:                     "Lingvoprefero estis ĝisdatigita",
	Logout_Done:                         "Ensalutinta de FastBuilder Uzanto-Centro.",
	Menu_BackButton:                     "< Reen",
	Menu_Cancel:                         "Nuligi",
	Menu_CurrentPath:                    "Nuna vojo",
	Menu_ExcludeCommandsOption:          "Ekskludi Komandojn",
	Menu_GetEndPos:                      "getEndPos",
	Menu_GetPos:                         "getPos",
	Menu_InvalidateCommandsOption:       "Nevalidigi Komandojn",
	Menu_Quit:                           "Forlasu Programon",
	Menu_StrictModeOption:               "Strikta Reĝimo",
	NotAnACMEFile:                       "Nevalida dosiero, ne ACME-strukturo.",
	Notice_CheckUpdate:                  "Kontrolante ĝisdatigon, bonvolu atendi...",
	Notice_iSH_Location_Service:         "Vi uzas iSH-simulilon, lokservo estas postulata por malfono, neniuj lokdatenoj estos konservitaj aŭ uzataj. Vi povas ĉesigi ĝin iam ajn.",
	Notice_OK:                           "OK\n",
	Notice_UpdateAvailable:              "Pli nova versio (%s) de PhoenixBuilder estas disponebla.\n",
	Notice_UpdateNotice:                 "Bonvolu ĝisdatigi.\n",
	Notice_ZLIB_CVE:                     "Via zlib-versio (%s) estas tro malnova ĉar ĝi enhavas konfirmitan CVE-vunereblecon, ĝisdatigo sugestita",
	Notify_NeedOp:                       "FastBuilder postulas funkciservan privilegion.",
	Notify_TurnOnCmdFeedBack:            "FastBuilder postulas gamerule sendcommandfeedback esti vera, ni jam ŝaltis ĝin, kaj memoru malŝalti ĝin",
	Omega_WaitingForOP:                  "Omega Sistemo atendas OP-Privilegion",
	Omega_Enabled:                       "Omega Sistemo Ebligita!",
	OpPrivilegeNotGrantedForOperation:   "Op-privilegio ne koncedita por ĉi tiu operacio, bonvolu doni bot-op-privilegion",
	Parsing_UnterminatedEscape:          "Nefinita fuĝo",
	Parsing_UnterminatedQuotedString:    "Nefinita citita ĉeno",
	PositionGot:                         "Pozicio akiris",
	PositionGot_End:                     "fin Pozicio akiris",
	PositionSet:                         "Pozicio aro",
	PositionSet_End:                     "fin Pozicio aro",
	QuitCorrectly:                       "Forlasu ĝuste",
	Sch_FailedToResolve:                 "Malsukcesis solvi dosieron",
	SelectLanguageOnConsole:             "Bonvolu elekti vian novan preferatan lingvon en la konzolo.",
	ServerCodeTrans:                     "Servilo",
	SimpleParser_Int_ParsingFailed:      "Analizilo: malsukcesis analizi int-argumenton",
	SimpleParser_InvEnum:                "Analizilo: nevalida enumvaloro, permesitaj valoroj estas: %s.", ///
	SimpleParser_Invalid_decider:        "Analizilo: Nevalida decidilo",
	SimpleParser_Too_few_args:           "Analizisto: Tro malmultaj argumentoj",
	Special_Startup:                     "Enŝaltita lingvo: angla\n",
	TaskCreated:                         "Tasko Kreita",
	TaskDisplayModeSet:                  "Reĝimo de montra stato de tasko agordita al: %s.",
	TaskFailedToParseCommand:            "Malsukcesis analizi komandon: %v",
	TaskNotFoundMessage:                 "Ne eblis trovi validan taskon per provizita taskoid.",
	TaskPausedNotice:                    "[Tasko %d] - Paŭzita",
	TaskResumedNotice:                   "[Tasko %d] - Rekomencita",
	TaskStateLine:                       "ID %d - CommandLine:\"%s\", Ŝtato: %s, Prokrasto: %d, DelayMode: %s, DelayThreshold: %d",
	TaskStoppedNotice:                   "[Tasko %d] - Haltis",
	TaskTTeIuKoto:                       "Tasko",
	TaskTotalCount:                      "Entute: %d",
	TaskTypeCalculating:                 "Kalkulado",
	TaskTypeDied:                        "Mortis",
	TaskTypePaused:                      "Paŭzita",
	TaskTypeRunning:                     "Kurante",
	TaskTypeSpecialTaskBreaking:         "Speciala Tasko: Rompi",
	TaskTypeSwitchedTo:                  "Taska krea tipo agordita al: %s.",
	TaskTypeUnknown:                     "Nekonata",
	Task_D_NothingGenerated:             "[Tasko %d] Nenio generita.",
	Task_DelaySet:                       "[Tasko %d] - Prokrasto fiksita: %d",
	Task_ResumeBuildFrom:                "Rekomencu Konstruadon El Bloko-Numero %v ",
	Task_SetDelay_Unavailable:           "[setdelay] estas nedisponebla kun prokrasta reĝimo: neniu",
	Task_Summary_1:                      "[Tasko %d] %v loko(j) estas ŝanĝita(j).",
	Task_Summary_2:                      "[Tasko %d] Uzata tempo: %v sekundo(j)",
	Task_Summary_3:                      "[Tasko %d] Averaĝa rapideco: %v blokoj/sekundo",
	UnsupportedACMEVersion:              "Nesubtenata ACME-strukturversio. Nur acme dosiero versio 1.2 estas subtenata.",
	Warning_ACME_Deprecated:             "AVERTO - `acme' estas malnoveca kaj estos forigita en la estonteco, bonvolu migri al alia formato anstataŭe.\n",
	Warning_UserHomeDir:                 "AVERTO - Malsukcesis akiri la hejman dosierujon de la uzanto. farita homedir=\".\";\n",
}
