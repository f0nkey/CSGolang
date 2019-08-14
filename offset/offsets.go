package offset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var Signatures *signatures
var Netvars *netvars

func MarshalOffsets() {
	resp, err := http.Get("https://raw.githubusercontent.com/frk1/hazedumper/master/csgo.json")
	if err != nil {
		log.Println("MarshalOffsets:", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("MarshallOffsets:", err)

	}

	fmt.Println("bod", body)

	var off offsets
	err = json.Unmarshal(body, &off)
	if err != nil {
		log.Println("MarshalOffsets:", err)
	}
	fmt.Println("MarshalOffsets:", err)

	if offsetsOutdated(off) {
		fmt.Println("OFFSETS POSSIBLY OUTDATED")
	}

	fmt.Println("of", off)

	Signatures = &off.Signatures
	Netvars = &off.Netvars
}

func offsetsOutdated(Offsets offsets) bool {
	offsetsPostDate := time.Unix(int64(Offsets.Timestamp), 0)
	csgoLatestDate, err := time.Parse("2006-01-02", strings.Replace(getCurrentCSGOUpdate(), ".", "-", 4))
	if err != nil {
		log.Println("MarshalOffsets: bad time parsing", err)
	}

	return csgoLatestDate.After(offsetsPostDate)
}

func getCurrentCSGOUpdate() string {
	resp, err := http.Get("https://blog.counter-strike.net/index.php/category/updates/")
	if err != nil {
		log.Println("MarshalOffsets:", err)
	}
	preBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("MarshallOffsets:", err)

	}
	body := string(preBody)
	postDateIndex := strings.Index(body, `class="post_date"`)
	date := body[postDateIndex+18 : postDateIndex+28]
	return date
}

type offsets struct {
	Timestamp  int        `json:"timestamp"`
	Signatures signatures `json:"signatures"`
	Netvars    netvars    `json:"netvars"`
}

