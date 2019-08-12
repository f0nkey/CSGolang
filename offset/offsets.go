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

	fmt.Println("bod",body)

	var off offsets
	err = json.Unmarshal(body, &off)
	if err != nil {
		log.Println("MarshalOffsets:", err)
	}
	fmt.Println("MarshalOffsets:", err)

	if offsetsOutdated(off) {
		fmt.Println("OFFSETS POSSIBLY OUTDATED")
	}

	fmt.Println("of",off)

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
	Signatures signatures    `json:"signatures"`
	Netvars netvars `json:"netvars"`
}

type netvars struct {
	CsGamerulesData                uintptr `json:"cs_gamerules_data"`
	MArmorValue                    uintptr `json:"m_ArmorValue"`
	MCollision                     uintptr `json:"m_Collision"`
	MCollisionGroup                uintptr `json:"m_CollisionGroup"`
	MLocal                         uintptr `json:"m_Local"`
	MMoveType                      uintptr `json:"m_MoveType"`
	MOriginalOwnerXuidHigh         uintptr `json:"m_OriginalOwnerXuidHigh"`
	MOriginalOwnerXuidLow          uintptr `json:"m_OriginalOwnerXuidLow"`
	MSurvivalGameRuleDecisionTypes uintptr `json:"m_SurvivalGameRuleDecisionTypes"`
	MSurvivalRules                 uintptr `json:"m_SurvivalRules"`
	MAimPunchAngle                 uintptr `json:"m_aimPunchAngle"`
	MAimPunchAngleVel              uintptr `json:"m_aimPunchAngleVel"`
	MAngEyeAnglesX                 uintptr `json:"m_angEyeAnglesX"`
	MAngEyeAnglesY                 uintptr `json:"m_angEyeAnglesY"`
	MBBombPlanted                  uintptr `json:"m_bBombPlanted"`
	MBFreezePeriod                 uintptr `json:"m_bFreezePeriod"`
	MBGunGameImmunity              uintptr `json:"m_bGunGameImmunity"`
	MBHasDefuser                   uintptr `json:"m_bHasDefuser"`
	MBHasHelmet                    uintptr `json:"m_bHasHelmet"`
	MBInReload                     uintptr `json:"m_bInReload"`
	MBIsDefusing                   uintptr `json:"m_bIsDefusing"`
	MBIsQueuedMatchmaking          uintptr `json:"m_bIsQueuedMatchmaking"`
	MBIsScoped                     uintptr `json:"m_bIsScoped"`
	MBIsValveDS                    uintptr `json:"m_bIsValveDS"`
	MBSpotted                      uintptr `json:"m_bSpotted"`
	MBSpottedByMask                uintptr `json:"m_bSpottedByMask"`
	MBStartedArming                uintptr `json:"m_bStartedArming"`
	MClrRender                     uintptr `json:"m_clrRender"`
	MDwBoneMatrix                  uintptr `json:"m_dwBoneMatrix"`
	MFAccuracyPenalty              uintptr `json:"m_fAccuracyPenalty"`
	MFFlags                        uintptr `json:"m_fFlags"`
	MFlC4Blow                      uintptr `json:"m_flC4Blow"`
	MFlDefuseCountDown             uintptr `json:"m_flDefuseCountDown"`
	MFlDefuseLength                uintptr `json:"m_flDefuseLength"`
	MFlFallbackWear                uintptr `json:"m_flFallbackWear"`
	MFlFlashDuration               uintptr `json:"m_flFlashDuration"`
	MFlFlashMaxAlpha               uintptr `json:"m_flFlashMaxAlpha"`
	MFlLastBoneSetupTime           uintptr `json:"m_flLastBoneSetupTime"`
	MFlLowerBodyYawTarget          uintptr `json:"m_flLowerBodyYawTarget"`
	MFlNextAttack                  uintptr `json:"m_flNextAttack"`
	MFlNextPrimaryAttack           uintptr `json:"m_flNextPrimaryAttack"`
	MFlSimulationTime              uintptr `json:"m_flSimulationTime"`
	MFlTimerLength                 uintptr `json:"m_flTimerLength"`
	MHActiveWeapon                 uintptr `json:"m_hActiveWeapon"`
	MHMyWeapons                    uintptr `json:"m_hMyWeapons"`
	MHObserverTarget               uintptr `json:"m_hObserverTarget"`
	MHOwner                        uintptr `json:"m_hOwner"`
	MHOwnerEntity                  uintptr `json:"m_hOwnerEntity"`
	MIAccountID                    uintptr `json:"m_iAccountID"`
	MIClip1                        uintptr `json:"m_iClip1"`
	MICompetitiveRanking           uintptr `json:"m_iCompetitiveRanking"`
	MICompetitiveWins              uintptr `json:"m_iCompetitiveWins"`
	MICrosshairID                  uintptr `json:"m_iCrosshairId"`
	MIEntityQuality                uintptr `json:"m_iEntityQuality"`
	MIFOV                          uintptr `json:"m_iFOV"`
	MIFOVStart                     uintptr `json:"m_iFOVStart"`
	MIGlowIndex                    uintptr `json:"m_iGlowIndex"`
	MIHealth                       uintptr `json:"m_iHealth"`
	MIItemDefinitionIndex          uintptr `json:"m_iItemDefinitionIndex"`
	MIItemIDHigh                   uintptr `json:"m_iItemIDHigh"`
	MIMostRecentModelBoneCounter   uintptr `json:"m_iMostRecentModelBoneCounter"`
	MIObserverMode                 uintptr `json:"m_iObserverMode"`
	MIShotsFired                   uintptr `json:"m_iShotsFired"`
	MIState                        uintptr `json:"m_iState"`
	MITeamNum                      uintptr `json:"m_iTeamNum"`
	MLifeState                     uintptr `json:"m_lifeState"`
	MNFallbackPauintptrKit         uintptr `json:"m_nFallbackPauintptrKit"`
	MNFallbackSeed                 uintptr `json:"m_nFallbackSeed"`
	MNFallbackStatTrak             uintptr `json:"m_nFallbackStatTrak"`
	MNForceBone                    uintptr `json:"m_nForceBone"`
	MNTickBase                     uintptr `json:"m_nTickBase"`
	MRgflCoordinateFrame           uintptr `json:"m_rgflCoordinateFrame"`
	MSzCustomName                  uintptr `json:"m_szCustomName"`
	MSzLastPlaceName               uintptr `json:"m_szLastPlaceName"`
	MThirdPersonViewAngles         uintptr `json:"m_thirdPersonViewAngles"`
	MVecOrigin                     uintptr `json:"m_vecOrigin"`
	MVecVelocity                   uintptr `json:"m_vecVelocity"`
	MVecViewOffset                 uintptr `json:"m_vecViewOffset"`
	MViewPunchAngle                uintptr `json:"m_viewPunchAngle"`
}

