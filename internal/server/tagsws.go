package server

// define some tags for requests that may be received from clients (the same definitions be be found on client code)
const JsonTagPingMsg string = "ping"
const JsonTagCreateGameMsg string = "creategame"
const JsonTagCreateLobbyMsg string = "createlobby"
const JsonTagCheckLobbyMsg string = "checklobby"
const JsonTagAdmissionMsg string = "admission"
const JsonTagAddPlayerMsg string = "addplayergame"
const JsonTagAddPlayerLobby string = "addplayerlobby"
const JsonTagRemPlayerLobby string = "leavelobby"
const JsonTagRemPlayerGame string = "leavegame"
const JsonTagPlayerEvent string = "playeraction"
const JsonTagBallEvent string = "ballstate"
const JsonTagGame string = "game"