type netvars struct {
	CsGamerulesData                int32 `json:"cs_gamerules_data"`
	MArmorValue                    int32 `json:"m_ArmorValue"`
	MCollision                     int32 `json:"m_Collision"`
	MCollisionGroup                int32 `json:"m_CollisionGroup"`
	MLocal                         int32 `json:"m_Local"`
	MMoveType                      int32 `json:"m_MoveType"`
	MOriginalOwnerXuidHigh         int32 `json:"m_OriginalOwnerXuidHigh"`
	MOriginalOwnerXuidLow          int32 `json:"m_OriginalOwnerXuidLow"`
	MSurvivalGameRuleDecisionTypes int32 `json:"m_SurvivalGameRuleDecisionTypes"`
	MSurvivalRules                 int32 `json:"m_SurvivalRules"`
	MAimPunchAngle                 int32 `json:"m_aimPunchAngle"`
	MAimPunchAngleVel              int32 `json:"m_aimPunchAngleVel"`
	MAngEyeAnglesX                 int32 `json:"m_angEyeAnglesX"`
	MAngEyeAnglesY                 int32 `json:"m_angEyeAnglesY"`
	MBBombPlanted                  int32 `json:"m_bBombPlanted"`
	MBFreezePeriod                 int32 `json:"m_bFreezePeriod"`
	MBGunGameImmunity              int32 `json:"m_bGunGameImmunity"`
	MBHasDefuser                   int32 `json:"m_bHasDefuser"`
	MBHasHelmet                    int32 `json:"m_bHasHelmet"`
	MBInReload                     int32 `json:"m_bInReload"`
	MBIsDefusing                   int32 `json:"m_bIsDefusing"`
	MBIsQueuedMatchmaking          int32 `json:"m_bIsQueuedMatchmaking"`
	MBIsScoped                     int32 `json:"m_bIsScoped"`
	MBIsValveDS                    int32 `json:"m_bIsValveDS"`
	MBSpotted                      int32 `json:"m_bSpotted"`
	MBSpottedByMask                int32 `json:"m_bSpottedByMask"`
	MBStartedArming                int32 `json:"m_bStartedArming"`
	MClrRender                     int32 `json:"m_clrRender"`
	MDwBoneMatrix                  int32 `json:"m_dwBoneMatrix"`
	MFAccuracyPenalty              int32 `json:"m_fAccuracyPenalty"`
	MFFlags                        int32 `json:"m_fFlags"`
	MFlC4Blow                      int32 `json:"m_flC4Blow"`
	MFlDefuseCountDown             int32 `json:"m_flDefuseCountDown"`
	MFlDefuseLength                int32 `json:"m_flDefuseLength"`
	MFlFallbackWear                int32 `json:"m_flFallbackWear"`
	MFlFlashDuration               int32 `json:"m_flFlashDuration"`
	MFlFlashMaxAlpha               int32 `json:"m_flFlashMaxAlpha"`
	MFlLastBoneSetupTime           int32 `json:"m_flLastBoneSetupTime"`
	MFlLowerBodyYawTarget          int32 `json:"m_flLowerBodyYawTarget"`
	MFlNextAttack                  int32 `json:"m_flNextAttack"`
	MFlNextPrimaryAttack           int32 `json:"m_flNextPrimaryAttack"`
	MFlSimulationTime              int32 `json:"m_flSimulationTime"`
	MFlTimerLength                 int32 `json:"m_flTimerLength"`
	MHActiveWeapon                 int32 `json:"m_hActiveWeapon"`
	MHMyWeapons                    int32 `json:"m_hMyWeapons"`
	MHObserverTarget               int32 `json:"m_hObserverTarget"`
	MHOwner                        int32 `json:"m_hOwner"`
	MHOwnerEntity                  int32 `json:"m_hOwnerEntity"`
	MIAccountID                    int32 `json:"m_iAccountID"`
	MIClip1                        int32 `json:"m_iClip1"`
	MICompetitiveRanking           int32 `json:"m_iCompetitiveRanking"`
	MICompetitiveWins              int32 `json:"m_iCompetitiveWins"`
	MICrosshairID                  int32 `json:"m_iCrosshairId"`
	MIEntityQuality                int32 `json:"m_iEntityQuality"`
	MIFOV                          int32 `json:"m_iFOV"`
	MIFOVStart                     int32 `json:"m_iFOVStart"`
	MIGlowIndex                    int32 `json:"m_iGlowIndex"`
	MIHealth                       int32 `json:"m_iHealth"`
	MIItemDefinitionIndex          int32 `json:"m_iItemDefinitionIndex"`
	MIItemIDHigh                   int32 `json:"m_iItemIDHigh"`
	MIMostRecentModelBoneCounter   int32 `json:"m_iMostRecentModelBoneCounter"`
	MIObserverMode                 int32 `json:"m_iObserverMode"`
	MIShotsFired                   int32 `json:"m_iShotsFired"`
	MIState                        int32 `json:"m_iState"`
	MITeamNum                      int32 `json:"m_iTeamNum"`
	MLifeState                     int32 `json:"m_lifeState"`
	MNFallbackPauintptrKit         int32 `json:"m_nFallbackPauintptrKit"`
	MNFallbackSeed                 int32 `json:"m_nFallbackSeed"`
	MNFallbackStatTrak             int32 `json:"m_nFallbackStatTrak"`
	MNForceBone                    int32 `json:"m_nForceBone"`
	MNTickBase                     int32 `json:"m_nTickBase"`
	MRgflCoordinateFrame           int32 `json:"m_rgflCoordinateFrame"`
	MSzCustomName                  int32 `json:"m_szCustomName"`
	MSzLastPlaceName               int32 `json:"m_szLastPlaceName"`
	MThirdPersonViewAngles         int32 `json:"m_thirdPersonViewAngles"`
	MVecOrigin                     int32 `json:"m_vecOrigin"`
	MVecVelocity                   int32 `json:"m_vecVelocity"`
	MVecViewOffset                 int32 `json:"m_vecViewOffset"`
	MViewPunchAngle                int32 `json:"m_viewPunchAngle"`
}