type signatures struct {
	ClientstateChokedCommands      uintptr `json:"clientstate_choked_commands"`
	ClientstateDeltaTicks          uintptr `json:"clientstate_delta_ticks"`
	ClientstateLastOutgoingCommand uintptr `json:"clientstate_last_outgoing_command"`
	ClientstateNetChannel          uintptr `json:"clientstate_net_channel"`
	ConvarNameHashTable            uintptr `json:"convar_name_hash_table"`
	DwClientState                  uintptr `json:"dwClientState"`
	DwClientStateGetLocalPlayer    uintptr `json:"dwClientState_GetLocalPlayer"`
	DwClientStateIsHLTV            uintptr `json:"dwClientState_IsHLTV"`
	DwClientStateMap               uintptr `json:"dwClientState_Map"`
	DwClientStateMapDirectory      uintptr `json:"dwClientState_MapDirectory"`
	DwClientStateMaxPlayer         uintptr `json:"dwClientState_MaxPlayer"`
	DwClientStatePlayerInfo        uintptr `json:"dwClientState_PlayerInfo"`
	DwClientStateState             uintptr `json:"dwClientState_State"`
	DwClientStateViewAngles        uintptr `json:"dwClientState_ViewAngles"`
	DwEntityList                   uintptr `json:"dwEntityList"`
	DwForceAttack                  uintptr `json:"dwForceAttack"`
	DwForceAttack2                 uintptr `json:"dwForceAttack2"`
	DwForceBackward                uintptr `json:"dwForceBackward"`
	DwForceForward                 uintptr `json:"dwForceForward"`
	DwForceJump                    uintptr `json:"dwForceJump"`
	DwForceLeft                    uintptr `json:"dwForceLeft"`
	DwForceRight                   uintptr `json:"dwForceRight"`
	DwGameDir                      uintptr `json:"dwGameDir"`
	DwGameRulesProxy               uintptr `json:"dwGameRulesProxy"`
	DwGetAllClasses                uintptr `json:"dwGetAllClasses"`
	DwGlobalVars                   uintptr `json:"dwGlobalVars"`
	DwGlowObjectManager            uintptr `json:"dwGlowObjectManager"`
	DwInput                        uintptr `json:"dwInput"`
	DwuintptrerfaceLinkList        uintptr `json:"dwuintptrerfaceLinkList"`
	DwLocalPlayer                  uintptr `json:"dwLocalPlayer"`
	DwMouseEnable                  uintptr `json:"dwMouseEnable"`
	DwMouseEnablePtr               uintptr `json:"dwMouseEnablePtr"`
	DwPlayerResource               uintptr `json:"dwPlayerResource"`
	DwRadarBase                    uintptr `json:"dwRadarBase"`
	DwSensitivity                  uintptr `json:"dwSensitivity"`
	DwSensitivityPtr               uintptr `json:"dwSensitivityPtr"`
	DwSetClanTag                   uintptr `json:"dwSetClanTag"`
	DwViewMatrix                   uintptr `json:"dwViewMatrix"`
	DwWeaponTable                  uintptr `json:"dwWeaponTable"`
	DwWeaponTableIndex             uintptr `json:"dwWeaponTableIndex"`
	DwYawPtr                       uintptr `json:"dwYawPtr"`
	DwZoomSensitivityRatioPtr      uintptr `json:"dwZoomSensitivityRatioPtr"`
	DwbSendPackets                 uintptr `json:"dwbSendPackets"`
	DwppDirect3DDevice9            uintptr `json:"dwppDirect3DDevice9"`
	ForceUpdateSpectatorGlow       uintptr `json:"force_update_spectator_glow"`
	uintptrerfaceEngineCvar        uintptr `json:"uintptrerface_engine_cvar"`
	IsC4Owner                      uintptr `json:"is_c4_owner"`
	MBDormant                      uintptr `json:"m_bDormant"`
	MPStudioHdr                    uintptr `json:"m_pStudioHdr"`
	MPitchClassPtr                 uintptr `json:"m_pitchClassPtr"`
	MYawClassPtr                   uintptr `json:"m_yawClassPtr"`
	ModelAmbientMin                uintptr `json:"model_ambient_min"`
	SetAbsAngles                   uintptr `json:"set_abs_angles"`
	SetAbsOrigin                   uintptr `json:"set_abs_origin"`
}
