import "./styles.css"
import "./chess/cburnett/move.svg"
import { createGame, importGame, fetchMove } from "./chess/board.js"
import newGame, { existingGame } from "./websockets.js"
import { viewStartHistory, viewBackHistory, viewForwardHistory, viewCurrentGame, undoMove, changePlayerSide, flipBoard } from './chess/modify.js'


const game = await existingGame(313, document.querySelector('#root'))
addControls(game)
console.log(game.id)

function addControls(chessGame){
    document.addEventListener('keydown', (e) => {
        if(e.code === "ArrowLeft") {
            viewBackHistory(chessGame)
        }
        if(e.code === "ArrowRight") {
            viewForwardHistory(chessGame)
        }
        if(e.code === "ArrowUp") {
            viewStartHistory(chessGame)
        }
        if(e.code === "ArrowDown") {
            viewCurrentGame(chessGame)
        } 
        if(e.code === "KeyF") {
            flipBoard(chessGame)
        }
        if(e.code === "KeyZ") {
            undoMove(chessGame)
        }
    })
    document.querySelector('#current').addEventListener('click', () => viewCurrentGame(chessGame))
    document.querySelector('#start').addEventListener('click', () => viewStartHistory(chessGame))
    document.querySelector('#back').addEventListener('click', () => viewBackHistory(chessGame))
    document.querySelector('#forwards').addEventListener('click', () => viewForwardHistory(chessGame))
    document.querySelector('#flip').addEventListener('click', () => flipBoard(chessGame))
    document.querySelector('#takeback').addEventListener('click', () => undoMove(chessGame))
}

export function updateToast(text) {
    const toast = document.querySelector('#toast')
    if(text.result) {
        if(text.result === "Stalemate") {
            toast.innerHTML = `Game ended in Stalemate due to ${text.reason}.`
        } else {
            toast.innerHTML = `${text.result} has won due to ${text.reason}`
        }
    } else {
        toast.innerHTML = text
    }   
}





