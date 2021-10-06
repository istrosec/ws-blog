package com.istrosec.bloh.ws.server

import io.ktor.http.cio.websocket.*

internal suspend fun passShellMessages(
    sender: WebSocketSession,
    senders: ShellSessions,
    receivers: ShellSessions
) {
    val agent = receiveAgentInfoOrClose(sender) ?: return
    senders.register(agent, sender)
    try {
        for (frame in sender.incoming) {
            if(frame is Frame.Text) {
                val message = frame.readText()
                receivers[agent]?.send(message)
            }
        }
    } catch (e: Exception) {
        e.printStackTrace()
    } finally {
        senders.unregister(agent)
    }
}

private suspend fun receiveAgentInfoOrClose(ws: WebSocketSession): Agent? {
    return try {
        val message = (ws.incoming.receive() as Frame.Text).readText()
        json.readValue(message, Agent::class.java)
    } catch (e: Exception) {
        ws.close(
            CloseReason(
                code = CloseReason.Codes.CANNOT_ACCEPT,
                message = "Expected Agent message as JSON inside the first text Frame."
            )
        )
        e.printStackTrace()
        null
    }
}
