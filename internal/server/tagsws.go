package server

// define some tags for requests that may be received from clients (the same definitions be be found on client code)
const JsonTagPingRequest string = "ping"
const JsonTagCreateGameRequest string = "creategame"
const JsonTagCreateLobbyRequest string = "createlobby"
const JsonTagAdmissionRequest string = "admission"
const JsonTagAddPlayerRequest string = "addplayergame"
const JsonTagAddPlayerLobby string = "addplayerlobby"
const JsonTagPlayerAction string = "playeraction"
const JsonTagBall string = "ballstate"
const JsonTagGame string = "game"
