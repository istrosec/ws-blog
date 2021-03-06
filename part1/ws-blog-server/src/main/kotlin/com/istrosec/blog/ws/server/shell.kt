package com.istrosec.blog.ws.server

import io.ktor.http.cio.websocket.*

internal suspend fun passShellMessages(
    sender: WebSocketSession,
    senders: ShellSessions,
    receivers: ShellSessions
) {
    val agent = receiveAgentInfoOrClose(sender) ?: return
    println("Received agent $agent from ${senders.name}")
    senders.register(agent, sender)
    try {
        for (frame in sender.incoming) {
            if(frame is Frame.Text) {
                println("Receiving for agent $agent from ${senders.name} for ${receivers.name}")
                val message = frame.readText()
                println("Received for agent $agent from ${senders.name} for ${receivers.name}: $message")
                receivers[agent]?.send(message)
            }
        }
    } catch (e: Exception) {
        e.printStackTrace()
    } finally {
        println("Ending $senders session for $agent")
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
