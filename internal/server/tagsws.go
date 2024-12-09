package server

// define some tags for requests that may be received from clients (the same definitions be be found on client code)
const JsonTagPingRequest string = "ping"
const JsonTagCreateGameRequest string = "creategame"
const JsonTagCreateLobbyRequest string = "createlobby"
const JsonTagAddPlayerRequest string = "addplayer"
const JsonTagPlayer string = "playerstate"
const JsonTagBall string = "ballstate"
const JsonTagGame string = "game"