type signatures struct {
	ClientstateChokedCommands      int32 `json:"clientstate_choked_commands"`
	ClientstateDeltaTicks          int32 `json:"clientstate_delta_ticks"`
	ClientstateLastOutgoingCommand int32 `json:"clientstate_last_outgoing_command"`
	ClientstateNetChannel          int32 `json:"clientstate_net_channel"`
	ConvarNameHashTable            int32 `json:"convar_name_hash_table"`
	DwClientState                  int32 `json:"dwClientState"`
	DwClientStateGetLocalPlayer    int32 `json:"dwClientState_GetLocalPlayer"`
	DwClientStateIsHLTV            int32 `json:"dwClientState_IsHLTV"`
	DwClientStateMap               int32 `json:"dwClientState_Map"`
	DwClientStateMapDirectory      int32 `json:"dwClientState_MapDirectory"`
	DwClientStateMaxPlayer         int32 `json:"dwClientState_MaxPlayer"`
	DwClientStatePlayerInfo        int32 `json:"dwClientState_PlayerInfo"`
	DwClientStateState             int32 `json:"dwClientState_State"`
	DwClientStateViewAngles        int32 `json:"dwClientState_ViewAngles"`
	DwEntityList                   int32 `json:"dwEntityList"`
	DwForceAttack                  int32 `json:"dwForceAttack"`
	DwForceAttack2                 int32 `json:"dwForceAttack2"`
	DwForceBackward                int32 `json:"dwForceBackward"`
	DwForceForward                 int32 `json:"dwForceForward"`
	DwForceJump                    int32 `json:"dwForceJump"`
	DwForceLeft                    int32 `json:"dwForceLeft"`
	DwForceRight                   int32 `json:"dwForceRight"`
	DwGameDir                      int32 `json:"dwGameDir"`
	DwGameRulesProxy               int32 `json:"dwGameRulesProxy"`
	DwGetAllClasses                int32 `json:"dwGetAllClasses"`
	DwGlobalVars                   int32 `json:"dwGlobalVars"`
	DwGlowObjectManager            int32 `json:"dwGlowObjectManager"`
	DwInput                        int32 `json:"dwInput"`
	DwInterfaceLinkList            int32 `json:"dwInterfaceLinkList"`
	DwLocalPlayer                  int32 `json:"dwLocalPlayer"`
	DwMouseEnable                  int32 `json:"dwMouseEnable"`
	DwMouseEnablePtr               int32 `json:"dwMouseEnablePtr"`
	DwPlayerResource               int32 `json:"dwPlayerResource"`
	DwRadarBase                    int32 `json:"dwRadarBase"`
	DwSensitivity                  int32 `json:"dwSensitivity"`
	DwSensitivityPtr               int32 `json:"dwSensitivityPtr"`
	DwSetClanTag                   int32 `json:"dwSetClanTag"`
	DwViewMatrix                   int32 `json:"dwViewMatrix"`
	DwWeaponTable                  int32 `json:"dwWeaponTable"`
	DwWeaponTableIndex             int32 `json:"dwWeaponTableIndex"`
	DwYawPtr                       int32 `json:"dwYawPtr"`
	DwZoomSensitivityRatioPtr      int32 `json:"dwZoomSensitivityRatioPtr"`
	DwbSendPackets                 int32 `json:"dwbSendPackets"`
	DwppDirect3DDevice9            int32 `json:"dwppDirect3DDevice9"`
	ForceUpdateSpectatorGlow       int32 `json:"force_update_spectator_glow"`
	InterfaceEngineCvar            int32 `json:"interface_engine_cvar"`
	IsC4Owner                      int32 `json:"is_c4_owner"`
	MBDormant                      int32 `json:"m_bDormant"`
	MPStudioHdr                    int32 `json:"m_pStudioHdr"`
	MPitchClassPtr                 int32 `json:"m_pitchClassPtr"`
	MYawClassPtr                   int32 `json:"m_yawClassPtr"`
	ModelAmbientMin                int32 `json:"model_ambient_min"`
	SetAbsAngles                   int32 `json:"set_abs_angles"`
	SetAbsOrigin                   int32 `json:"set_abs_origin"`
}
